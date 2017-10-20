package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type indexHandler struct{}

func (h indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	msg := "helloworld"
	fmt.Fprint(w, msg)
}

type user struct {
	Name   string `json:"name"`
	Gender string `json:"gender"`
}

type exampleJSONHandler struct{}

func (h exampleJSONHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userData := user{Name: "ann", Gender: "female"}
	encoder := json.NewEncoder(w)
	encoder.Encode(userData)
}

func main() {
	http.Handle("/", indexHandler{})
	http.Handle("?name=ann&response=json", exampleJSONHandler{})
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
