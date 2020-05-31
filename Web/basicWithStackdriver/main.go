package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"go.opencensus.io/stats"

	sd "github.com/TV4/logrus-stackdriver-formatter"
	"github.com/sirupsen/logrus"
	"google.golang.org/genproto/googleapis/api/monitoredres"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

var ServiceName = "Basic-With-Stackdriver"
var Version = "0.1.0"

var (
	requestCounter = stats.Int64("request_count", "Number of requests by path", stats.UnitDimensionless)
)

func main() {
	logger := logrus.New()
	logger.Formatter = sd.NewFormatter(
		sd.WithService(ServiceName),
		sd.WithVersion(Version),
	)
	logger.Level = logrus.InfoLevel
	logger.Info("Application Start Up")
	defer logger.Info("Application Ended")

	if err := view.Register(ochttp.DefaultServerViews...); err != nil {
		log.Fatalf("Failed to register the view: %v", err)
	}

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
		logger.Fatal(err)
	}
	defer exporter.Flush()
	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	if err := exporter.StartMetricsExporter(); err != nil {
		logger.Fatalf("Error starting metric exporter: %v", err)
	}
	defer exporter.StopMetricsExporter()

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

	http.Handle("/version", ochttp.WithRouteTag(VersionHandler{Logger: logger}, "/version"))
	http.Handle("/", ochttp.WithRouteTag(MainHandler{Logger: logger, Client: client}, "/"))
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
