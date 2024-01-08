package main

import (
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

type signinPage struct{}

func (b signinPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("started signin-page handler")
	defer log.Println("ended signin-page handler")

	tmpl := template.Must(template.ParseFiles("layout.html"))
	tmpl.Execute(w, nil)
}

type signin struct{}

func (b signin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("started signin handler")
	defer log.Println("ended signin handler")

	name := r.FormValue("name")
	password := r.FormValue("password")
	if name == "admin" && password == "password" {
		cookie := http.Cookie{
			Name:  "test",
			Value: "test-cookie",
			Path:  "/",
		}
		http.SetCookie(w, &cookie)
		w.Write([]byte("successfully login"))
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("unauthorized login"))
}

type auth struct{}

func (a auth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("test")
	if err == nil {
		log.Println("cookie found, will return 200 ok")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("cookie found - successfully in"))
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("invalid"))
}

func main() {
	log.Print("Auth started")

	r := mux.NewRouter()
	r.Handle("/", signinPage{})
	r.Handle("/signin", signin{})
	r.Handle("/auth", auth{})
	srv := http.Server{
		Handler: r,
		Addr:    "0.0.0.0:8080",
	}
	log.Fatal(srv.ListenAndServe())
}
