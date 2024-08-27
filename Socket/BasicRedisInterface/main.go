package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var enableHelloCommand bool = false
var enableBrokenUpCommands bool = false

func main() {
	fmt.Println("Start server")
	defer fmt.Println("Stop server")

	_, enableHelloCommand = os.LookupEnv("ENABLE_HELLO_COMMAND")
	_, enableBrokenUpCommands = os.LookupEnv("ENABLE_BROKEN_UP_COMMANDS")

	ln, _ := net.Listen("tcp", ":9999")

	sstore := map[string]string{}
	store := map[string][]string{}

	for {
		conn, _ := ln.Accept()
		scanner := bufio.NewScanner(conn)
		c := Command{sstore: sstore, store: store, conn: conn}
		// inputEntry := 0
		expectedBrokenUpItems := 0
		zzz := []string{}
		for {
			if ok := scanner.Scan(); !ok {
				break
			}
			rawInput := scanner.Text()
			fmt.Printf("Obtained input: %v\n", rawInput)

			if enableBrokenUpCommands {
				if string(rawInput[0]) == "*" {
					fmt.Println("Broken up mode detected...")
					processedInput := strings.Replace(rawInput, "*", "", -1)
					expectedBrokenUpItems, _ = strconv.Atoi(processedInput)
					continue
				}

				if expectedBrokenUpItems > 0 {
					fmt.Println("In broken up command mode")
					if string(rawInput[0]) == "$" {
						fmt.Println("we're processing this crap")
						continue
					}

					zzz = append(zzz, rawInput)

					if len(zzz) < expectedBrokenUpItems {
						fmt.Println("Insufficient values for broken up input")
						continue
					}

					expectedBrokenUpItems = 0
				}
			}

			if len(zzz) == 0 {
				zzz = strings.Split(rawInput, " ")
			}
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
			zzz = []string{}
			// inputEntry = inputEntry + 1
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
		c.conn.Write([]byte("+PONG\r\n"))
		return
	}
	if c.name == "set" {
		if len(c.inputs) < 2 {
			c.conn.Write([]byte(errorStr))
			return
		}
		if len(c.inputs) == 3 {
			fmt.Println("we will ignore time expiry for now")
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

	if enableHelloCommand {
		if c.name == "hello" {
			if len(c.inputs) != 1 {
				c.conn.Write([]byte("-ERR"))
				return
			}
			items := [][]string{
				{"server", "redis"},
				{"version", "7.2.5"},
				{"proto", ":3"},
				{"id", ":13"},
				{"mode", "standalone"},
				{"role", "master"},
				{"modules", "*0"},
			}
			c.conn.Write([]byte(fmt.Sprintf("%%%v\r\n", len(items))))
			for _, aaa := range items {
				k := aaa[0]
				v := aaa[1]
				c.conn.Write([]byte(fmt.Sprintf("$%v\r\n", len(k))))
				c.conn.Write([]byte(k + "\r\n"))
				if k == "proto" || k == "id" || k == "modules" {
					c.conn.Write([]byte(v + "\r\n"))
				} else {
					c.conn.Write([]byte(fmt.Sprintf("$%v\r\n", len(v))))
					c.conn.Write([]byte(v + "\r\n"))
				}
			}
			return

			// 		%7
			// $6
			// server
			// $5
			// redis
			// $7
			// version
			// $5
			// 7.2.5
			// $5
			// proto
			// :3
			// $2
			// id
			// :13
			// $4
			// mode
			// $10
			// standalone
			// $4
			// role
			// $6
			// master
			// $7
			// modules
			// *0
		}
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

	// We can return or not return this...
	// responseText := fmt.Sprintf("-ERR unknown command '%v', with args beginning with %v:\r\n", c.name, c.inputs)
	// c.conn.Write([]byte(responseText))
}
