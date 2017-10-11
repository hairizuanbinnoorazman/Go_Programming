/*
Arrays are usually manipulated via reference
However, if passed to reference, a copy will be created for the function - making it slow
Array sizes cannot be manipulated
*/

package main

import "fmt"

func main() {
	myArray := [4]int{1, 2, 3, -4}
	threeD := [2][2][2]int{{{1, 2}, {3, 4}}, {{5, 6}, {7, 8}}}
	twoD := [3][3]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	fmt.Println(myArray)
	fmt.Println(threeD)
	fmt.Println(myArray[0])
	fmt.Println(twoD[2][1])
	fmt.Println(threeD[0][1][1])
}