package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func random(min, max int) int {
	return rand.Intn(max - min) + min
}

func main() {
	Min := 0
	Max := 0
	Total := 0
	if len(os.Args) > 3 {
		Min, _ = strconv.Atoi(os.Args[1])
		Max, _ = strconv.Atoi(os.Args[2])
		Total, _ = strconv.Atoi(os.Args[3])
	} else {
		// Interesting point - os.Args[0] actually mean the name of the program
		// Read the documentation - it only says that
		fmt.Println("Usage:", os.Args[0], "MIN MAX TOTAL")
		os.Exit(-1)
	}

	rand.Seed(time.Now().Unix())
	for i := 0; i < Total; i++ {
		myRand := random(Min, Max)
		fmt.Print(myRand)
		fmt.Print(" ")
	}
	fmt.Println()
}

