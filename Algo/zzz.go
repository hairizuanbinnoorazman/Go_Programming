package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	zzz := bufio.NewReader(os.Stdin)
	aaa, _ := zzz.ReadString('\n')
	fmt.Println(aaa)

	aaa = strings.TrimSpace(aaa)
	aa := strings.Split(aaa, " ")
	fmt.Println(aa)

	lol := os.Args
	fmt.Println(lol[1])
	fmt.Println(lol[2])
}
