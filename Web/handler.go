package main

import ( 
	"fmt"
	"encoding/json"
	"log"
	"net/http"
)


type helloWorldResponse struct {
	Message string `json:"message"`
}


type helloWorldRequest struct {
	Message string `json:"message"`
}


type helloWorldHandler struct {}

func (h helloWorldHandler) ServeHTTP (rw http.ResponseWriter, r *http.Request) {
	response := helloWorldResponse{Message: "Hello"}

	encoder := json.NewEncoder(rw)
	encoder.Encode(response)
}

func newHelloWorldHandler() http.Handler {
	return helloWorldHandler{}
}



type validationHandler struct {
	next http.Handler
}

func (h validationHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request){
	h.next.ServeHTTP(rw, r)
}

func newValidationHandler(next http.Handler) http.Handler {
	return validationHandler{next: next}
}

func main() {
	port := 8080

	handler := newValidationHandler(newHelloWorldHandler())
	http.Handle("/helloworld", handler)

	log.Printf("Server starting on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}