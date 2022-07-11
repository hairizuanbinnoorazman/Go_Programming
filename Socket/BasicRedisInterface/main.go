package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
)

func main() {
	fmt.Println("Start server")
	defer fmt.Println("Stop server")
	ln, _ := net.Listen("tcp", ":6379")

	sstore := map[string]string{}
	store := map[string][]string{}

	for {
		conn, _ := ln.Accept()
		scanner := bufio.NewScanner(conn)
		c := Command{sstore: sstore, store: store, conn: conn}
		inputEntry := 0
		for {
			if ok := scanner.Scan(); !ok {
				break
			}
			rawInput := scanner.Text()
			fmt.Println(rawInput)
			if inputEntry == 2 {
				c.name = rawInput
			}
			if inputEntry == 4 {
				c.listName = rawInput
			}
			if inputEntry == 6 {
				c.itemVal = rawInput
			}
			if inputEntry == 8 {
				c.itemVal2 = rawInput
			}
			c.Run()
			inputEntry = inputEntry + 1
		}
	}
}

type Command struct {
	sstore   map[string]string
	store    map[string][]string
	conn     net.Conn
	name     string
	listName string
	itemVal  string
	itemVal2 string
}

func (c *Command) Run() {
	if c.name == "ping" {
		PrintPong(c.conn)
	}
	if c.name == "set" && c.listName != "" && c.itemVal != "" {
		c.sstore[c.listName] = c.itemVal
		c.conn.Write([]byte("+OK\r\n"))
	}
	if c.name == "get" && c.listName != "" {
		val := c.sstore[c.listName]
		processed := fmt.Sprintf("$%v\r\n%v\r\n", len(val), val)
		c.conn.Write([]byte(processed))
	}
	if c.name == "lpush" && c.listName != "" && c.itemVal != "" {
		c.store[c.listName] = append([]string{c.itemVal}, c.store[c.listName]...)
		c.conn.Write([]byte(fmt.Sprintf(":%v\r\n", len(c.store[c.listName]))))
	}
	if c.name == "lrange" && c.listName != "" && c.itemVal != "" && c.itemVal2 != "" {
		// c.itemVal is "starting value"
		// c.itemVal2 is "ending value" - if more -> it would mean everything
		startIdx, _ := strconv.Atoi(c.itemVal)
		endIdx, _ := strconv.Atoi(c.itemVal2)
		items := c.store[c.listName]
		if startIdx >= len(items) {
			c.conn.Write([]byte("*0\r\n"))
		} else if endIdx < len(items) && endIdx >= 0 {
			items = items[startIdx : endIdx+1]
		} else if endIdx >= len(items) || endIdx < 0 {
			items = items[startIdx:]
		}
		processed := fmt.Sprintf("*%v\r\n", len(items))
		for _, j := range items {
			processed = processed + fmt.Sprintf("$%v\r\n%v\r\n", len(j), j)
		}
		fmt.Println(processed)
		c.conn.Write([]byte(processed))
	}
}

func PrintPong(conn net.Conn) {
	conn.Write([]byte("+PONG\r\n"))
}
