/*
Basic web server implementation in go
To run this file, run the following command

go run json_example_1.go

Some learning points:
- Introduction to struct. Important thing to note is that the property has to be in captial letters, else it will be ignored
- The json marshall function will return two outputs; a string and error. If there are no errors, it will return nil
  Go to the following site for more details: https://golang.org/pkg/encoding/json/#Marshal
- Marshalling struct -> Converts struct data into string
- Panic immediately kills the server
- Essentially, all errors should be handled accordingly - at the local level (But there are cases where errors need to be handled )
- No native way to declaratively say which field is optional and which field is compulsory. 
*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type helloWorldResponse struct {
	Message string
}

type helloWorldType2Response struct {
	Message string `json:"message"`
}

type helloWorldComplexResponse struct {
	Message string    `json:"message"`
	Author  string    `json:"-"`
	Date    string    `json:",omitempty"`
	Id      int       `json:"id,string"`
	IdMax   int       `json:"idMax,string"`
	Floater float64   `json:"floating"`
	Structure struct {
		Message string    `json:"message"`
		Author  string    `json:"-"`
		Date    string    `json:",omitempty"`
	}
}

func main() {
	port := 8080

	http.HandleFunc("/helloworld", helloWorldHandler)
	http.HandleFunc("/helloworldtype2", helloWorldType2Handler)
	http.HandleFunc("/complex", helloWorldComplexHandler)

	log.Printf("Server starting on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	response := helloWorldResponse{Message: "Hello World"}
	data, err := json.Marshal(response)
	if err != nil {
		panic("Ooops")
	}
	fmt.Fprint(w, string(data))
}

func helloWorldType2Handler(w http.ResponseWriter, r *http.Request) {
	response := helloWorldType2Response{Message: "Hello World The Second"}
	data, err := json.Marshal(response)
	if err != nil {
		panic("Ooops")
	}
	fmt.Fprint(w, string(data))
}

func helloWorldComplexHandler(w http.ResponseWriter, r *http.Request) {
	response := helloWorldComplexResponse{Message: "Hello World The Second", Author: "Arther"}
	data, err := json.Marshal(response)
	if err != nil {
		panic("Ooops")
	}
	fmt.Fprint(w, string(data))
	fmt.Printf(response.Author)
}