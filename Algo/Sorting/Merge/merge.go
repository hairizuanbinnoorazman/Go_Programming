package main

import "fmt"

/*
Merge sort:
- Split array into two until it cannot be split anymore
- While combining it, put them back in the right order
*/

func MergeSort(items []int) []int {
	if len(items) <= 1 {
		return items
	}

	leftSide := MergeSort(items[0 : len(items)/2])
	rightSide := MergeSort(items[len(items)/2 : len(items)])

	i := 0
	j := 0
	combined := []int{}
	for i < len(leftSide) || j < len(rightSide) {
		if i >= len(leftSide) {
			combined = append(combined, rightSide[j:]...)
			j = len(rightSide)
			continue
		}
		if j >= len(rightSide) {
			combined = append(combined, leftSide[i:]...)
			i = len(leftSide)
			continue
		}
		if leftSide[i] < rightSide[j] {
			combined = append(combined, leftSide[i])
			i = i + 1
			continue
		}
		combined = append(combined, rightSide[j])
		j = j + 1
	}

	return combined
}

func main() {
	values := []int{4, 3, 2, 1}
	sorted := MergeSort(values)
	fmt.Println(sorted)
}
