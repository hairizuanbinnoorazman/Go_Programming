// this is just a test to see if the compiler would complain about the file in this folder.
// apparently go compiler ignores files that has test in its name? -  It needs to end with test
// This file will not be picked up at all

package miao_test

import "fmt"

func Goddess() {
	fmt.Println("There is a goodness in everyone")
}