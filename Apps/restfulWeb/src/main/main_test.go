package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestIndexHandler(t *testing.T) {
	server := httptest.NewServer(indexHandler{})
	defer server.Close()

	res, err := http.Get(server.URL)
	if err != nil {
		t.Error("Unexpected Error from server", err.Error())
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err.Error())
	}
	stringData := string(data)
	log.Println(stringData)
	if stringData != "helloworld" {
		t.Error("Expected helloworld but received", stringData)
	}
}

func TestExampleJSONHandler(t *testing.T) {
	server := httptest.NewServer(exampleJSONHandler{})
	defer server.Close()

	u, err := url.Parse(server.URL)
	if err != nil {
		log.Println("Unable to parse the server url")
	}
	q := u.Query()
	q.Set("name", "ann")
	q.Set("response", "json")
	u.RawPath = q.Encode()
	t.Log(u.String())

	res, err := http.Get(u.String())
	if err != nil {
		t.Error("Unexpected error with reaching server")
	}
	data, err := ioutil.ReadAll(res.Body)
	log.Println(string(data))
	if err != nil {
		t.Error("Unable to read value from server", err.Error())
	}
	var user user
	err = json.Unmarshal(data, &user)
	if err != nil {
		t.Error("Unable to convert the values accordingly", err.Error())
	}
	if user.Name != "ann" {
		t.Error("Expected ann but received", user.Name)
	}
	if user.Gender != "female" {
		t.Error("Expected female but received", user.Gender)
	}
}
