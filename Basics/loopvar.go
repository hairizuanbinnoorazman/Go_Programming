package main

import (
	"fmt"
	"time"
)

// Basic example for bug - closures in loops
// https://go.dev/blog/loopvar-preview
// Should be fixed go 1.22

func main() {
	for i := 0; i < 10; i++ {
		go func() {
			fmt.Println(i)
		}()
	}
	time.Sleep(5 * time.Second)
}
