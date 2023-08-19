package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/hairizuanbinnoorazman/basic-memory/store"
)

func main() {
	z := store.NewMemoryStore()
	hoho(z)
	time.Sleep(15 * time.Second)
	fmt.Println(len(z.View()))
	hoho(z)
	time.Sleep(15 * time.Second)
	fmt.Println(len(z.View()))
	hoho(z)
	time.Sleep(15 * time.Second)
	fmt.Println(len(z.View()))
}

func Adder(a store.Store, x int) {
	for i := 0; i < 100; i++ {
		a.Store(x)
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	}
}

func hoho(a store.Store) {
	for i := 0; i < 100; i++ {
		go Adder(a, i)
	}
}
