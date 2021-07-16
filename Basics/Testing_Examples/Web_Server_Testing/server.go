package main

import (
	"fmt"
	"log"
	"net/http"
)

type helloWorldHandler struct{}

func (h helloWorldHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving hello world handler")
	fmt.Fprint(w, "Miao")
}

func main() {
	port := 8080

	http.Handle("/admin", helloWorldHandler{})

	log.Printf("Server starting on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}
