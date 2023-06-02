package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	f, _ := os.OpenFile("test.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	bw := bufio.NewWriter(f)
	iter := 0
	expectedCounter := 0

	for {
		time.Sleep(1 * time.Second)
		select {
		case a := <-sigs:
			fmt.Printf("Received signal to end application :: %v\n", a)
			fmt.Printf("%v items left in buffer. Flushing it\n", iter)
			bw.Flush()
			iter = 0
			fmt.Printf("expected number of lines: %v\n", expectedCounter)
			os.Exit(0)
		default:
			iter = iter + 1
			fmt.Println("generated new item")
			expectedCounter = expectedCounter + 1
			dataVal := fmt.Sprintf("generated val: %v\n", rand.Int())
			bw.WriteString(dataVal)
			if iter >= 10 {
				bw.Flush()
				iter = 0
			}
		}
	}
}
