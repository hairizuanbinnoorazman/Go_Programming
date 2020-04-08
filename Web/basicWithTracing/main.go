package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

func handler(w http.ResponseWriter, r *http.Request) {
	tracer := opentracing.GlobalTracer()
	spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	serverSpan := tracer.StartSpan("server", ext.RPCServerOption(spanCtx))
	defer serverSpan.Finish()

	log.Print("Hello world received a request.")
	defer log.Print("End hello world request")
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
		req, _ := http.NewRequest("GET", url, nil)

		ext.SpanKindRPCClient.Set(serverSpan)
		ext.HTTPUrl.Set(serverSpan, url)
		ext.HTTPMethod.Set(serverSpan, "GET")

		tracer.Inject(serverSpan.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
		http.DefaultClient.Do(req)
	}
}

func main() {
	log.Print("Hello world sample started.")

	jaegerCollector := os.Getenv("JAEGER_COLLECTOR")
	if jaegerCollector == "" {
		log.Println("Will be utilizing local Jaeger agent")
	}

	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			CollectorEndpoint: jaegerCollector,
			LogSpans:          true,
		},
	}

	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "NOT SPECIFIED"
	}

	// Initialize tracer with a logger and a metrics factory
	closer, _ := cfg.InitGlobalTracer(
		serviceName,
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	defer closer.Close()

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
