package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/yi-jiayu/datamall"
)

func main() {
	log.Println("Server start")

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Println("No API Key to access LTA datamall dataset")
		os.Exit(1)
	}

	// appSecret := os.Getenv("APP_SECRET")
	// if appSecret == "" {
	// 	log.Println("No API Secret for server")
	// 	os.Exit(1)
	// }

	c := http.DefaultClient
	c.Timeout = time.Second * 10

	datamallClient := datamall.NewClient(apiKey, c)

	r := mux.NewRouter()
	r.Handle("/", admin{}).Methods(http.MethodGet)
	s := r.PathPrefix("/api/lta-datamall/v1").Subrouter()
	s.Handle("/healthz", admin{})
	s.Handle("/bus-arrival", arrival{apiKey: apiKey, dc: datamallClient})

	http.ListenAndServe(":8080", r)
}

type admin struct{}

func (a admin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Start Admin Endpoint")
	defer log.Println("End Admin Endpoint")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

type arrival struct {
	apiKey string
	dc     datamall.APIClient
}

func (a arrival) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Start Bus Arrival Endpoint")
	defer log.Println("End Bus Arrival Endpoint")

	// ak := r.Header.Get("Authorization")
	// if ak != a.appSecret {
	// 	log.Println("bad authorization token passed in. rejected")
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	resp := map[string]string{"status": "unauthorized"}
	// 	rawResp, _ := json.Marshal(resp)
	// 	w.Write(rawResp)
	// 	return
	// }

	busStopID := r.URL.Query().Get("bus-stop-id")
	if len(busStopID) != 5 {
		log.Printf("bad bus stop id passed in: %v\n", r.URL.Query().Encode())
		w.WriteHeader(http.StatusBadRequest)
		resp := map[string]string{"status": "bad bus stop id"}
		rawResp, _ := json.Marshal(resp)
		w.Write(rawResp)
		return
	}

	if busStopID == "99999" || busStopID == "99998" || busStopID == "99997" {
		bar := BusArrivalResponse{}
		bar.BusStopID = busStopID
		bar.Services = []BusServiceArrivals{
			BusServiceArrivals{
				ServiceNo: "55",
				NextBus:   5,
				NextBus2:  10,
				NextBus3:  15,
			},
			BusServiceArrivals{
				ServiceNo: "44",
				NextBus:   3,
				NextBus2:  6,
				NextBus3:  9,
			},
		}
		rawArrivals, _ := json.Marshal(bar)
		w.WriteHeader(http.StatusOK)
		w.Write(rawArrivals)
		return
	}

	arrivals, err := a.dc.GetBusArrivalV2(busStopID, "")
	if err != nil {
		log.Printf("unable to get bus arrival results. Err: %v \n", err)
		w.WriteHeader(http.StatusInternalServerError)
		resp := map[string]string{"status": "unable to retrieve results"}
		rawResp, _ := json.Marshal(resp)
		w.Write(rawResp)
		return
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc)

	bar := BusArrivalResponse{}
	bar.BusStopID = arrivals.BusStopCode
	for _, val := range arrivals.Services {
		nextBusTime, nextBusErr := time.Parse("2006-01-02T15:04:05-07:00", val.NextBus.EstimatedArrival)
		nextBus2Time, nextBus2Err := time.Parse("2006-01-02T15:04:05-07:00", val.NextBus2.EstimatedArrival)
		nextBus3Time, nextBus3Err := time.Parse("2006-01-02T15:04:05-07:00", val.NextBus3.EstimatedArrival)
		bsa := BusServiceArrivals{}
		if nextBusErr != nil && nextBus2Err != nil && nextBus3Err != nil {
			log.Printf("All timings failed to be parsed")
			log.Printf("Service No: %v", val.ServiceNo)
			log.Printf("Next Bus Timing Parse Err: %v, OriginalVal: %v", nextBusErr, val.NextBus.EstimatedArrival)
			log.Printf("Next Bus 2 Timing Parse Err: %v, OriginalVal: %v", nextBus2Err, val.NextBus2.EstimatedArrival)
			log.Printf("Next Bus 3 Timing Parse Err: %v, OriginalVal: %v", nextBus3Err, val.NextBus3.EstimatedArrival)
			continue
		}
		bsa.ServiceNo = val.ServiceNo
		if nextBusErr == nil {
			bsa.NextBus = int(nextBusTime.Sub(now).Minutes())
		}
		if nextBus2Err == nil {
			bsa.NextBus2 = int(nextBus2Time.Sub(now).Minutes())
		}
		if nextBus3Err == nil {
			bsa.NextBus3 = int(nextBus3Time.Sub(now).Minutes())
		}
		bar.Services = append(bar.Services, bsa)
	}

	rawArrivals, err := json.Marshal(bar)
	if err != nil {
		log.Printf("unable to convert to raw JSON response. Err: %v \n", err)
		w.WriteHeader(http.StatusInternalServerError)
		resp := map[string]string{"status": "unable to retrieve results"}
		rawResp, _ := json.Marshal(resp)
		w.Write(rawResp)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(rawArrivals)
}

type BusArrivalResponse struct {
	Services  []BusServiceArrivals `json:"services"`
	BusStopID string               `json:"bus_stop_id"`
}

type BusServiceArrivals struct {
	ServiceNo string `json:"service_no"`
	NextBus   int    `json:"next_bus"`
	NextBus2  int    `json:"next_bus_2"`
	NextBus3  int    `json:"next_bus_3"`
}
