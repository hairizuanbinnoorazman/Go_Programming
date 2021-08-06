package main

import (
	"fmt"
	"log"
	"net/http"
)

var storeValue = [][]byte{}

func main() {
	port := 8888

	http.HandleFunc("/", helloWorldHandler)

	for i := 0; i < 5000000; i++ {
		a := []byte("abcdefghijklmnopqrstuvwxyz")
		storeValue = append(storeValue, a)
	}
	log.Printf("Size of data: %v", len(storeValue))

	log.Printf("Server starting on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("serving", r.URL)
	defer log.Println("End serving")
}
