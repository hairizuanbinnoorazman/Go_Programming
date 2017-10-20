/*
Testing script to check that all required endpoints are made available accordingly

We would be creating a user detail service that would serve details on a person's information
Endpoints to consider
- Get single person
- Get a list of people
- Post Create
- Patch Update
- Put Update
- Delete person

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
		{"1", "Ann", "Mcdonald", &Address{"Singapore", "Singapore"}},
	}

	server := httptest.NewServer(insertPersonHandler{})
	for idx, person := range cases {
		data, err := json.Marshal(person)
		if err != nil {
			t.Fatal("Unable to process cases properly. Please inspect accordingly")
		}
		res, err := http.Post(server.URL, "application/json", bytes.NewReader(data))
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

}
