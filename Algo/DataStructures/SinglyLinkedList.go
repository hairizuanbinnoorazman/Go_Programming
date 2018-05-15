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
	Next  *Node
}

type List struct {
	Root *Node
}

func Init(val Node) *List {
	return &List{
		Root: &val,
	}
}

func (l *List) Append(val Node) *List {
	nodeWalk := l.Root
	for nodeWalk.Next != nil {
		nodeWalk = nodeWalk.Next
	}
	nodeWalk.Next = &val
	return l
}

func (l *List) Print() {
	nodeWalk := l.Root
	for nodeWalk.Next != nil {
		fmt.Println(nodeWalk.Value)
		nodeWalk = nodeWalk.Next
	}
	fmt.Println(nodeWalk.Value)
}

func (l *List) Len() int {
	nodeWalk := l.Root
	if nodeWalk == nil {
		return 0
	}
	count := 1
	for nodeWalk.Next != nil {
		count = count + 1
		nodeWalk = nodeWalk.Next
	}
	return count
}

func main() {
	aa := Init(Node{Value: "caca"})
	bb := Node{Value: "akcnjkanclk"}
	cc := Node{Value: "kmackmcklamkkmslc"}
	aa = aa.Append(bb)
	aa = aa.Append(cc)
	aa.Print()
	fmt.Println(aa.Len())

	zz := Init(Node{Value: "kanjkca"})
	zz.Print()
	fmt.Println(zz.Len())
}
