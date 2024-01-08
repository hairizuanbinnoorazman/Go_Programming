package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type basic struct{}

func (b basic) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("started basic handler")
	defer log.Println("ended basic handler")
	w.Write([]byte("successfully called basic handler"))
}

func main() {
	log.Print("App started")

	r := mux.NewRouter()
	r.Handle("/", basic{})
	srv := http.Server{
		Handler: r,
		Addr:    "0.0.0.0:8080",
	}
	log.Fatal(srv.ListenAndServe())
}
