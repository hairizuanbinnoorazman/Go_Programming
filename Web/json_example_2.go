/*
Basic web server implementation in go
To run this file, run the following command

go run json_example_2.go

Some learning points:
- To declare a variable, it is: `var <variable name> <type>`
- Pointers: https://tour.golang.org/moretypes/1
  If u declare the following: i := 21, then to reference it via pointer, you need to generate/extract the pointer: p := &i
  If u want to reference to an already created pointer, use the * e.g. (Continue from above): m := p  (m is now also a pointer to i)
  
  The reason why you would use pointers is to prevent copying of data in memory, its great for optimizing program performance but very bad
  in terms of readability of code
- For some odd reason, raw encoders & decoders are faster than marshall functions that are provided by the json library
*/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type helloWorldResponse struct {
	Message string `json:"message"`
}

type helloWorldRequest struct {
	Name string `json:"name"`
}

func main() {
	port := 8080

	http.HandleFunc("/helloworld", helloWorldHandler)

	log.Printf("Server starting on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%v", port), nil))
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var request helloWorldRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Obtianing speed gains
	// var request helloWorldRequest
	// decoder := json.NewDecoder(r.Body)
	// err := decoder.Decode(&request)
	// if err != nil {
	// 	http.Error(w, "Bad Request", http.StatusBadRequest)
	// 	return
	// }

	response := helloWorldResponse{Message: "Hello " + request.Name}

	encoder := json.NewEncoder(w)
	encoder.Encode(response)
}
