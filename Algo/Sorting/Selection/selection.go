/*
Selection Algorithm

Selection Algorithm is an inplace comparison sort.

Time complexity: O(n2)

Description in plain english:
1. Find the min in array** Important diff between selection and insertion
2. Swap min to 1st place and the latter value to the original min place
3, Move to next element and search the array from that element on for min in array
4. Loop it
 */

package main

// Accepts an array of numbers that is to be sorted accordingly.
func selection_sort(arr []int) []int {
	arrLength :=  len(arr)
	var idx, num int
	for idx, num = range arr {
		// Find the lowest value
		lowestNum := num
		lowestIdx := idx
		var iterNum, iterIdx int
		for iterIdx, iterNum = range arr[idx:arrLength] {
			if lowestNum > iterNum {
				lowestNum = iterNum
				lowestIdx = iterIdx + idx
			}
		}
		arr[lowestIdx] = num
		arr[idx] = lowestNum
	}
	return arr
}



func main(){
	trialArr := []int{-1, 7, 1, 2, -8, 0,0, 5, 4, 3}
	selection_sort(trialArr)
}

