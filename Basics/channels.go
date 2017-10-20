/*
	Some learnings:

	To write a infinite loop, use a for {}.
	There seems to be no more while loops?

	In order to make a random sleeper -> Convert to type accordingly.
	Convert int to Duration and then multiply it to the time.Second to properly get the time value

	There is no length property for unbuffered channels. Only if you add buffer to it will you be able to get the length out.
*/

package main

import (
	"log"
	"fmt"
	"time"
	"math/rand"
)


func randomFireConsumer(c chan int, worker int) {
	for {
		log.Println("Random Fire Function Started for worker", worker)
		time.Sleep(time.Duration(rand.Intn(5) + 1) * time.Second)
		fmt.Println("Random Fire Function Generated the following number:", <- c)
	}
}

func randomFireGenerator(c chan int) {
	for {
		time.Sleep(100 * time.Millisecond)
		c <- rand.Intn(1000)
	}
}

func main() {
	a := make(chan int, 200)
	go randomFireGenerator(a)

	for i := 0; i < 6; i++ {
		go randomFireConsumer(a, i)
	}

	for {
		time.Sleep(5 * time.Second)
		fmt.Println("Number of i terms in the channel currently:", len(a))
	}
}