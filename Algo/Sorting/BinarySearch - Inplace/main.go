package main

import "fmt"

func main() {
	vals := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} // 5, 9,
	zz := binSearch(vals, 1, 0, len(vals)-1)
	fmt.Println(zz)
}

func binSearch(vals []int, finding, startIdx, endIdx int) bool {
	if startIdx == endIdx {
		if vals[startIdx] == finding {
			return true
		}
		return false
	}

	lhsEndIdx := (startIdx + endIdx) / 2
	rhsStartIdx := (startIdx+endIdx)/2 + 1
	found := false
	if finding <= vals[lhsEndIdx] {
		found = binSearch(vals, finding, startIdx, lhsEndIdx)
	} else {
		found = binSearch(vals, finding, rhsStartIdx, endIdx)
	}
	return found
}
