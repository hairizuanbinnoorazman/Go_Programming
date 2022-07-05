// This package is meant for building a self balancing BST (AVL)
package main

import "fmt"

type Node struct {
	Value int
	Left  *Node
	Right *Node
}

func InorderPrint(root *Node) {
	if root == nil {
		return
	}
	if root.Left != nil {
		InorderPrint(root.Left)
	}
	fmt.Println(root.Value)
	if root.Right != nil {
		InorderPrint(root.Right)
	}
}

func MaxDepth(root *Node) int {
	if root == nil {
		return 0
	}
	numL := MaxDepth(root.Left) + 1
	numR := MaxDepth(root.Right) + 1
	if numL >= numR {
		return numL
	}
	return numR
}

func PrintLevel(root *Node, currentLevel, level int) {
	if root == nil {
		return
	}
	if currentLevel == level {
		fmt.Println(root.Value)
	}
	PrintLevel(root.Left, currentLevel+1, level)
	PrintLevel(root.Right, currentLevel+1, level)
}

func Insert(root *Node, newNode *Node) *Node {
	if root == nil {
		return newNode
	}
	if newNode.Value <= root.Value {
		root.Left = Insert(root.Left, newNode)
	} else {
		root.Right = Insert(root.Right, newNode)
	}
	LH := MaxDepth(root.Left)
	RH := MaxDepth(root.Right)
	LHBalance := 0
	RHBalance := 0
	if root.Left != nil {
		LHBalance = MaxDepth(root.Left.Left) - MaxDepth(root.Left.Right)
	}
	if root.Right != nil {
		RHBalance = MaxDepth(root.Right.Left) - MaxDepth(root.Right.Right)
	}

	// Left hand side too heavy
	if (LH-RH) >= 2 && LHBalance >= 0 {
		newRoot := root.Left
		root.Left = newRoot.Right
		newRoot.Right = root
		return newRoot
	}

	// Right hand side too heavy
	if (LH-RH) <= -2 && RHBalance <= 0 {
		newRoot := root.Right
		root.Right = newRoot.Left
		newRoot.Left = root
		return newRoot
	}

	// Double rotation cases
	if (LH-RH) >= 2 && LHBalance < 0 {
		newRoot := root.Left.Right
		root.Left.Right = nil
		newRoot.Left = root.Left
		root.Left = newRoot.Right
		newRoot.Right = root
		return newRoot
	}

	// Double rotation cases
	if (LH-RH) <= -2 && RHBalance > 0 {
		newRoot := root.Right.Left
		root.Right.Left = nil
		newRoot.Right = root.Right
		root.Right = newRoot.Left
		newRoot.Left = root
		return newRoot
	}

	return root
}

func main() {
	aa := Node{Value: 30}
	bb := Node{Value: 20}
	cc := Node{Value: 10}
	dd := Node{Value: 15}
	ee := Node{Value: 17}
	ff := Node{Value: 18}

	zz := Insert(nil, &aa)
	zz = Insert(zz, &cc)
	zz = Insert(zz, &bb)
	zz = Insert(zz, &dd)

	for i := 1; i <= MaxDepth(zz); i++ {
		fmt.Printf("Print level %v\n", i)
		PrintLevel(zz, 1, i)
	}
	fmt.Println("Done")

	zz = Insert(zz, &ee)
	zz = Insert(zz, &ff)

	InorderPrint(zz)
	fmt.Println(MaxDepth(zz))

	for i := 1; i <= MaxDepth(zz); i++ {
		fmt.Printf("Print level %v\n", i)
		PrintLevel(zz, 1, i)
	}
}
