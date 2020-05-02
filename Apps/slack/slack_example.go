// Example of using Go to ping Slack
// This would ping a message by the text message passed via postMessage function
//
// In order to utilize this file, use the command: go run slack_example.go
// Else, generate a binary file by running the command: go build slack_example.go

package main

import (
	"fmt"
	
	"encoding/json"
	"bytes"
	"io/ioutil"

	"net/http"
)

type message struct {
	Text string    `json:"text"`
}

func postMessage(msg string) {
	slackUrl := "https://hooks.slack.com/services/<<KEYS>>"

	// Create a reader to be used by http.Post
	response := message{Text: msg}
	body, _ := json.Marshal(response)
	byteBody := bytes.NewReader(body)

	res, err := http.Post(slackUrl, "application/json", byteBody)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Try again later.")
	}

	fmt.Println("Status of response:", res.Status)
	fmt.Println("Status code of response:", res.StatusCode)
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(content))
}

func main() {
	fmt.Println("A test application to fire a message into Slack")
	postMessage("init")
}