/*
	Aim of this file is to demonstrate the use of templates

	Learnings:
	1. In order to use values in templates and assign them to certain values, they need to be exported state
	2. Plenty of tutorials just put plain . but if you put .<Variable name> -> That allows you to link accordidngly
*/

package main

import (
	"log"
	"net/http"
	"html/template"
)

type adminHTML struct {
	Miao string
	Heh string
}


func admin(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving admin")
	t, err := template.New("admin.html").ParseFiles("admin.html")
	if err != nil {
		log.Println(err.Error())
	}

	err = t.Execute(w, adminHTML{Miao:"cac", Heh:"jkacnkanc"})
	if err != nil {
		log.Println("Error:", err.Error())
	}
}


func main() {
	server := http.Server{
        Addr: "127.0.0.1:8080",
    }
	
	log.Println("Serving on", server.Addr)
	http.HandleFunc("/admin", admin)
	log.Fatal(server.ListenAndServe())
}