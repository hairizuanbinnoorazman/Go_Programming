package main

import (
	"fmt"
	"time"

	"github.com/hairizuanbinnoorazman/basic-memory-map/store"
)

func main() {
	z := store.NewMemoryStore()
	z.Store(store.ZZ{ID: "aa", Val: 12})
	fmt.Println(z.View())
	z.Store(store.ZZ{ID: "bb", Val: 13})
	z.Store(store.ZZ{ID: "cc", Val: 14})
	fmt.Println(z.View())
	z.Delete("bb")
	time.Sleep(5 * time.Second)
	fmt.Println(z.View())
}
