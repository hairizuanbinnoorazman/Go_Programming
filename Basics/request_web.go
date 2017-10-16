package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
)

func main() {
	res, err := http.Get("http://www.jpf.go.jp/")
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("What the heck")
	}

	fmt.Println(res.StatusCode)
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Failed at reading response")
	}
	miao := string(data)
	fmt.Println(miao)

	defer res.Body.Close()
}