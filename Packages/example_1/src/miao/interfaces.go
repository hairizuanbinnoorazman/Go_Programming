package miao

import (
	"fmt"
	"math"
)

type geometry interface {
	area() float64
	perim() float64
}

type Rect struct {
	Width, Height float64
}

type Circle struct {
	Radius float64
}

func (r Rect) area() float64 {
	return r.Width * r.Height
}

func (r Rect) perim() float64 {
	return 2*r.Width + 2*r.Height
}

func (c Circle) area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) perim() float64 {
	return 2 * math.Pi * c.Radius
}

func (c Circle) mioza() float64 {
	return c.area() + c.perim()
}

func Measure(g geometry) {
	fmt.Println(g)
	fmt.Println(g.area())
	fmt.Println(g.perim())
}
