/*
Insertion Sort

Description in English:
1. There is a unsorted/sorted portion in the array
2. After first card, see second card.
3. On second card, swap if first card is bigger
4. If it was a smaller number, keep swapping left
 */

package main

import "fmt"

func insertion_sort(arr []int) []int {
	var idx int
	for idx, _ = range arr {
		iter := idx
		for iter > 0 {
			if arr[iter] < arr[iter-1] {
				temp := arr[iter]

				arr[iter] = arr[iter-1]
				arr[iter-1] = temp
			}
			iter = iter - 1
			fmt.Println(iter)
		}
	}
	return arr
}

func main() {
	values := []int{2, 3, 6, 1, 0, -1, 8, -5}
	miao := insertion_sort(values)
	fmt.Println(miao)
}
