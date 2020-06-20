package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/profiler"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
	"gopkg.in/yaml.v2"

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
var ConfigFilePath = "/go/bin/miao/config.yaml"

var (
	ServerRequestCountView = &view.View{
		Name:        "opencensus.io/http/server/request_count",
		Description: "Count of HTTP requests started",
		Measure:     ochttp.ServerRequestCount,
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{ochttp.KeyServerRoute, ApplicationKeyTag, ApplicationVersionKeyTag},
	}
	ServerLatencyView = &view.View{
		Name:        "opencensus.io/http/server/latency",
		Description: "Latency distribution of HTTP requests",
		Measure:     ochttp.ServerLatency,
		Aggregation: ochttp.DefaultLatencyDistribution,
		TagKeys:     []tag.Key{ochttp.KeyServerRoute, ApplicationKeyTag, ApplicationVersionKeyTag},
	}
	ApplicationKeyTag        = tag.MustNewKey("application")
	ApplicationVersionKeyTag = tag.MustNewKey("version")

	config Config
)

type Config struct {
	Endpoints []string `yaml:"endpoints"`
}

func main() {
	// Make it easier to configure utilizing the same images but reporting different images
	ServiceName = os.Getenv("SERVICE_NAME")

	logger := logrus.New()
	logger.Formatter = sd.NewFormatter(
		sd.WithService(strings.ToLower(ServiceName)),
		sd.WithVersion(Version),
	)

	if err := profiler.Start(profiler.Config{
		Service:        ServiceName,
		ServiceVersion: Version,
		ProjectID:      os.Getenv("GOOGLE_CLOUD_PROJECT"),
	}); err != nil {
		logger.Errorf("Unable to start profiler")
	}

	// Read config file
	rawConfig, err := ioutil.ReadFile(ConfigFilePath)
	if err != nil {
		logger.Errorf("Unable to read file in. %v. ConfigFilePath: %v", err, ConfigFilePath)
		panic("Unable to read file")
	}
	err = yaml.Unmarshal(rawConfig, &config)

	logger.Level = logrus.InfoLevel
	logger.Info("Application Start Up")
	defer logger.Info("Application Ended")

	if err := view.Register(ServerRequestCountView, ServerLatencyView); err != nil {
		log.Fatalf("Failed to register the view: %v", err)
	}

	// Create and register a OpenCensus Stackdriver Trace exporter.
	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: os.Getenv("GOOGLE_CLOUD_PROJECT"),
		Resource: &monitoredres.MonitoredResource{
			Type: "gke_container",
			Labels: map[string]string{
				"project_id":     os.Getenv("GOOGLE_CLOUD_PROJECT"),
				"namespace_id":   os.Getenv("MY_POD_NAMESPACE"),
				"pod_id":         os.Getenv("MY_POD_NAME"),
				"cluster_name":   os.Getenv("CLUSTER_NAME"),
				"container_name": os.Getenv("CONTAINER_NAME"),
				"instance_id":    os.Getenv("INSTANCE_ID"),
				"zone":           os.Getenv("ZONE"),
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

	http.Handle("/version", MonitoringHandler{
		Zzz:   VersionHandler{Logger: logger},
		Route: "/version",
	})
	http.Handle("/", MonitoringHandler{
		Zzz:   MainHandler{Logger: logger, Client: client},
		Route: "/",
	})
	http.ListenAndServe(":8080", httpHandler)
}

type MonitoringHandler struct {
	Zzz   http.Handler
	Route string
}

func (m MonitoringHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	ctx, _ := tag.New(
		r.Context(),
		tag.Insert(ochttp.KeyServerRoute, m.Route),
		tag.Insert(ApplicationKeyTag, ServiceName),
		tag.Insert(ApplicationVersionKeyTag, Version),
	)
	defer func() {
		stats.Record(
			ctx,
			ochttp.ServerRequestCount.M(1),
			ochttp.ServerLatency.M(float64(time.Now().Sub(startTime).Milliseconds())),
		)
	}()
	m.Zzz.ServeHTTP(w, r)
}

type MainHandler struct {
	Logger *logrus.Logger
	Client *http.Client
}

func (m MainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, span := trace.StartSpan(r.Context(), fmt.Sprintf("%v mainHandler", ServiceName))
	defer span.End()

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

	if len(config.Endpoints) == 0 {
		m.Logger.Info("No endpoints has been configured here")
	}

	// Loop through the endpoints
	for _, extEndpoint := range config.Endpoints {
		_, span := trace.StartSpan(r.Context(), fmt.Sprintf("%v mainHandler - call endpoint %v", ServiceName, extEndpoint))
		defer span.End()
		req, _ := http.NewRequest("GET", extEndpoint, nil)
		req = req.WithContext(r.Context())
		resp, err := m.Client.Do(req)
		if err != nil {
			m.Logger.Errorf("Unable to contact endpoint: %v", extEndpoint)
			continue
		}
		zz, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			m.Logger.Errorf("HTTP response unsuccessful. StatusCode: %v. Response: %v", resp.StatusCode, string(zz))
			continue
		}
		m.Logger.Infof("HTTP response successful. Response: %v", string(zz))
	}
	fmt.Fprintf(w, "Hello World: %s!\n", target)
}

type VersionHandler struct {
	Logger *logrus.Logger
}

func (v VersionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Logger.WithField("kcmaclac", "acacacca").Info("Start version handler")
	defer v.Logger.Info("End version handler")
	fmt.Fprintf(w, "Version of app: %s\n", Version)
}
