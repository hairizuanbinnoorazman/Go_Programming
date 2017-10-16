package main

import (
	"testing"
)

func TestMessageParser(t *testing.T){
	data := messageParser(12, "cacc")
	if data !=  "12cacc"{
		t.Error("Return value from messageParser is wrong.")
	}
}