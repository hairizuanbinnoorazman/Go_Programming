package main

import (
	"fmt"
)

type Node struct {
	Value int
	Next *Node
}

var root = new(Node)

func addNode(t *Node, v int) int {
	if root == nil {
		t = &Node{v, nil}
		root = t
		return 0
	}

	if v == t.Value {
		fmt.Println("Node already exists:", v)
		return -1
	}

	if t.Next == nil {
		t.Next = &Node{v, nil}
		return -2
	}

	return addNode(t.Next, v)
}

func traverse(t *Node) {
	if t == nil {
		fmt.Println("Empty list!")
		return
	}

	for t != nil {
		fmt.Printf(" %d ", t.Value)
		t = t.Next
	}

	fmt.Println()
}

func main() {
	fmt.Println(root)
	addNode(root, 5)
	addNode(root, 3)
	addNode(root, 5)
	traverse(root)
}