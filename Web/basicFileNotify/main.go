package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
)

type config struct {
	Name string `json:"name"`
}

func readConfig() config {
	raw, _ := ioutil.ReadFile("config.json")
	var c config
	json.Unmarshal(raw, &c)
	return c
}

func main() {
	currentConfig := readConfig()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	watcher.Add("config.json")

	for {
		time.Sleep(1 * time.Second)
		select {
		case event := <-watcher.Events:
			fmt.Printf("Event happened: %v\n", event)
			if event.Has(fsnotify.Write) {
				currentConfig = readConfig()
				fmt.Println("configuration is updated")
			}
		default:
			fmt.Printf("%+v\n", currentConfig)
		}
	}
}
