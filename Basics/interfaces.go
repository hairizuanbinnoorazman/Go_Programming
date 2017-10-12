/*
Learning Go interfaces

In order to run and try the following file, run the following command:
go run interfaces.

Learnings:
- There are various approaches to establish relationships between objects
  http://spf13.com/post/is-go-object-oriented/
  - single/multiple inheritance
  - subtyping
  - object composition
- Acording to an article, it was mentioned that implementation inheritance (extends relationship)
  is the one causing issues. Some of the issues caused
  - Coupling (Change a base class and there is a change that changes need to be replicated downwards)
- It is possible to use un-capitalized keys in the structs. However, uncapitalized kind of mean private in go land and capitalized
  means public and can be exported.
*/

package main

import (
	"fmt"
	"math"
)

type geometry interface {
	area() float64
	perim() float64
}

type rect struct {
	width, height float64
}

type circle struct {
	radius float64
}

func (r rect) area() float64 {
	return r.width * r.height
}

func (r rect) perim() float64 {
	return 2*r.width + 2*r.height
}

func (c circle) area() float64 {
	return math.Pi * c.radius * c.radius
}

func (c circle) perim() float64 {
	return 2 * math.Pi * c.radius
}

func measure(g geometry) {
	fmt.Println(g)
	fmt.Println(g.area())
	fmt.Println(g.perim())
}

func main() {
	r := rect{width: 3, height: 4}
	c := circle{radius: 5}

	measure(r)
	measure(c)
}