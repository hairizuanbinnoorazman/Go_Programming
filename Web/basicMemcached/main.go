package main

import (
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

func main() {
	mc := memcache.New("localhost:11211")
	mc.Set(&memcache.Item{Key: "foo", Value: []byte("my value")})
	zz, err := mc.Get("foo")
	if err != nil {
		panic(fmt.Sprintf("didnt expect error from gettting values from memcached %v", err))
	}
	fmt.Printf("Value of foo: %v\n", string(zz.Value))

	addErr := mc.Add(&memcache.Item{Key: "foo", Value: []byte("new value")})
	if addErr != nil {
		fmt.Printf("Add error: %v\n", addErr)
	}

	appendErr := mc.Append(&memcache.Item{Key: "foo", Value: []byte("new value")})
	if addErr != nil {
		fmt.Printf("Add error: %v\n", appendErr)
	}

	pp, _ := mc.Get("foo")
	fmt.Printf("Value of foo: %v\n", string(pp.Value))

	mc.Set(&memcache.Item{Key: "yar", Value: []byte("yar"), Expiration: 10})
	time.Sleep(5 * time.Second)
	yy, err := mc.Get("yar")
	if err != nil {
		panic(fmt.Sprintf("didnt expect error from gettting values from memcached %v\n", err))
	}
	fmt.Printf("Value of yar: %v\n", string(yy.Value))

	time.Sleep(6 * time.Second)
	_, err = mc.Get("yar")
	if err != nil {
		fmt.Printf("Expeccted error: %v\n", err)
	}

}
