package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
	hehe = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "myapp_hehe",
		Help: "Hehe gauge",
	}, []string{"items", "date"})
)

type ZZZ struct{}

func (z *ZZZ) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	date := fmt.Sprintf(time.RFC3339, time.Now())
	items := r.URL.Query().Get("items")
	if items != "" {
		hehe.WithLabelValues(items, date).Set(1.0)
	}
	zzz := r.URL.Query().Get("delete")
	if zzz != "" {
		hehe.Reset()
	}

	w.Write([]byte("testing"))
}

func main() {
	fmt.Println("Start server")
	http.Handle("/hoho", &ZZZ{})
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
