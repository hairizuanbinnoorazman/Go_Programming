package main

import (
  "fmt"
  "miao"
)

func main(){
	fmt.Println("cacaca")
	miao.Miao()

	r := miao.Rect{Width: 3, Height: 4}
	c := miao.Circle{Radius: 5}

	miao.Measure(r)
	miao.Measure(c)
	miao.Heyza(r)
	miao.Heyza(c)
	miao.Zazzas()
}
