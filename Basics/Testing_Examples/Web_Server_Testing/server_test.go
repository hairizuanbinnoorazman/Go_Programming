package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"io/ioutil"
)

func TestHelloWorldHandler(t *testing.T) {
	server := httptest.NewServer(helloWorldHandler{})
	defer server.Close()

	log.Println(server.URL + "/admin")
	resp, err := http.Get(server.URL + "/admin")
	if err != nil{
		log.Println(err.Error())
		t.Error("Wrong Response Obtained?")
	}
	if resp.StatusCode != 200 {
		t.Error("Wrong Status Code")
	}
	data, _ := ioutil.ReadAll(resp.Body)
	if string(data) != "Miao" {
		t.Error("Wrong incoming message")
	}
}