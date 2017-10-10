package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	arguments := os.Args
	
	// This is both declaration and initialization step
	infoTag := false
	// Infotag can also be defined the following way:
	// var infoTag bool = false

	for i := 0; i < len(arguments); i++ {
		if strings.Compare(arguments[i], "-i") == 0 {
			infoTag = true
			break
		}
	}

	if infoTag {
		fmt.Println("Info tag enabled!")
		fmt.Print("y/n: ")
		
		// This is just declaration step
		var answer string
		fmt.Scanln(&answer)
		fmt.Println("You entered: ", answer)
	} else {
		fmt.Println("Info tag disabled by default!")
	}
}