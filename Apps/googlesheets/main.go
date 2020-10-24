package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	drive "google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	sheets "google.golang.org/api/sheets/v4"
)

func main() {
	// Added comment
	fmt.Println("Begin server")

	svc, err := sheets.NewService(context.Background(), option.WithCredentialsFile("key.json"), option.WithScopes("https://www.googleapis.com/auth/spreadsheets"))
	if err != nil {
		fmt.Println(err.Error())
	}
	dsvc, err := drive.NewService(context.Background(), option.WithCredentialsFile("key.json"), option.WithScopes("https://www.googleapis.com/auth/drive.file"))
	if err != nil {
		fmt.Println(err.Error())
	}

	http.Handle("/", ssHandler{sSvc: svc, dSvc: dsvc})
	log.Fatal(http.ListenAndServe(":9000", nil))
}

type ssHandler struct {
	sSvc *sheets.Service
	dSvc *drive.Service
}

func (s ssHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Begin serving request to create new spreadsheet")
	defer log.Println("End serving request to create new spreadsheet")

	call := s.sSvc.Spreadsheets.Create(&sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{
			Title: "Lol",
		},
	})

	ss, err := call.Do()
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Errored out: %v", err.Error())))
		return
	}

	log.Println(ss.SpreadsheetUrl)
	log.Println(ss.SpreadsheetId)
	log.Println(ss.Properties.Title)

	aa := s.dSvc.Permissions.Create(ss.SpreadsheetId, &drive.Permission{
		Type:         "user",
		Role:         "writer",
		EmailAddress: "hairizuanbinnoorazman@gmail.com",
	})
	_, err = aa.Do()
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Errored out: %v", err.Error())))
		return
	}

	bb := s.dSvc.Files.List()
	resp, err := bb.Do()
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Errored out: %v", err.Error())))
		return
	}
	for key, x := range resp.Files {
		log.Println(key)
		log.Println(x.Id)
		log.Println(x.Name)
	}

	w.WriteHeader(200)
	w.Write([]byte("New sspreadsheet created"))
	return
}
