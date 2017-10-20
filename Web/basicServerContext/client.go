package main

import (
	"fmt"
	"net/http"
	"log"
	"io/ioutil"
	"context"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	req = req.WithContext(ctx)

	// res, err := http.Get("http://localhost:8080")
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}

	value, _ := ioutil.ReadAll(res.Body)

	fmt.Println(string(value))
}

