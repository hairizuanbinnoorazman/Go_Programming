package main

import "fmt"

type Node struct {
	Key   string
	Value string
}

type HashedMaps struct {
	items [100][]Node
}

func NewHashedMaps() HashedMaps {
	return HashedMaps{items: [100][]Node{}}
}

func (h HashedMaps) GetHash(key string) int {
	totalSum := 0
	for _, v := range key {
		totalSum = totalSum + int(v)
	}
	hashKey := totalSum % 100
	return hashKey
}

func (h *HashedMaps) Set(key, val string) {
	hashedKey := h.GetHash(key)
	if len(h.items[hashedKey]) == 0 {
		h.items[hashedKey] = []Node{Node{Key: key, Value: val}}
	}
	for _, i := range h.items[hashedKey] {
		if i.Key == key {
			return
		}
	}
	h.items[hashedKey] = append(h.items[hashedKey], Node{Key: key, Value: val})
}

func (h *HashedMaps) Get(key string) string {
	hashedKey := h.GetHash(key)
	if len(h.items[hashedKey]) == 0 {
		return ""
	}
	for _, item := range h.items[hashedKey] {
		if item.Key == key {
			return item.Value
		}
	}
	return ""
}

func main() {
	aa := NewHashedMaps()
	aa.Set("a", "sample value")
	aa.Set("ABB", "unknown sample value")

	fmt.Println(aa.Get("a"))
	fmt.Println(aa.Get("ABB"))
}
