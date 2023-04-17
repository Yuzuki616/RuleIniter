package main

import (
	"log"
	"testing"
)

func init() {
	LoadConfig("./config.json")
}

func TestCheckMediaUnlock(t *testing.T) {
	log.Print(CheckMediaUnlock("global"))
}
