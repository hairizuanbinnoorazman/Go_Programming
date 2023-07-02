package main

import (
	"fmt"
	"io/ioutil"
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
	w.Write([]byte(target))
	return
}

func configPrinter() {
	log.Println("Start Config Printer")
	waitTimeEnv := os.Getenv("WAIT_TIME")
	waitTime, _ := strconv.Atoi(waitTimeEnv)
	for {
		time.Sleep(time.Duration(waitTime) * time.Second)
		configfileLocation := os.Getenv("CONFIG_FILE_LOCATION")
		rawFile, err := ioutil.ReadFile(configfileLocation)
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Println(string(rawFile))
	}

}

func main() {
	log.Print("Hello world sample started.")

	go configPrinter()

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
