package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

type zzz struct{}

func (z zzz) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	type miao struct {
		Name    string `json:"name"`
		Address string `json:"address"`
		Date    string `json:"date"`
	}
	mm := miao{
		Name:    "zzz",
		Address: "yyy",
		Date:    "mmm",
	}
	raw, _ := json.Marshal(mm)
	w.WriteHeader(http.StatusOK)
	w.Write(raw)
}

type miao struct{}

func (m miao) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tmplFile := "miao.tmpl"
	tmpl, err := template.New(tmplFile).ParseFiles(tmplFile)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		panic(err)
	}
}

type yy struct{}

func (y yy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tmplFile := "yyy.tmpl"
	tmpl, err := template.New(tmplFile).ParseFiles(tmplFile)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		panic(err)
	}
}

func main() {
	http.Handle("/api/v1/zzz", zzz{})
	http.Handle("/miao", miao{})
	http.Handle("/yyy", yy{})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
