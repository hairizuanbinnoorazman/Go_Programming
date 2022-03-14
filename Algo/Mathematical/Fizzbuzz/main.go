package main

import "fmt"

// Problem
// Print numbers but if number is multiple of 3, print Fizz, if multiple of 5, print Buzz, if multiple of 3 and 5, print Fizzbuzz
// Accepts an input
func main() {
	fizzbuzz(25)
}

func fizzbuzz(n int) {
	for i := 1; i <= n; i++ {
		if i%3 == 0 && i%5 == 0 {
			fmt.Println("fizzbuzz")
		} else if i%3 == 0 {
			fmt.Println("fizz")
		} else if i%5 == 0 {
			fmt.Println("buzz")
		} else {
			fmt.Println(i)
		}
	}
}
