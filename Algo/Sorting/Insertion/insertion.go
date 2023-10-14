/*
Insertion Sort

Description in English:
1. There is a unsorted/sorted portion in the array
2. After first card, see second card.
3. On second card, swap if first card is bigger
4. If it was a smaller number, keep swapping left

Refer to the following page for details:
https://www.geeksforgeeks.org/insertion-sort/
*/

package main

import "fmt"

func insertion_sort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}
	for i := 1; i < len(arr); i++ {
		reversingIdx := i
		for arr[reversingIdx] < arr[reversingIdx-1] {
			tempVal := arr[reversingIdx-1]
			arr[reversingIdx-1] = arr[reversingIdx]
			arr[reversingIdx] = tempVal
			reversingIdx = reversingIdx - 1
			if reversingIdx == 0 {
				break
			}
		}
	}

	return arr
}

func main() {
	values := []int{2, 3, 6, 1, 0, -1, 8, -5}
	miao := insertion_sort(values)
	fmt.Println(miao)
}
