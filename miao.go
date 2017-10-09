package main

import "fmt"
import "os"

func main() {
	arguments := os.Args
	for i := 0; i < len(arguments); i++ {
		fmt.Println(arguments[i])
	}
}
