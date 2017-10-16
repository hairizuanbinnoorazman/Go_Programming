package main

import (
	"fmt"
)


func messageParser(num int, lol string) string {
	miao := fmt.Sprintf("%d%v", num, lol)
	return miao
}


func main() {
	fmt.Println("This is an example")
	fmt.Println(messageParser(12, "caca"))
}