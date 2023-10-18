/*
For practise - need to be able to remove and add elements to an array in place with no issue

Problems:
1. Given an array - remove an item from it via index, and then move the rest of the elements down
2. Given an array - move one number to the right hand side of the array
*/
package main

import "fmt"

func main() {
	items := []int{1, 2, 3, 4, 5}
	removeElements(items, 3, 5)
	fmt.Println(items)
	removeElements(items, 2, 4)
	fmt.Println(items)
	moveVal(items, 2)
	fmt.Println(items)
}

func removeElements(a []int, idx, actualLen int) {
	for i := idx; i < actualLen-1; i++ {
		a[idx] = a[idx+1]
	}
	a[actualLen-1] = -1
}

func moveVal(a []int, val int) {
	sortHappened := true
	for sortHappened {
		sortHappened = false
		for k, _ := range a {
			if k == len(a)-1 {
				continue
			}
			if a[k] == val {
				a[k] = a[k+1]
				a[k+1] = val
				sortHappened = true
			}
		}
	}
}

func moveLeft(aa []int, idx int) {
	counter := idx
	for i := len(aa) - 1 - idx; i > 0; i-- {
		aa[counter] = aa[counter+1]
		counter += 1
	}
}

func moveRight(aa []int, idx int) {
	previousVal := 0
	nextVal := aa[idx]
	for i := idx; i < len(aa); i++ {
		aa[i] = previousVal
		previousVal = nextVal
		if i+1 < len(aa) {
			nextVal = aa[i+1]
		}
	}
}
