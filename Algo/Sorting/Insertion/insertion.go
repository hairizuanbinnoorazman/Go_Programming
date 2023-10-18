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

func insertion_sort(arr []int) {
	for k := 1; k < len(arr); k++ {
		for m := k; m > 0; m-- {
			if arr[m] < arr[m-1] {
				tempVal := arr[m-1]
				arr[m-1] = arr[m]
				arr[m] = tempVal
			} else {
				break
			}
		}
	}
}

func main() {
	values := []int{2, 3, 6, 1, 0, -1, 8, -5}
	insertion_sort(values)
	fmt.Println(values)
}
