package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/rs/cors"
)

var hashKey = []byte("very-secret")
var blockKey = []byte("a-lot-secret")
var s = securecookie.New(hashKey, blockKey)

type HomeHandler struct{}

func (h HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	value := map[string]string{
		"foo": "bar",
	}
	if encoded, err := s.Encode("cookie-name", value); err == nil {
		cookie := &http.Cookie{
			Name:     "cookie-name",
			Value:    encoded,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
		}
		log.Printf("Cookie Generated :: %v", encoded)
		http.SetCookie(w, cookie)
	}
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

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"*"},
	})

	srv := &http.Server{
		Handler: c.Handler(r),
		Addr:    "0.0.0.0:8080",
	}

	log.Fatal(srv.ListenAndServe())
}
