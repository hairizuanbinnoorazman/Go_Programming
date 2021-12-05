package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type HomeHandler struct{}

func (h HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Home Handler endpoint reached")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

type Websocket struct {
	hub *Hub
}

func (a Websocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: a.hub, conn: conn, send: make(chan []byte, 256)}
	a.hub.register <- client

	go client.readPump()
	go client.writePump()
}

func main() {
	log.Println("Application Start")
	h := newHub()
	go h.run()
	r := mux.NewRouter()
	r.Handle("/", HomeHandler{})
	r.Handle("/ws", Websocket{hub: h})
	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:8080",
	}

	log.Fatal(srv.ListenAndServe())
}
