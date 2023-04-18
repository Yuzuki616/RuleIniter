package main

import "testing"

func init() {
	LoadConfig("./config.json")
}

func TestParseRouteConf(t *testing.T) {
	m, _ := CheckMediaUnlock("global")
	ParseRouteConf(m)
}
