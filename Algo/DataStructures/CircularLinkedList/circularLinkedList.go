package main

import "fmt"

type Node struct {
	Value    int
	Next     *Node
	Previous *Node
}

func Print(root *Node) {
	if root == nil {
		return
	}
	fmt.Println(root.Value)
	a := root.Next
	for a != root {
		fmt.Println(a.Value)
		a = a.Next
	}
}

func main() {
	aa := Node{Value: 1}
	bb := Node{Value: 2}
	cc := Node{Value: 3}

	aa.Next = &bb
	aa.Previous = &cc
	bb.Next = &cc
	bb.Previous = &aa
	cc.Next = &aa
	cc.Previous = &bb

	Print(&aa)
}
