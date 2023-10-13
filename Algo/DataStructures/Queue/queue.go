package main

import "fmt"

type Queue struct {
	queue []int
}

func NewQueue() Queue {
	return Queue{queue: []int{}}
}

func (q *Queue) Enqueue(item int) {
	q.queue = append(q.queue, item)
}

func (q *Queue) Dequeue() (int, error) {
	if len(q.queue) == 0 {
		return 0, fmt.Errorf("no more items in the queue")
	}
	item := q.queue[0]
	q.queue = q.queue[1:]
	return item, nil
}

func main() {
	q := NewQueue()
	q.Enqueue(1)
	q.Enqueue(2)
	fmt.Println(q.Dequeue())
	fmt.Println(q.Dequeue())
}
