package main

import (
	"context"
	"embed"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

//go:embed web/index.html
var webFS embed.FS

type server struct {
	kube      *kubeService
	namespace string
	template  *template.Template
}

func main() {
	var (
		addr       = flag.String("listen", env("LISTEN_ADDR", ":8080"), "HTTP listen address")
		namespace  = flag.String("namespace", env("SANDBOX_NAMESPACE", "agent-sandbox-snapshots"), "namespace to manage")
		templateID = flag.String("sandbox-template", env("SANDBOX_TEMPLATE", "snapshot-python"), "SandboxTemplate name")
		warmPool   = flag.String("warm-pool", env("SANDBOX_WARM_POOL", "snapshot-python"), "SandboxWarmPool name")
		kubeconfig = flag.String("kubeconfig", "", "optional kubeconfig for local development")
	)
	flag.Parse()

	cfg, err := kubernetesConfig(*kubeconfig)
	if err != nil {
		slog.Error("load Kubernetes configuration", "error", err)
		os.Exit(1)
	}
	coreClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		slog.Error("create Kubernetes client", "error", err)
		os.Exit(1)
	}
	dynamicClient, err := dynamic.NewForConfig(cfg)
	if err != nil {
		slog.Error("create dynamic Kubernetes client", "error", err)
		os.Exit(1)
	}

	tmpl, err := template.ParseFS(webFS, "web/index.html")
	if err != nil {
		slog.Error("parse UI template", "error", err)
		os.Exit(1)
	}

	s := &server{
		namespace: *namespace,
		template:  tmpl,
		kube:      newKubeService(coreClient, dynamicClient, *namespace, *templateID, *warmPool),
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", s.dashboard)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})
	mux.HandleFunc("POST /actions", s.action)

	httpServer := &http.Server{
		Addr:              *addr,
		Handler:           requestLog(mux),
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	go func() {
		slog.Info("sandbox demo listening", "address", *addr, "namespace", *namespace)
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server failed", "error", err)
			os.Exit(1)
		}
	}()

	stop, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	<-stop.Done()
	ctx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	_ = httpServer.Shutdown(ctx)
}

func (s *server) dashboard(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	data, err := s.kube.dashboard(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("load Kubernetes state: %v", err), http.StatusInternalServerError)
		return
	}
	data.Message = r.URL.Query().Get("message")
	data.Error = r.URL.Query().Get("error")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.template.ExecuteTemplate(w, "index.html", data); err != nil {
		slog.Error("render dashboard", "error", err)
	}
}

func (s *server) action(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		redirectResult(w, r, "", "invalid form")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 6*time.Minute)
	defer cancel()

	action, name := r.FormValue("action"), strings.TrimSpace(r.FormValue("name"))
	var (
		message string
		err     error
	)
	switch action {
	case "create":
		name, err = s.kube.createClaim(ctx, name)
		message = "created claim " + name
	case "destroy":
		err = s.kube.destroyClaim(ctx, name)
		message = "destroy requested for claim " + name
	case "suspend":
		err = s.kube.setReplicas(ctx, name, 0)
		message = "suspended sandbox " + name
	case "resume":
		err = s.kube.setReplicas(ctx, name, 1)
		message = "resumed sandbox " + name
	case "snapshot":
		var snapshot string
		snapshot, err = s.kube.snapshot(ctx, name)
		message = "memory snapshot ready: " + snapshot
	case "snapshot-suspend":
		var snapshot string
		snapshot, err = s.kube.snapshot(ctx, name)
		if err == nil {
			err = s.kube.setReplicas(ctx, name, 0)
		}
		message = "snapshot " + snapshot + " captured; sandbox suspended"
	case "restore":
		var snapshot string
		snapshot, err = s.kube.restoreLatest(ctx, name)
		message = "restoring sandbox " + name + " from " + snapshot
	default:
		err = fmt.Errorf("unknown action %q", action)
	}
	if err != nil {
		redirectResult(w, r, "", err.Error())
		return
	}
	redirectResult(w, r, message, "")
}

func redirectResult(w http.ResponseWriter, r *http.Request, message, errorMessage string) {
	q := make(map[string]string)
	q["message"], q["error"] = message, errorMessage
	target := "/?"
	if errorMessage != "" {
		target += "error=" + urlQueryEscape(errorMessage)
	} else {
		target += "message=" + urlQueryEscape(message)
	}
	http.Redirect(w, r, target, http.StatusSeeOther)
}

func urlQueryEscape(value string) string {
	r := strings.NewReplacer("%", "%25", " ", "+", "&", "%26", "?", "%3F", "#", "%23", "=", "%3D")
	return r.Replace(value)
}

func kubernetesConfig(path string) (*rest.Config, error) {
	if path != "" {
		return clientcmd.BuildConfigFromFlags("", path)
	}
	cfg, err := rest.InClusterConfig()
	if err == nil {
		return cfg, nil
	}
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{}).ClientConfig()
}

func env(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

func requestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		next.ServeHTTP(w, r)
		slog.Info("request", "method", r.Method, "path", r.URL.Path, "duration", time.Since(started))
	})
}
