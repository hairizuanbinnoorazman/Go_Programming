package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	for {
		fmt.Println("Start")
		for _, e := range os.Environ() {
			fmt.Println(e)
		}
		fmt.Printf("All environment variables printed\n\n\n")
		time.Sleep(5 * time.Second)
	}

}
