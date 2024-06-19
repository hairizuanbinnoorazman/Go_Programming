package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func main() {
	fmt.Println("Start server")
	defer fmt.Println("Stop server")
	ln, _ := net.Listen("tcp", ":9999")

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
			fmt.Printf("Obtained input: %v\n", rawInput)
			zzz := strings.Split(rawInput, " ")
			c.name = zzz[0]
			c.inputs = zzz[1:]
			// if inputEntry == 2 {
			// 	c.name = rawInput
			// }
			// if inputEntry == 4 {
			// 	c.listName = rawInput
			// }
			// if inputEntry == 6 {
			// 	c.itemVal = rawInput
			// }
			// if inputEntry == 8 {
			// 	c.itemVal2 = rawInput
			// }
			c.Run()
			inputEntry = inputEntry + 1
		}
	}
}

type Command struct {
	sstore map[string]string
	store  map[string][]string
	conn   net.Conn
	name   string
	inputs []string
	// listName string
	// itemVal  string
	// itemVal2 string

}

func (c *Command) Run() {
	errorStr := "-ERR syntax error\r\n"

	if c.name == "ping" {
		PrintPong(c.conn)
		return
	}
	if c.name == "set" {
		if len(c.inputs) != 2 {
			c.conn.Write([]byte(errorStr))
			return
		}
		c.sstore[c.inputs[0]] = c.inputs[1]
		c.conn.Write([]byte("+OK\r\n"))
	}
	if c.name == "get" {
		if len(c.inputs) != 1 {
			c.conn.Write([]byte(errorStr))
			return
		}
		val := c.sstore[c.inputs[0]]
		processed := fmt.Sprintf("$%v\r\n%v\r\n", len(val), val)
		c.conn.Write([]byte(processed))
		return
	}
	if c.name == "lpush" {
		if len(c.inputs) < 2 {
			c.conn.Write([]byte("-ERR wrong number of arguments for 'lpush' command\r\n"))
			return
		}
		for _, v := range c.inputs[1:] {
			c.store[c.inputs[0]] = append([]string{v}, c.store[c.inputs[0]]...)
		}
		c.conn.Write([]byte(fmt.Sprintf(":%v\r\n", len(c.store[c.inputs[0]]))))
		return
	}
	// if c.name == "lrange" && c.listName != "" && c.itemVal != "" && c.itemVal2 != "" {
	// 	// c.itemVal is "starting value"
	// 	// c.itemVal2 is "ending value" - if more -> it would mean everything
	// 	startIdx, _ := strconv.Atoi(c.itemVal)
	// 	endIdx, _ := strconv.Atoi(c.itemVal2)
	// 	items := c.store[c.listName]
	// 	if startIdx >= len(items) {
	// 		c.conn.Write([]byte("*0\r\n"))
	// 	} else if endIdx < len(items) && endIdx >= 0 {
	// 		items = items[startIdx : endIdx+1]
	// 	} else if endIdx >= len(items) || endIdx < 0 {
	// 		items = items[startIdx:]
	// 	}
	// 	processed := fmt.Sprintf("*%v\r\n", len(items))
	// 	for _, j := range items {
	// 		processed = processed + fmt.Sprintf("$%v\r\n%v\r\n", len(j), j)
	// 	}
	// 	fmt.Println(processed)
	// 	c.conn.Write([]byte(processed))
	// }
	fmt.Printf("-ERR unknown command '%v', with args beginning with %v:\r\n", c.name, c.inputs)
}

func PrintPong(conn net.Conn) {
	conn.Write([]byte("+PONG\r\n"))
}
