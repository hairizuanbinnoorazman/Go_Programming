/*
Basic web server implementation in go.
To run this file run the following command:

go run basic_http_web.go

Some notes to think about
- Difference between Printf and Fprint: Printf prints a bunch of stuff into stdout.
  Fprint prints a bunch of stuff into a file/writer
  Sprintf returns a string
  Printf writes the stuff into stdout and returns bytes return and any write error
- Printf function to handle values can use %v. It replaces it with values. It can take on most values but ->
  It will convert to string first though
- You feed a handler function for the endpoint specified in the go file
*/

package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := 8080

	http.HandleFunc("/", helloWorldHandler)

	log.Printf("Server starting on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("serving", r.URL)
	fmt.Fprint(w, "This is a test. Hello World Miaoza!!\n")
}
