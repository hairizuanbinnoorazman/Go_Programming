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

func InorderPrint(root *Node) {
	if root == nil {
		return
	}
	if root.left != nil {
		InorderPrint(root.left)
	}
	fmt.Println(root.data)
	if root.right != nil {
		InorderPrint(root.right)
	}
}

func InsertNode(root *Node, item *Node) *Node {
	if root == nil {
		return item
	}
	if item.data < root.data {
		if root.left != nil {
			InsertNode(root.left, item)
		} else {
			root.left = item
		}
	}
	if item.data > root.data {
		if root.right != nil {
			InsertNode(root.right, item)
		} else {
			root.right = item
		}
	}
	return root
}

func ReverseBinaryTree(root *Node) {
	if root == nil {
		return
	}
	temp := root.left
	temp2 := root.right
	root.right = temp
	root.left = temp2
	if root.left != nil {
		ReverseBinaryTree(root.left)
	}
	if root.right != nil {
		ReverseBinaryTree(root.right)
	}
	return
}

func main() {
	// hehe := Tree{}
	// hehe.Insert("A")
	// hehe.Insert("B")
	// hehe.Insert("C")
	// hehe.Insert("D")
	// hehe.Insert("E")
	// hehe.Print()

	A := Node{data: 9}
	B := Node{data: 2}
	C := Node{data: 12}
	D := Node{data: 4}
	E := Node{data: 5}
	F := Node{data: -2}
	G := Node{data: 7}

	// A.left = &B
	// A.right = &C
	// B.left = &D
	// B.right = &E
	// C.left = &F
	// C.right = &G

	aa := InsertNode(nil, &D)
	aa = InsertNode(aa, &E)
	aa = InsertNode(aa, &C)
	aa = InsertNode(aa, &A)
	aa = InsertNode(aa, &B)
	aa = InsertNode(aa, &G)
	aa = InsertNode(aa, &F)

	InorderPrint(aa)
	fmt.Println("reversing")
	ReverseBinaryTree(aa)
	fmt.Println("2nd print")
	InorderPrint(aa)
}
