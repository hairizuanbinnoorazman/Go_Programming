package main

import "fmt"

func quicksort(values []int) []int {
	if len(values) <= 1 {
		return values
	}

	leftSide := []int{}
	rightSide := []int{}

	pivot := values[len(values)-1]

	for _, v := range values[0 : len(values)-1] {
		if v < pivot {
			leftSide = append(leftSide, v)
			continue
		}
		rightSide = append(rightSide, v)
	}

	sortedLeftSide := quicksort(leftSide)
	sortedRightSide := quicksort(rightSide)

	sorted := append(sortedLeftSide, pivot)
	sorted = append(sorted, sortedRightSide...)

	return sorted
}

func main() {
	values := []int{5, 3, 4, 1, 2}
	sorted := quicksort(values)
	fmt.Println(sorted)
}
