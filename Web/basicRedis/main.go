package main

import (
	"context"
	"fmt"
	"time"

	redis "github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		// Addr:     "localhost:9999",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	fmt.Println("Setting foo value")
	status := rdb.Set(context.TODO(), "foo", "zzz", 20*time.Second)
	if status.Err() != nil {
		panic(fmt.Sprintf("error observed: %v\n", status.Err()))
	}
	fmt.Printf("%+v\n", status)

	fmt.Println("Doing a ping command")
	rdb.Ping(context.TODO())

	fmt.Println("Doing a get command")
	val := rdb.Get(context.TODO(), "foo")
	fmt.Printf("Value of foo: %v\n", val.Val())
	fmt.Printf("Value of foo: %v\n", val.String())

	// zz := rdb.HSet(context.TODO(), "zzz", map[string]interface{}{"aa": "qcaca", "aqq": 12})
	// if zz.Err() != nil {
	// 	panic(fmt.Sprintf("zz error observed: %v\n", zz.Err()))
	// }
	// yy := rdb.HGet(context.TODO(), "zzz", "aqq")
	// fmt.Printf("Value of zzz-aqq: %v\n", yy.Val())
}
