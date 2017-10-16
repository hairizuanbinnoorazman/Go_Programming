package main

import (
	"fmt"
)

func main() {
	// Need make to create the actual item in memory
	// Map is a hash map - almost like your dictionary in python

	// This is just the declaration step here
	aMap := make(map[string]int) 

	aMap["Mon"] = 0 
	aMap["Tue"] = 1 
	aMap["Wed"] = 2 
	aMap["Thu"] = 3 
	aMap["Fri"] = 4 
	aMap["Sat"] = 5 
	aMap["Sun"] = 6 

	fmt.Printf("Sunday is the %dth day of the week.\n", aMap["Sun"]) 

	count := 0
	for key, value := range aMap {
		count++
		fmt.Println("Value of key:", key)
		fmt.Println("Value of value:", value)
	}
	fmt.Printf("There are %d elements in aMap\n", count)

	// It is possible to declare and initialize the variable here
	anotherMap := map[string]int{
		"here": 12,
		"cmac": 145,
	}
	fmt.Println("Value of item here in anotherMap map: ", anotherMap["here"])
}