package main

import (
	"fmt"
)

func change(x []int) {
	x[3] = -2
}

func print_array(x []int) {
	for _, number := range x {
		fmt.Printf("%d", number)
	}
	fmt.Println("")
}

func main() {
	aSlice := []int{1, 2, 3, 4, 5, 6, 7}
	fmt.Println("Before Change:")
	print_array(aSlice)
	fmt.Println("After Change:")
	change(aSlice)
	print_array(aSlice)
}