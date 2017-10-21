package main

import (
	"encoding/json"
	"fmt"
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
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(replyData)
}

type getPersonHandler struct{}

func (h getPersonHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	log.Println(id)
	log.Println(len(people))

	var person Person
	for _, value := range people {
		log.Println(value.ID)
		log.Println(value.Address.City)
		log.Println(value.ID == string(id))
		if value.ID == id {
			person = value
			encoder := json.NewEncoder(w)
			encoder.Encode(person)
			return
		}
	}
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w)
}

func router() *mux.Router {
	// Reassign value for people variable
	people = people[:0]
	people = append(people, Person{"999", "miao", "miao", &Address{"miao", "mia"}})

	// Set up router
	router := mux.NewRouter()
	router.Handle("/people/", insertPersonHandler{}).Methods("POST")
	router.Handle("/people/", getPeopleHandler{}).Methods("GET")
	router.Handle("/people/{id}/", getPersonHandler{}).Methods("GET")
	router.Handle("/people/{id}", deletePeopleHandler{}).Methods("DELETE")
	return router
}

type getPeopleHandler struct{}

func (h getPeopleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	encoder.Encode(people)
}

type deletePeopleHandler struct{}

func (h deletePeopleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Since this is a poor implementation of stored data, there might be duplicates
	// We will be deleting all instances of records that contains 999
	vars := mux.Vars(r)
	id := vars["id"]

	beforeLength := len(people)

	var tempPeople []Person

	for _, value := range people {
		if value.ID != id {
			tempPeople = append(tempPeople, value)
		}
	}

	people = tempPeople
	afterLength := len(people)

	type deleteResponse struct {
		Id           string `json:"id"`
		DeleteStatus string `json:"delete_status"`
	}

	deleteStatus := "failed"
	if beforeLength > afterLength {
		deleteStatus = "success"
	}
	response := deleteResponse{Id: id, DeleteStatus: deleteStatus}

	encoder := json.NewEncoder(w)
	encoder.Encode(response)
}

func main() {
	log.Println("server started")
	log.Println("Create fake data")

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", router()))
}
