package main

import "fmt"

type Node struct {
	Value    int
	Next     *Node
	Previous *Node
}

func Print(root *Node) {
	n := root
	for n != nil {
		fmt.Println(n.Value)
		n = n.Next
	}
}

func Insert(root *Node, newNode *Node, loc int) *Node {
	counter := 0
	n := root
	var p *Node
	for n != nil {
		if counter == loc {
			newNode.Next = n
			newNode.Previous = p
			if n.Next != nil {
				n.Next.Previous = newNode
			}
			if p != nil {
				p.Next = newNode
				return root
			}
			return newNode
		}
		p = n
		n = n.Next
		counter = counter + 1
	}
	if counter == loc {
		newNode.Previous = p
		p.Next = newNode
		return root
	}
	return nil
}

func main() {
	aa := Node{Value: 1}
	bb := Node{Value: 2}
	cc := Node{Value: 3}

	aa.Next = &bb
	bb.Next = &cc
	bb.Previous = &aa
	cc.Previous = &bb

	Print(&aa)
	dd := Node{Value: 4}
	hh := Insert(&aa, &dd, 3)
	Print(hh)
}
