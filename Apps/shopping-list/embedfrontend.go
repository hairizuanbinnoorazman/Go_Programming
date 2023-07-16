//go:build embedfrontend

package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
)

// Reference: https://github.com/golang/go/issues/44484
//
//go:embed static/*
var frontendFolder embed.FS

func addFrontendRoutes(r *mux.Router) {
	folderRoot, err := fs.Sub(frontendFolder, "static")
	if err != nil {
		panic(fmt.Sprintf("Not handled correctly %v", err))
	}

	ingressPath := "/tools/shopping-list"

	// r.PathPrefix("/").Handler(http.FileServer(http.Dir("/root/frontend")))
	r.Handle(ingressPath+"/", http.StripPrefix(ingressPath, http.FileServer(http.FS(folderRoot))))
	r.Handle(ingressPath+"/main.js", http.StripPrefix(ingressPath, http.FileServer(http.FS(folderRoot))))
	r.Handle(ingressPath+"/favicon.ico", http.StripPrefix(ingressPath, http.FileServer(http.FS(folderRoot))))
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
	tmpl, err := template.ParseFS(frontendFolder, "static/*")
	if err != nil {
		h.Logger.Errorf("Error: %v", err)
		return
	}
	err = tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		h.Logger.Errorf("Error: %v", err)
		return
	}
}
