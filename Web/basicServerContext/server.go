package main

import (
	"log"
	"net/http"
	"fmt"
	"time"
)



func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handler Started")
	defer log.Printf("handler Ended")

	ctx := r.Context()

	select {
	case <- time.After(5 * time.Second):
		fmt.Fprint(w, "hello")
	case <- ctx.Done():
		err := ctx.Err()
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}


func main() {
	log.Println("Server Started")

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}