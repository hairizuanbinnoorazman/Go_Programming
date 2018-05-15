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

import "fmt"

func partition(arr []int, low, high int) int {
	fmt.Println("Low:", low, "High:", high)

	// Get the pivot to sort based on
	pivotValue := arr[high]

	// The different pointers
	idx := low // idx is to control overall loop
	jdx := low // jdx is to control loop to note the left hand side and right side of pivot

	for idx < high {
		if arr[idx] <= pivotValue {
			temp := arr[idx]
			arr[idx] = arr[jdx]
			arr[jdx] = temp
			jdx = jdx + 1
		}
		idx = idx + 1
	}

	temp := arr[jdx]
	arr[jdx] = arr[high]
	arr[high] = temp
	return jdx
}

func quicksort(arr []int, low, high int) []int {
	if low < high {
		key := partition(arr, low, high)

		arr = quicksort(arr, low, key - 1)
		arr = quicksort(arr, key + 1, high)
	}
	return arr
}


func main(){
	values := []int{10,80, 20, 30, 65, 90, 100, 110, 60, 70}
	fmt.Println(partition(values, 0, 9))
	miao := quicksort(values, 0, 9)
	fmt.Println(miao)
}