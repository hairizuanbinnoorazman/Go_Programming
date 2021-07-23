package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	port := 8080

	secretFilePath := os.Getenv("SECRET_PATH")
	fmt.Println(secretFilePath)

	ss := SecretFilePrinter{Path: secretFilePath}

	http.HandleFunc("/", helloWorldHandler)
	http.HandleFunc("/secretenv", getEnviromentVar)
	http.Handle("/secretfile", ss)

	log.Printf("Server starting on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

type SecretFilePrinter struct {
	Path string
}

func (s SecretFilePrinter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("start SecretFilePrinter")
	defer log.Println("end SecretFilePrinter")
	raw, err := ioutil.ReadFile(s.Path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to find file"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(raw)
	return
}

func getEnviromentVar(w http.ResponseWriter, r *http.Request) {
	log.Println("start getEnviromentVar")
	defer log.Println("end getEnviromentVar")
	key := r.URL.Query().Get("env")
	if key == "" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("No environment var provided"))
		return
	}
	z := os.Getenv(key)
	if z == "" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Env var: %v - environment var not set", key)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Env var: %v - Value: %v", key, z)))
	return
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("serving", r.URL)
	fmt.Fprint(w, "This is a test. Hello World Miaoza!!\n")
}
