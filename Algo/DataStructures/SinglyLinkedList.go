/*
Singly linked list implementation in Golang

Some of the operations that needs to be done
1. Traversing
2. Searching
3. Insertion
4. Deletion
5. Sorting
6. Merging

NOTE: There is already a implementation of linked list in the container package from the standard library
 */

package main

import "fmt"

type Node struct {
	Value string
	Next *Node
}

// Append operation should only work at the end of the linked list
func (n *Node) append(node *Node) {
	for n.Next != nil {
		n = n.Next
	}
	n.Next = node
}

func (n *Node) print () {
	for n.Next != nil {
		fmt.Println(n.Value)
		n = n.Next
	}
	fmt.Println(n.Value)
}

func (n *Node) length() int {
	count :=  0
	for n.Next != nil {
		count = count + 1
		n = n.Next
	}
	return count + 1
}



func main() {
	first := Node{"a", nil}
	second := Node{"b", nil}
	third := Node{"c", nil}
	first.append(&second)
	first.append(&third)
	first.print()
	fmt.Println(first.length())
}


