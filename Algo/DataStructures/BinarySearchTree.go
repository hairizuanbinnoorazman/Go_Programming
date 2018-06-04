/*
Binary search tree data structure

- Left node has data that is less that the node's key
- Right node has data that is more than the node's key
- All must be a single binary tree all the way down

This is a node/object based implementation
*/

package main

import "fmt"

type Node struct {
	data  int
	left  *Node
	right *Node
}

func (n *Node) Insert(data int) {
	if data >= n.data {
		if n.right == nil {
			n.right = &Node{data: data}
			return
		}
		tempNode := n.right
		tempNode.Insert(data)
	}

	if data <= n.data {
		if n.left == nil {
			n.left = &Node{data: data}
			return
		}
		tempNode := n.left
		tempNode.Insert(data)
	}
}

func (n *Node) Print() {
	fmt.Println(n.data)
	if n.left != nil {
		tempNode := n.left
		tempNode.Print()
	}
	if n.right != nil {
		tempNode := n.right
		tempNode.Print()
	}
}

type Tree struct {
	root *Node
}

func (t *Tree) Insert(data int) {
	if t.root == nil {
		t.root = &Node{data: data}
		return
	}
	lol := t.root
	lol.Insert(data)
}

func (t *Tree) Print() {
	if t.root == nil {
		fmt.Println("No elements inside tree")
	}
	t.root.Print()
}

func main() {
	hehe := Tree{}
	hehe.Insert(6)
	hehe.Insert(8)
	hehe.Insert(10)
	hehe.Insert(12)
	hehe.Insert(3)
	hehe.Print()
}
