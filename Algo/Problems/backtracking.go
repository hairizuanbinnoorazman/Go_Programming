/*
This is a basic backtracking algorithm

*/

package main

import "fmt"

func permutate(arr, chosen []string) {
	if len(chosen) == 3 {
		fmt.Println(chosen)
		return
	}
	for idx, val := range arr {
		chosenz := append(chosen, val)
		lol := make([]string, len(arr))
		copy(lol, arr)
		lol = append(lol[:idx], lol[idx+1:]...)
		// fmt.Println(idx, val, chosenz, lol)

		permutate(lol, chosenz)
	}
}

func main() {
	lol := []string{"A", "B", "C", "D"}
	permutate(lol, []string{})
}
