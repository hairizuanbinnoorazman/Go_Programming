package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	endpoint := os.Getenv("QUEUE_ITEM_SERVER")
	if endpoint == "" {
		endpoint = "http://localhost:8080/dec"
	}

	queueItemConsumptionRaw := os.Getenv("QUEUE_ITEM_CONSUMPTION")
	if queueItemConsumptionRaw == "" {
		queueItemConsumptionRaw = "1"
	}

	queueItemConsumption, err := strconv.ParseFloat(queueItemConsumptionRaw, 64)
	if err != nil {
		queueItemConsumption = 1.0
	}
	durationRest := 1000.0 / queueItemConsumption

	log.Printf("Endpoint: %v", endpoint)

	for {
		time.Sleep(time.Duration(durationRest) * time.Millisecond)
		resp, err := http.Get(endpoint)
		if err != nil {
			log.Println("Request failed")
		}
		if resp.StatusCode != http.StatusOK {
			log.Println("Request failed")
		}
	}
}
