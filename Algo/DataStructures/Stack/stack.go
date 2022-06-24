package main

import "fmt"

type Stack struct {
	stack []int
}

func NewStack() Stack {
	return Stack{stack: []int{}}
}

func (s *Stack) AddToStack(item int) {
	s.stack = append(s.stack, item)
}

func (s *Stack) RemoveItemFromStack() (int, error) {
	if len(s.stack) == 0 {
		return 0, fmt.Errorf("stack is empty")
	}
	temp := s.stack[len(s.stack)-1]
	s.stack = s.stack[0 : len(s.stack)-1]
	return temp, nil
}

func main() {
	aa := NewStack()
	aa.AddToStack(1)
	aa.AddToStack(2)
	fmt.Println(aa.stack)
	fmt.Println(aa.RemoveItemFromStack())
	fmt.Println(aa.RemoveItemFromStack())
}
