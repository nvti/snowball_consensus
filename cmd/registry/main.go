package main

import (
	"flag"
	"log"
	"snowball/app"
)

var peers = []string{}
var host string
var port int

func init() {
	flag.StringVar(&host, "host", "127.0.0.1", "Listen on interface")
	flag.IntVar(&port, "host", 5001, "Listen on port")

	flag.Parse()
}

func main() {
	_, err := app.CreateService(app.ServiceConfig{})
	if err != nil {
		log.Fatal(err)
	}
}
