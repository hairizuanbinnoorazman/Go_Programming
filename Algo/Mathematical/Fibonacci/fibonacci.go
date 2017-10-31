package main

import "fmt"

func fibonacci(n int) int {
	//fmt.Println(n)
	if n > 1 {
		value := fibonacci(n-1) + fibonacci(n-2)
		return value
	}
	if n == 1 {
		return 1
	}
	return 0
}

func main() {
	fmt.Println(fibonacci(10))
}