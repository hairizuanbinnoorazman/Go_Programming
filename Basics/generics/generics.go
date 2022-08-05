package main

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

func min[T constraints.Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}

func tester[T constraints.Ordered](x int, y T) T {
	if x < 1 {
		return y
	}
	return x
}

func main() {
	fmt.Println("start")
	fmt.Println(min[int](12, 13))
}
