/*
This example is for understanding how to do convert and clean data input on
command line utilities
*/

package main

import (
  "fmt"
  "os"
  "strconv"
)

func main(){
  arguments := os.Args
  sum := 0
  for i := 0; i < len(arguments); i++ {
    temp, _ := strconv.Atoi(arguments[i])
    sum = sum + temp
  }
  fmt.Println("Sum:", sum)
}
