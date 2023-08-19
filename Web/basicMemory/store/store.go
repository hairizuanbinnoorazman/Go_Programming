package store

import "fmt"

type Store interface {
	Store(x int)
	View() []int
}

type MemoryStore struct {
	items []int
	zz    chan int
}

func NewMemoryStore() *MemoryStore {
	m := MemoryStore{
		items: []int{},
		zz:    make(chan int),
	}
	go m.start()
	return &m
}

func (m *MemoryStore) start() {
	for {
		select {
		case x := <-m.zz:
			m.items = append(m.items, x)
		}
	}
}

func (m *MemoryStore) Store(x int) {
	m.zz <- x
}

func (m *MemoryStore) Printer() {
	fmt.Println(m.items)
}

func (m *MemoryStore) View() []int {
	dst := make([]int, len(m.items))
	copy(dst, m.items)
	return dst
}
