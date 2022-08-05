package search

import (
	"errors"
	"time"
)

type Result struct{ Title, URL string }

func Search(query string) ([]Result, error) {
	results := []Result{
		Web(query),
		Image(query),
		Video(query),
	}
	return results, nil
}

func SearchParallel(query string) ([]Result, error) {
	c := make(chan Result)
	go func() { c <- Web(query) }()
	go func() { c <- Image(query) }()
	go func() { c <- Video(query) }()

	return []Result{<-c, <-c, <-c}, nil
}

func SearchTimeout(query string, timeout time.Duration) ([]Result, error) {
	timer := time.After(timeout)
	c := make(chan Result, 3)
	go func() { c <- Web(query) }()
	go func() { c <- Image(query) }()
	go func() { c <- Video(query) }()

	var results []Result
	for i := 0; i < 3; i++ {
		select {
		case result := <-c:
			results = append(results, result)
		case <-timer:
			return results, errors.New("timed out")
		}
	}
	return results, nil
}
