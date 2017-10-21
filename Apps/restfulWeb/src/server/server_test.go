/*
Testing script to check that all required endpoints are made available accordingly

We would be creating a user detail service that would serve details on a person's information
Endpoints to consider
- Get single person
- Get a list of people
- Post Create
- Delete person
- Patch person (TODO)
- Put person (TODO)

Important point to take note:
- Ensure that items in struct is capitalized - if not capitalized, will have issues to convert between json and struct

Some of the stuff that is being added here
- Try out table driven tests
*/

package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreate(t *testing.T) {
	cases := []Person{
		{"1000", "Ann", "Mcdonald", &Address{"Singapore", "Singapore"}},
	}

	server := httptest.NewServer(router())
	defer server.Close()
	for idx, person := range cases {
		data, err := json.Marshal(person)
		if err != nil {
			t.Fatal("Unable to process cases properly. Please inspect accordingly")
		}
		res, err := http.Post(server.URL+"/people/", "application/json", bytes.NewReader(data))
		if err != nil {
			t.Error("Unable to add the following case", idx, "Error:", err.Error())
		}
		if res.StatusCode != 201 {
			t.Error("Expected status code: 201. Received:", res.StatusCode)
		}

		// Process data
		t.Log(res.Header)
		resData, err := ioutil.ReadAll(res.Body)
		t.Log(string(resData))
		type extractedResponse struct {
			Id     string
			Status string
		}
		var resDataStruct extractedResponse
		json.Unmarshal(resData, &resDataStruct)

		if resDataStruct.Id != person.ID {
			t.Error("Expected id:", person.ID, "Received", resDataStruct.Id)
		}
		if resDataStruct.Status != "success" {
			t.Error("Expected status: success. Received", resDataStruct.Status)
		}
	}
}

func TestGet(t *testing.T) {
	server := httptest.NewServer(router())
	defer server.Close()

	data, err := http.Get(server.URL + "/people/999/")
	t.Log(server.URL + "/people/999/")
	if err != nil {
		t.Error("Server endpoint not available.", err.Error())
	}

	dataRaw, err := ioutil.ReadAll(data.Body)
	if err != nil {
		t.Error("Unexpected error in trying to get body of the data", err.Error())
	}

	var response Person
	json.Unmarshal(dataRaw, &response)

	if response.ID != "999" {
		t.Error("Expected 999. Received", response.ID)
	}
	if response.Lastname != "miao" {
		t.Error("Expected miao. Received", response.Lastname)
	}
	if response.Address.City != "miao" {
		t.Error("Expected miao. Received", response.Address.City)
	}
}

func TestGetAll(t *testing.T) {
	server := httptest.NewServer(router())
	defer server.Close()

	res, err := http.Get(server.URL + "/people/")
	if err != nil {
		t.Fatal("Unable to hit the people endpoint")
	}

	if res.StatusCode != 200 {
		t.Error("Wrong error status")
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal("Unable to read the response from people endpoint")
	}

	var jsonData []Person
	json.Unmarshal(data, &jsonData)

	if len(jsonData) != 1 {
		t.Error("Expected length of data to be returned to be 1 got:", len(jsonData))
	}
	if len(jsonData) > 0 {
		if jsonData[0].ID != "999" {
			t.Error("Expected ID: 999 Received:", jsonData[0].ID)
		}
	}
}

func TestDelete(t *testing.T) {
	cases := []string{
		"999",
	}

	server := httptest.NewServer(router())
	defer server.Close()

	client := &http.Client{}

	for _, value := range cases {
		req, err := http.NewRequest(http.MethodDelete, server.URL+"/people/"+value, nil)
		if err != nil {
			t.Fatal("Unable to hit endpoint")
		}
		res, err := client.Do(req)
		if err != nil {
			t.Fatal("Unable to hit delete endpoint")
		}
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal("Unable to parse data from the delete endpoint")
		}

		type deleteResponse struct {
			Id           string `json:"id"`
			DeleteStatus string `json:"delete_status"`
		}

		var status deleteResponse
		err = json.Unmarshal(data, &status)
		if err != nil {
			t.Fatal("Unable to parse data from delete endpoint")
		}
		if status.Id != value {
			t.Error("Expected Deleted ID:", value, "Received:", status.Id)
		}
		if status.DeleteStatus != "success" {
			t.Error("Expected DeleteStatus: success Received:", status.DeleteStatus)
		}
	}

}
