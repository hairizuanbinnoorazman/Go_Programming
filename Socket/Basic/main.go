// Reference from: https://gist.github.com/hakobe/6f70d69b8c5243117787fd488ae7fbf2
package main

import (
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var SocketFile = "/tmp/go.sock"

func echoServer(c net.Conn) {
	for {
		buf := make([]byte, 512)
		nr, err := c.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Println("END OF FILE")
				return
			}
			log.Println("error in trying to read data")
			return
		}

		data := buf[0:nr]
		println("Server got:", string(data))
		_, err = c.Write(data)
		if err != nil {
			log.Fatal("Writing client error: ", err)
		}
	}
}

func main() {
	log.Println("Starting echo server")
	ln, err := net.Listen("unix", SocketFile)
	if err != nil {
		log.Fatal("Listen error: ", err)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	go func(ln net.Listener, c chan os.Signal) {
		sig := <-c
		log.Printf("Caught signal %s: shutting down.", sig)
		ln.Close()
		os.Exit(0)
	}(ln, sigc)

	for {
		fd, err := ln.Accept()
		if err != nil {
			log.Fatal("Accept error: ", err)
		}

		go echoServer(fd)
	}
}
