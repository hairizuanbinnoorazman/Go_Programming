package main

import "fmt"

func BinarySearch(finding int, values []int) bool {
	if len(values) == 0 {
		return false
	}
	if len(values) == 1 {
		if values[0] == finding {
			return true
		}
		return false
	}

	found := false
	leftHalf := values[0 : len(values)/2]
	rightHalf := values[len(values)/2 : len(values)]

	if finding >= rightHalf[0] {
		found = BinarySearch(finding, rightHalf)
	} else {
		found = BinarySearch(finding, leftHalf)
	}

	return found
}

func main() {
	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	found := BinarySearch(8, items)
	fmt.Printf("Was 8 found: %v\n", found)
}
