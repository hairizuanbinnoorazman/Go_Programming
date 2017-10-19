/*
	Some learnings:

	To write a infinite loop, use a for {}.
	There seems to be no more while loops?

	In order to make a random sleeper -> Convert to type accordingly.
	Convert int to Duration and then multiply it to the time.Second to properly get the time value
*/

package main

import (
	"log"
	"fmt"
	"time"
	"math/rand"
)

func randomFire(c chan int) {
	log.Println("Random Fire Function Started")
	time.Sleep(time.Duration(rand.Intn(5) + 1) * time.Second)
	randomNumber := rand.Intn(1000)
	fmt.Println("Random Fire Function Generated the following number:", randomNumber)
	c <- randomNumber
}

func randomFireLoop(c chan int) {
	for {
		time.Sleep(2 * time.Second)
		go randomFire(c)
	}
}

func main() {
	a := make(chan int)
	go randomFireLoop(a)
	for {
		fmt.Println(<- a)
	}
}