/*
Quicksort algorithm

There are two parts to the quick sort algorithm
1. Sorting based on a pivot (There are several ways to pivot the sort:
   a. Last Element
   b. First Element
   c. Random Element
   d. Median Element
2. Partition based off the pivot and repeat

This algorithm is a in place algorithm, which is why you would see that it plays quite a lot with pointers.
*However, tweaks need to be done in order to convert this to a proper in place algorithm. Using pointers and references within the code.
*/

package main

import (
	"fmt"
)

// returns pivot split point
func partition(arr []int, low, high int) int {
	pivotNum := arr[high]
	initialHigh := high
	high = high - 1
	for low < high {
		// Find for "value that needs to go right"
		for arr[low] < pivotNum && low <= high {
			low = low + 1
		}
		if low >= high {
			break
		}
		// FInd for value that needs to go left
		if arr[high] > pivotNum && low <= high {
			high = high - 1
		}
		if low >= high {
			break
		}
		temp := arr[low]
		arr[low] = arr[high]
		arr[high] = temp
	}
	if arr[initialHigh] < arr[high] {
		temp := arr[high]
		arr[high] = pivotNum
		arr[initialHigh] = temp
	}
	return low
}

func quicksort(arr []int, low, high int) {
	if low < high {
		splitCounter := partition(arr, low, high)
		quicksort(arr, low, splitCounter-1)
		quicksort(arr, splitCounter+1, high)
	}
	return
}

func main() {
	values := []int{10, 80, 20, 30, 65, 90, 100, 110, 60, 70}
	quicksort(values, 0, 9)
	fmt.Println(values)
}
