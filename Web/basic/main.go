package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
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
	fmt.Fprintf(w, "Hello World: %s!\n", target)
}

type HandleViaStruct struct{}

func (*HandleViaStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Print("Hello world received a request.")
	defer log.Print("End hello world request")
	fmt.Fprintf(w, "Hello World via Struct")
}

type DoHttpReq struct{}

func (d *DoHttpReq) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Start DoHTTPReq")
	defer log.Println("End DoHTTPReq")
	z := r.URL.Query().Get("url")
	log.Printf("Attempting to query the following url %v\n", z)

	resp, err := http.Get(z)
	if err != nil {
		log.Printf("unable to get data from url %v\n", err)
	} else {
		y, errx := io.ReadAll(resp.Body)
		if errx != nil {
			log.Printf("unable to get reading body %v\n", errx)
		} else {
			log.Printf("%v\n", string(y))
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("attempt to query done"))
}

func main() {
	log.Print("Hello world sample started.")

	http.HandleFunc("/", handler)
	http.Handle("/struct", &HandleViaStruct{})
	http.Handle("/query", &DoHttpReq{})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
