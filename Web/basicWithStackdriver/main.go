package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	sd "github.com/TV4/logrus-stackdriver-formatter"
	"github.com/sirupsen/logrus"
	"google.golang.org/genproto/googleapis/api/monitoredres"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

var ServiceName = "Basic-With-Stackdriver"
var Version = "0.1.0"

func main() {
	logger := logrus.New()
	logger.Formatter = sd.NewFormatter(
		sd.WithService(ServiceName),
		sd.WithVersion(Version),
	)
	logger.Level = logrus.InfoLevel
	logger.Info("Application Start Up")
	defer logger.Info("Application Ended")

	// Create and register a OpenCensus Stackdriver Trace exporter.
	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: os.Getenv("GOOGLE_CLOUD_PROJECT"),
		Resource: &monitoredres.MonitoredResource{
			Type: "gke_container",
			Labels: map[string]string{
				"project_id":   os.Getenv("GOOGLE_CLOUD_PROJECT"),
				"namespace_id": os.Getenv("MY_POD_NAMESPACE"),
				"pod_id":       os.Getenv("MY_POD_NAME"),
			},
		},
		DefaultMonitoringLabels: &stackdriver.Labels{},
	})
	if err != nil {
		log.Fatal(err)
	}
	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	client := &http.Client{
		Transport: &ochttp.Transport{
			// Use Google Cloud propagation format.
			Propagation: &propagation.HTTPFormat{},
		},
	}

	httpHandler := &ochttp.Handler{
		// Use the Google Cloud propagation format.
		Propagation: &propagation.HTTPFormat{},
	}

	http.Handle("/version", VersionHandler{Logger: logger})
	http.Handle("/", MainHandler{Logger: logger, Client: client})
	http.ListenAndServe(":8080", httpHandler)
}

type MainHandler struct {
	Logger *logrus.Logger
	Client *http.Client
}

func (m MainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Logger.WithField("acac", "akcnkackacm").Info("Hello world received a request.")
	defer m.Logger.Infof("End hello world request")
	target := os.Getenv("TARGET")
	if target == "" {
		target = "NOT SPECIFIED"
	}
	waitTimeEnv := os.Getenv("WAIT_TIME")
	waitTime, _ := strconv.Atoi(waitTimeEnv)
	m.Logger.Infof("Sleeping for %v", waitTime)
	time.Sleep(time.Duration(waitTime) * time.Second)
	fmt.Fprintf(w, "Hello World: %s!\n", target)
}

type VersionHandler struct {
	Logger *logrus.Logger
}

func (v VersionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	context.WithValue(ctx, "acacaca", "kcamlcamclmac")
	v.Logger.WithField("kcmaclac", "acacacca").Info("Start version handler")
	defer v.Logger.Info("End version handler")
	fmt.Fprintf(w, "Version of app: %s\n", Version)
}
