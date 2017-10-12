/*
Go routines

To run the following script, run the following command:
go run goroutines.go
*/

package main

import "fmt"

func f(s string) {
	for i := 0; i < 3; i++ {
		fmt.Println(s, ":", i)
	}
}

func main() {
	f("sync")

	go f("goroutine mode")

	go func(msg string) {fmt.Println(msg)}("going")
	go func(msg string) {fmt.Println(msg)}("going2")
	go func(msg string) {fmt.Println(msg)}("going3")

	f("miao")

	var input string
	fmt.Scanln(&input)
	fmt.Println("done")
}
