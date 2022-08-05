package main

import "fmt"

var store = map[int]int{
	1: 1,
	2: 1,
}

func fibonacci(n int) int {
	if n <= 0 {
		return 0
	}
	if store[n] != 0 {
		return store[n]
	}
	val := fibonacci(n-1) + fibonacci(n-2)
	store[n] = val
	return val
}

func fibonacciTabulate(n int) int {
	if n <= 0 {
		return 0
	}
	if n <= 2 {
		return 1
	}
	previous1 := 1
	previous2 := 1
	currentVal := 0
	for i := 3; i <= n; i++ {
		currentVal = previous1 + previous2
		previous1 = previous2
		previous2 = currentVal
	}
	return currentVal
}

// 0, 1, 1, 2, 3, 5, 8, 13, 21,...
func main() {
	fmt.Println(fibonacci(100))
	fmt.Println(fibonacciTabulate(100))
}
