package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)


type statusHandler struct {}

func (h statusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	log.Println("Serving the status handler")
	fmt.Fprint(w, "miaoMax")
}

func main() {
	r := mux.NewRouter()
	r.Handle("/api/v2/status", statusHandler{})
	log.Fatal(http.ListenAndServe(":8080", r))
}

