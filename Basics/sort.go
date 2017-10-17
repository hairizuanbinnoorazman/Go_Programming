/*
	The sort.slice function has the following signature: Slice(slice interface{}, less func(i, j int) bool)
	You need to provide a function which it would then use to do the comparison etc

	append is part of the builtin set of functions - there is no need to append a package name in front of the append function.

	Fyi -> Quick sort kind of slowly breaks down the data structure and sorts each part piece by piece.

*/

package main

import (
	"fmt"
	"sort"
)

type aStructure struct {
	person string
	height int
	weight int
}

func main() {
	mySlice := make([]aStructure, 0)
	a := aStructure{"Mihalis", 180, 90}

	mySlice = append(mySlice, a)
	a = aStructure{"Dimitris", 180, 95}
	mySlice = append(mySlice, a)
	a = aStructure{"Marietta", 155, 45}
	mySlice = append(mySlice, a)
	a = aStructure{"Bill", 134, 40}
	mySlice = append(mySlice, a)

	fmt.Println("0:", mySlice)
	sort.Slice(mySlice, func(i, j int) bool {
			return mySlice[i].weight < mySlice[j].weight
		})
	fmt.Println("<:", mySlice)
	sort.Slice(mySlice, func(i, j int) bool {
			return mySlice[i].weight > mySlice[j].weight
		})
	fmt.Println(">:", mySlice)
}