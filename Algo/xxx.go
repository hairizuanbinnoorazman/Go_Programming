package main

import (
	"container/heap"
	"fmt"
)

type hoho struct {
	Key int
	Val int
}

type hohoList struct {
	aa []hoho
}

func (h *hohoList) Len() int {
	return len(h.aa)
}

func (h *hohoList) Less(i, j int) bool {
	if h.aa[i].Key > h.aa[j].Key {
		return true
	}
	return false
}

func (h *hohoList) Swap(i, j int) {
	temp := h.aa[i]
	h.aa[i] = h.aa[j]
	h.aa[j] = temp
}

func (h *hohoList) Push(xx any) {
	h.aa = append(h.aa, xx.(hoho))
}

func (h *hohoList) Pop() any {
	if len(h.aa) > 0 {
		y := h.aa[0]
		h.aa = h.aa[1:]
		return y
	}
	return hoho{}
}

func main() {
	aa := hoho{Key: 1, Val: 1}
	bb := hoho{Key: 2, Val: 2}
	cc := hoho{Key: 3, Val: 3}
	yy := hohoList{aa: []hoho{cc, bb, aa}}
	fmt.Println(yy)
	heap.Init(&yy)
	fmt.Println(yy)
	for i := 0; i < 3; i++ {
		fmt.Println(yy.Pop())
		fmt.Println(yy)
	}
}
