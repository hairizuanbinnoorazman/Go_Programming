/*

*/

package main

import (
	"os"
	"bufio"
	"fmt"
	"context"
	"time"
	"log"
)


func sleepAndTalk(ctx context.Context, t time.Duration, msg string) {
	log.Println("sleepAndTalk_First starts")
	defer log.Println("sleepAndTalk_First end")

	select {
	case <- time.After(t):
		fmt.Println(msg)
		// sleepAndTalkSecond(ctx, t)
	case <- ctx.Done():
		log.Println("Error on the first sleep and talk", ctx.Err())
	}
}


func sleepAndTalkSecond(ctx context.Context, t time.Duration) {
	log.Println("sleepAndTalk_Second starts")
	defer log.Println("sleepAndTalk_Second end")

	select {
	case <- time.After(t):
		fmt.Println("the second sleep and talk")
	case <- ctx.Done():
		log.Println("Error on the sleep and talk the second", ctx.Err())
	}
}



func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		s := bufio.NewScanner(os.Stdin)
		s.Scan()
		cancel()
	}()

	sleepAndTalk(ctx, 5*time.Second, "hello")
	sleepAndTalkSecond(ctx, 5*time.Second)
}