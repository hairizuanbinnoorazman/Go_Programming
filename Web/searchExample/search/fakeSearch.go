package search

import (
	"fmt"
	"math/rand"
	"time"
)

type SearchFunc func(query string) Result

func FakeSearch(kind, title, url string) SearchFunc {
	return func(query string) Result {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return Result{
			Title: fmt.Sprintf("%s(%q): %s", kind, query, title),
			URL:   url,
		}
	}
}

var (
	Web   = FakeSearch("web", "The Go Programming Language", "http://golang.org")
	Image = FakeSearch("image", "The Go gopher", "https://blog.golang.org/gopher/gopher.png")
	Video = FakeSearch("video", "Concurrentcy is not Parellelism", "https://www.youtube.com/watch?v=cN_DpYBzKso")
)
