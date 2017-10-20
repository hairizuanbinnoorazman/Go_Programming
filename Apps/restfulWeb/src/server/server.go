package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

var people []Person

type Person struct {
	ID        string   `json:"id,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}

type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"'`
}

type insertPersonHandler struct{}

func (h insertPersonHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err.Error())
	}

	var person Person
	json.Unmarshal(data, &person)
	people = append(people, person)

	type reply struct {
		Id     string `json:"id"`
		Status string `json:"status"`
	}

	replyData := reply{Id: person.ID, Status: "success"}
	if err != nil {
		log.Println("Error miao", err.Error())
	}
	w.WriteHeader(201)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(replyData)
}

func main() {
	log.Println("server started")
	router := mux.NewRouter()
	router.Handle("/people", insertPersonHandler{}).Methods("POST")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", router))
}
