package main

import (
	"log"
	"io/ioutil"

	"testing"
	"net/http"
	"net/http/httptest"
)

func TestStatusHandler(t *testing.T) {
	server := httptest.NewServer(statusHandler{})
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/v2/status")
	if err != nil {
		log.Println(err.Error(), "Unable to hit the server")
		t.Error(err.Error(), "Unable to hit the server")
	}
	if resp.StatusCode != 200 {
		log.Println("Server is unhealthy")
		t.Error(err.Error(), "Server is unhealthy")
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error(), "Data is wrong")
		t.Error(err.Error(), "Data is Wrong")
	}

	stringData := string(data)
	log.Println(stringData)
	if stringData != "miaoMax" {
		t.Error("Unexpected Output:", stringData, "instead of miaoMax")
	}
}
