/*
Example of using interfaces to construct/compose an object

To run the following file, run the following command:
go run interfaces_2.go

Learnings:
- There is no "this" keyword. This is replace by p in `func (p *Person) lol() {...}`
- You are passing the values around by reference - for python is self, for plenty of other languages, its this
- By put the (p *Person) in front of the function, you are kind of saying that this function belongs to this "struct"/object
*/

package main

import (
	"fmt"
	"encoding/json"
)

type Person struct {
	Name string
	Address Address
}

type Address struct {
	Number string
	Street string
	City string
	State string
	Zip string
	Manga string
}

func (p *Person) Talk() {
	fmt.Println("Hi, my name is", p.Name)
}

func (p *Person) Location() {
	fmt.Println("I'm at", p.Address.Number, p.Address.Street,
		p.Address.City, p.Address.State, p.Address.Zip)
}

func (p *Person) FullDetails() {
	miao, eh := json.MarshalIndent(p, "", "  ")
	if eh != nil {
		panic("Again again!!")
	}
	s := string(miao)
	fmt.Println(s)
}

func main() {
	p  := Person{
		Name: "Steve",
		Address: Address{
			Number: "13",
			Street: "Main",
			City: "Manhattan",
			State: "NY",
			Zip: "01313",
		},
	}
	p.Talk()
	p.Location()
	p.FullDetails()
}