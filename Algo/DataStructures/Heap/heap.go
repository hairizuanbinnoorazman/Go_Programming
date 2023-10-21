package main

import "fmt"

func main() {
	a := []int{1, 3, 5, 3, 6, 13, 10, 9, 8, 15, 17}
	fmt.Println(a)

	for i := (len(a) - 1) / 2; i >= 0; i-- {
		ArrHeapify(a, i)
	}
	fmt.Println(a)
}

// Slice implementation
// Root = 0
// Left handside = 2s + 1
// Right handside = 2s + 2
// Parent = (s-1)/2
func ArrHeapify(nums []int, node int) {
	lhsIdx := 2*node + 1
	rhsIdx := 2*node + 2
	largestIdx := node

	if lhsIdx < len(nums) {
		if nums[lhsIdx] > nums[largestIdx] {
			largestIdx = lhsIdx
		}
	}

	if rhsIdx < len(nums) {
		if nums[rhsIdx] > nums[largestIdx] {
			largestIdx = rhsIdx
		}
	}

	if largestIdx != node {
		tempVal := nums[node]
		nums[node] = nums[largestIdx]
		nums[largestIdx] = tempVal
		ArrHeapify(nums, largestIdx)
	}
}

// Node implementation
// Represented using nodes:
func nodeImplementation() {
	leftz := Node{value: 3}
	leftLeftz := Node{value: 4}
	rightz := Node{value: 2, right: &leftLeftz}
	topz := Node{value: 1, left: &leftz, right: &rightz}
	Printer(&topz)
	aa := Heapify(&topz)
	fmt.Println("after")
	Printer(aa)
	fmt.Println(aa.value)
}

type Node struct {
	value int
	left  *Node
	right *Node
}

func Heapify(n *Node) *Node {
	if n.left == nil && n.right == nil {
		return n
	}

	if n.left != nil {
		n.left = Heapify(n.left)
	}

	if n.right != nil {
		n.right = Heapify(n.right)
	}

	if n.left != nil {
		if n.value < n.left.value {
			tempLeft := n.left.left
			tempRight := n.left.right
			currentRight := n.right
			currentLeft := n.left
			currentLeft.right = currentRight
			currentLeft.left = n
			n.left = tempLeft
			n.right = tempRight
			n = currentLeft
		}

	}

	if n.right != nil {
		if n.value < n.right.value {
			tempLeft := n.right.left
			tempRight := n.right.right
			currentRight := n.right
			currentLeft := n.left
			currentRight.right = n
			currentRight.left = currentLeft
			n.left = tempLeft
			n.right = tempRight
			n = currentRight
		}

	}
	return n
}

func Printer(n *Node) {
	if n.left != nil {
		Printer(n.left)
	}
	fmt.Println(n.value)
	if n.right != nil {
		Printer(n.right)
	}
}
