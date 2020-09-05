package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	logger "github.com/sirupsen/logrus"
)

var queueItemGauge = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "testservice_generated_queue_item",
	Help: "Current number of items in queue",
})

type generateRateSpike struct {
	// Maximum should be 1000
	Rate int
	// In seconds
	Duration int
}

var items int
var queue chan int
var spikeQueue chan generateRateSpike
var durationRest float64

func main() {
	logger.Info("Hello world sample started.")
	queue = make(chan int)
	spikeQueue = make(chan generateRateSpike, 10)

	queueItemGenerationRaw := os.Getenv("QUEUE_ITEM_GENERATION")
	if queueItemGenerationRaw == "" {
		queueItemGenerationRaw = "1"
	}

	queueItemGeneration, err := strconv.ParseFloat(queueItemGenerationRaw, 64)
	if err != nil {
		queueItemGeneration = 1.0
	}
	durationRest = 1000.0 / queueItemGeneration

	go func() {
		for {
			time.Sleep(time.Duration(durationRest) * time.Millisecond)
			queue <- 1
		}
	}()
	go itemProcessor()
	go spikeProcessor()

	http.HandleFunc("/dec", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Queue item decreased")
		queue <- -1
		w.WriteHeader(http.StatusOK)
	})
	http.HandleFunc("/spike", func(w http.ResponseWriter, r *http.Request) {
		increasedRateRaw := r.URL.Query().Get("rate")
		durationRaw := r.URL.Query().Get("duration")
		increasedRate, _ := strconv.Atoi(increasedRateRaw)
		duration, _ := strconv.Atoi(durationRaw)
		spikeQueue <- generateRateSpike{
			Rate:     increasedRate,
			Duration: duration,
		}
		w.WriteHeader(http.StatusOK)
	})
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}

func spikeProcessor() {
	for {
		select {
		case spikeItem := <-spikeQueue:
			logger.Infof("Processing spike. Spike: %v", spikeItem)
			originalDurationRest := durationRest
			if spikeItem.Rate != 0 && spikeItem.Rate < 1000 {
				logger.Infof("Spike Item Rate: %v", spikeItem.Rate)
				durationRest = 1000.0 / float64(spikeItem.Rate)
			}
			time.Sleep(time.Duration(spikeItem.Duration) * time.Second)
			durationRest = originalDurationRest
		}
	}
}

func itemProcessor() {
	for {
		select {
		case val := <-queue:
			logger.Info("Running addition")
			if val >= 1 {
				items = items + 1
			}
			if val <= 0 {
				if items > 0 {
					items = items - 1
				} else {
					items = 0
				}
			}
			queueItemGauge.Set(float64(items))
		}
	}
}
