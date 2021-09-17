package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

var (
	requestsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "requests_total",
		Help: "The total number of processed events",
	})
)

type LogTransport struct {
	Transport http.RoundTripper
}

func (l LogTransport) RoundTrip(a *http.Request) (*http.Response, error) {
	startTime := time.Now().String()
	resp, err := l.Transport.RoundTrip(a)
	defer log.Printf("Method: %v URL: %v StartTime: %v EndTime: %v ResponseCode: %v\n",
		a.Method, a.URL.String(), startTime, time.Now().String(), resp.StatusCode)
	return resp, err
}

type TraceTransport struct {
	Transport http.RoundTripper
}

func (t TraceTransport) RoundTrip(a *http.Request) (*http.Response, error) {
	tracer := opentracing.GlobalTracer()
	ctx := a.Context()
	parentSpan := opentracing.SpanFromContext(ctx)
	var childSpan opentracing.Span
	if parentSpan == nil {
		childSpan = tracer.StartSpan("client")
	} else {
		childSpan = tracer.StartSpan("client", opentracing.ChildOf(parentSpan.Context()))
	}
	traceID := childSpan.Context().(jaeger.SpanContext)
	defer childSpan.Finish()
	defer log.WithField("traceID", traceID).Info("Done with requests")
	ext.SpanKindRPCClient.Set(childSpan)
	ext.HTTPUrl.Set(childSpan, a.URL.String())
	ext.HTTPMethod.Set(childSpan, a.Method)
	tracer.Inject(childSpan.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(a.Header))
	resp, err := t.Transport.RoundTrip(a)
	if err != nil {
		ext.Error.Set(childSpan, true)
	}
	ext.HTTPStatusCode.Set(childSpan, uint16(resp.StatusCode))
	return resp, err
}

func handler(w http.ResponseWriter, r *http.Request) {
	tracer := opentracing.GlobalTracer()
	spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	serverSpan := tracer.StartSpan("server", ext.RPCServerOption(spanCtx))
	defer serverSpan.Finish()

	log.Print("Hello world received a request.")
	defer log.Print("End hello world request")
	defer requestsTotal.Inc()
	target := os.Getenv("TARGET")
	if target == "" {
		target = "NOT SPECIFIED"
	}
	waitTimeEnv := os.Getenv("WAIT_TIME")
	waitTime, _ := strconv.Atoi(waitTimeEnv)
	log.Printf("Sleeping for %v", waitTime)
	time.Sleep(time.Duration(waitTime) * time.Second)
	fmt.Fprintf(w, "Hello: %s!\n", target)

	clientURL := os.Getenv("CLIENT_URL")
	if clientURL != "" {
		url := clientURL
		ctx := opentracing.ContextWithSpan(r.Context(), serverSpan)
		req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

		traceClient := http.Client{
			Transport: TraceTransport{
				LogTransport{
					Transport: http.DefaultTransport,
				},
			},
		}

		traceClient.Do(req)
	}
}

type StatusHandler struct {
	StatusType string
}

func (s StatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Infof("Status Handler: %v started", s.StatusType)
	zz := map[string]string{"status": "ok"}
	aa, _ := json.Marshal(zz)
	w.WriteHeader(http.StatusOK)
	w.Write(aa)
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.Info("Application started")

	cfg, err := jaegercfg.FromEnv()
	if err != nil {
		// parsing errors might happen here, such as when we get a string where we expect a number
		log.Errorf("Could not parse Jaeger env vars: %s", err.Error())
		panic("Unable to parse jaeger stuff")
	}

	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "NOT SPECIFIED"
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	// Initialize tracer with a logger and a metrics factory
	closer, _ := cfg.InitGlobalTracer(
		serviceName,
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	defer closer.Close()

	http.HandleFunc("/", handler)
	http.Handle("/healthz", StatusHandler{StatusType: "healthz"})
	http.Handle("/readyz", StatusHandler{StatusType: "readyz"})
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(fmt.Sprintf(":%v", serverPort), nil)
}
