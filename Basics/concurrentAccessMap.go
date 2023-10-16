package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {
	a := NewMemoryMapStore()

	for i := 0; i < 10000; i++ {
		val := strconv.Itoa(i)
		go a.Store(val, val)
	}

	time.Sleep(5 * time.Second)
	fmt.Println(len(a.items))

	for i := 0; i < 1000; i++ {
		val := strconv.Itoa(i)
		go a.Delete(val)
	}

	time.Sleep(5 * time.Second)
	fmt.Println(len(a.items))

}

type storeItem struct {
	Key   string
	Value string
}

type MemoryMapStore struct {
	items      map[string]string
	addChan    chan storeItem
	deleteChan chan string
}

func NewMemoryMapStore() *MemoryMapStore {
	initMap := map[string]string{}
	aChan := make(chan storeItem)
	dChan := make(chan string)
	m := MemoryMapStore{
		items:      initMap,
		addChan:    aChan,
		deleteChan: dChan,
	}
	go m.runner()
	return &m
}

func (m *MemoryMapStore) runner() {
	for {
		select {
		case x := <-m.addChan:
			m.items[x.Key] = x.Value
		case y := <-m.deleteChan:
			delete(m.items, y)
		}
	}
}
func (m *MemoryMapStore) Store(key, value string) {
	m.addChan <- storeItem{key, value}
}

func (m *MemoryMapStore) Get(key string) (value string) {
	return m.items[key]
}
func (m *MemoryMapStore) Delete(key string) {
	m.deleteChan <- key
}
