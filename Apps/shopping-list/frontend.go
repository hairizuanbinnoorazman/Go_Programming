//go:build !embedfrontend

package main

import (
	"net/http"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
)

func addFrontendRoutes(r *mux.Router) {
	// r.PathPrefix("/").Handler(http.FileServer(http.Dir("/root/frontend")))
	r.Handle("/", http.FileServer(http.Dir("/root/frontend")))
	r.Handle("/main.js", http.FileServer(http.Dir("/root/frontend")))
	r.Handle("/favicon.ico", http.FileServer(http.Dir("/root/frontend")))
}

type NotFound struct {
	Logger Logger
}

func (h NotFound) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Logger.Infof("Not found path: %v", r.URL.Path)
	if strings.HasPrefix(r.URL.Path, "/api/v1/") {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("api route not found"))
		return
	}
	// http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	tmpl, err := template.ParseFiles("/root/frontend/index.html")
	if err != nil {
		h.Logger.Errorf("Error: %v", err)
	}
	err = tmpl.Execute(w, "a")
	if err != nil {
		h.Logger.Errorf("Error: %v", err)
	}
}
