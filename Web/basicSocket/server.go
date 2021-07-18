package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func main() {
	log.Println("Start Server")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/v1/ws", func(w http.ResponseWriter, r *http.Request) {
		var conn, _ = upgrader.Upgrade(w, r, nil)
		go func(conn *websocket.Conn) {
			for {
				mType, msg, err := conn.ReadMessage()
				if err != nil {
					log.Println(err.Error())
				}

				log.Println(string(msg))
				conn.WriteMessage(mType, msg)
			}
		}(conn)
	})

	http.HandleFunc("/v2/ws", func(w http.ResponseWriter, r *http.Request) {
		var conn, _ = upgrader.Upgrade(w, r, nil)
		go func(conn *websocket.Conn) {
			ch := time.Tick(time.Second)

			for range ch {
				conn.WriteJSON(myStruct{Username: "cnkamclacmajnejra"})
			}
		}(conn)
	})

	log.Fatal(http.ListenAndServe(":3000", nil))
}

type myStruct struct {
	Username string `json:"username"`
}
