package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func firer() {
	hostName := os.Getenv("SERVER_HOST")
	if hostName == "" {
		fmt.Println("hostname not defined. will exit")
		os.Exit(1)
	}
	for {
		ips, err := net.LookupIP(hostName)
		if err != nil {
			fmt.Printf("unexpected error while looking up ips: %v", err)
		}
		for _, ip := range ips {
			fmt.Printf("%v ips found. Will contact ip: %v", len(ips), ip.String())
			time.Sleep(2 * time.Second)
			resp, err := http.Get(fmt.Sprintf("%v:8080", ip.String()))
			if err != nil {
				fmt.Printf("unexpected error when contacting: %v\n", err)
			}
			raw, _ := io.ReadAll(resp.Body)
			fmt.Printf("Output from ip: %v, %v", ip.String(), string(raw))
		}
	}

}

func server() {
	port := 8080

	http.HandleFunc("/", helloWorldHandler)

	log.Printf("Server starting on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("serving", r.URL)
	fmt.Fprintf(w, "This is a test. Hello World Miaoza!! Time: %v\n", time.Now())
}

func main() {
	mode := os.Getenv("MODE")
	if mode == "firer" {
		firer()
	} else if mode == "server" {
		server()
	} else {
		panic("Mode not properly defined. Will terminate")
	}
}
