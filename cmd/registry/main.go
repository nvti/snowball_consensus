package main

import (
	"flag"
	"snowball/pkg/log"
	"snowball/pkg/p2p/http"
)

var host string
var port int

func init() {
	flag.StringVar(&host, "host", "127.0.0.1", "Listen on interface")
	flag.IntVar(&port, "port", 5001, "Listen on port")

	flag.Parse()
}

func main() {
	service, err := http.CreateRegistry(http.RegistryConfig{
		Host: host,
		Port: port,
	})

	if err != nil {
		log.Fatal(err)
	}

	if err := service.Start(); err != nil {
		log.Fatal(err)
	}
}
