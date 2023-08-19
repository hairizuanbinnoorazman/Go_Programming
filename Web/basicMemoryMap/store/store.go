package store

import (
	"fmt"
)

type ZZ struct {
	ID  string
	Val int
}

type Store interface {
	Store(x ZZ)
	View() []ZZ
	Delete(id string)
}

type MemoryStore struct {
	items map[string]ZZ
	zz    chan ZZ
	del   chan string
}

func NewMemoryStore() *MemoryStore {
	m := MemoryStore{
		items: make(map[string]ZZ),
		zz:    make(chan ZZ),
		del:   make(chan string),
	}
	go m.start()
	return &m
}

func (m *MemoryStore) start() {
	for {
		select {
		case x := <-m.zz:
			m.items[x.ID] = x
		case y := <-m.del:
			fmt.Println("im here")
			delete(m.items, y)
			// default:
			// 	fmt.Println("Nothing available")
			// 	time.Sleep(1 * time.Second)
		}
	}
}

func (m *MemoryStore) Store(x ZZ) {
	m.zz <- x
}

func (m *MemoryStore) Printer() {
	fmt.Println(m.items)
}

func (m *MemoryStore) View() []ZZ {
	lister := []ZZ{}
	for _, v := range m.items {
		lister = append(lister, v)
	}
	return lister
}

func (m *MemoryStore) Delete(id string) {
	m.del <- id
}
