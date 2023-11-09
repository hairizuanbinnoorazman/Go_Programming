package main

import (
	"fmt"
	"sort"
)

type ppp struct {
	Key int
	Val int
}

type pppList struct {
	ll []ppp
}

func (p *pppList) Len() int {
	return len(p.ll)
}

func (p *pppList) Less(i, j int) bool {
	if p.ll[i].Key < p.ll[j].Key {
		return true
	}
	return false
}

func (p *pppList) Swap(i, j int) {
	tempVal := p.ll[i]
	p.ll[i] = p.ll[j]
	p.ll[j] = tempVal
}

func main() {
	aa := ppp{Key: 9, Val: 89}
	bb := ppp{Key: 7, Val: 11}
	cc := ppp{Key: 4, Val: 12}
	zz := pppList{ll: []ppp{aa, bb, cc}}
	fmt.Println(zz)
	sort.Sort(&zz)
	fmt.Println(zz)
}
