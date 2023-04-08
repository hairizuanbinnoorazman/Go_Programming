package main

import (
	"fmt"
	"strings"

	cooklang "github.com/justintout/cooklang-go"
)

func main() {
	fmt.Println("Begin golang code")

	zzz := cooklang.MustParseFile("zzz.cook")
	// fmt.Printf("Test %+v\n", zzz)
	// fmt.Printf("%+v", zzz.Ingredients)
	for i, _ := range zzz.Ingredients {
		fmt.Println(i)
	}

	for j, k := range zzz.Metadata {
		if j == "tags" {
			zz := fmt.Sprintf("%v", j)
			for _, a := range strings.Split(k, ",") {
				zz = fmt.Sprintf("%v == %v", zz, a)
			}
			fmt.Println(zz)
		}
	}
}
