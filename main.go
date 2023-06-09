package main

import (
	"flag"
	"log"
)

var path = flag.String("conf", "./config.json", "config file")
var region = flag.String("region", "all", "region")
var version = "default"

func main() {
	flag.Parse()
	log.Println("RuleIniter Ver", version)
	err := LoadConfig(*path)
	if err != nil {
		log.Printf("load config error: %s", err)
		return
	}
	m, err := CheckMediaUnlock(*region)
	if err != nil {
		log.Printf("some check error: %s", err)
	}
	log.Println("Can not unlock list:")
	for i := range m {
		log.Println(m[i])
	}
	err = ParseRouteConf(m)
	if err != nil {
		log.Printf("parse error: %s", err)
	}
}
