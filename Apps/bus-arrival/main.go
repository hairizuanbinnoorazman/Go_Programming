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

	appSecret := os.Getenv("APP_SECRET")
	if appSecret == "" {
		log.Println("No API Secret for server")
		os.Exit(1)
	}

	c := http.DefaultClient
	c.Timeout = time.Second * 10

	datamallClient := datamall.NewClient(apiKey, c)

	r := mux.NewRouter()
	r.Handle("/", admin{}).Methods(http.MethodGet)
	s := r.PathPrefix("/api/lta-datamall/v1").Subrouter()
	s.Handle("/healthz", admin{})
	s.Handle("/bus-arrival", arrival{apiKey: apiKey, appSecret: appSecret, dc: datamallClient})

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
	apiKey    string
	appSecret string
	dc        datamall.APIClient
}

func (a arrival) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Start Bus Arrival Endpoint")
	defer log.Println("End Bus Arrival Endpoint")

	ak := r.Header.Get("Authorization")
	if ak != a.appSecret {
		log.Println("bad authorization token passed in. rejected")
		w.WriteHeader(http.StatusUnauthorized)
		resp := map[string]string{"status": "unauthorized"}
		rawResp, _ := json.Marshal(resp)
		w.Write(rawResp)
		return
	}

	busStopID := r.URL.Query().Get("bus-stop-id")
	if len(busStopID) != 5 {
		log.Printf("bad bus stop id passed in: %v\n", busStopID)
		w.WriteHeader(http.StatusBadRequest)
		resp := map[string]string{"status": "bad bus stop id"}
		rawResp, _ := json.Marshal(resp)
		w.Write(rawResp)
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

	rawArrivals, err := json.Marshal(arrivals)
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
