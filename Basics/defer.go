/*
Go file to understand defer

To run:
go run defer.go

Output:
2 1 0
3 3 3
2 1 0

Fyi:
Add another quick test to see how many defer statements can be put into a function. Multiple defer statements can be added
but it goes in reverse order -> Like in a stack.
*/

package main

import (
	"fmt"
)

func a1() {
	for i := 0; i < 3; i++ {
		defer fmt.Print(i, " ")
	}
}

// Only runs the defer function after the for loop
// i would have become 3 but the defer function has been triggered 3 times by now
// It fetches the current version of i
func a2() {
	for i := 0; i < 3; i++ {
		defer func() { fmt.Print(i, " ") }()
	}
}

func a3() {
	for i := 0; i < 3; i++ {
		defer func(n int){ fmt.Print(n, " ") }(i)
	}
}

func a4() {
	defer fmt.Println("a4 test 1")
	defer fmt.Println("a4 test 2")
	defer fmt.Println("a4 test 3")
	defer fmt.Println("a4 test 4")
}

func main(){
	a1()
	fmt.Println()
	a2()
	fmt.Println()
	a3()
	fmt.Println()
	a4()
}