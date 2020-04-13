// Reference from: https://gist.github.com/hakobe/6f70d69b8c5243117787fd488ae7fbf2
package main

import (
	"io"
	"log"
	"net"
	"time"
)

func reader(r io.Reader) {
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			return
		}
		println("Client got:", string(buf[0:n]))
	}
}

func main() {
	c, err := net.Dial("unix", "/tmp/go.sock")
	if err != nil {
		log.Fatal("Dial error", err)
	}
	defer c.Close()

	go reader(c)
	for {
		msg := "hi"
		_, err := c.Write([]byte(msg))
		if err != nil {
			log.Fatal("Write error:", err)
			break
		}
		println("Client sent:", msg)
		time.Sleep(1 * time.Second)
	}
}
