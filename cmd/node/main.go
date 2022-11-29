package main

import (
	"flag"
	"math/rand"
	"snowball/app"
	"snowball/pkg/log"
	"snowball/pkg/p2p"
	"snowball/pkg/snowball"
	"time"
)

const protocolID = "/snowball/1.0.0"
const serviceName = "snowball"

var (
	name     string
	host     string
	port     int
	k        int
	alpha    int
	beta     int
	maxStep  int
	chainLen int
	nChoices int
)

func init() {
	rand.Seed(time.Now().UnixNano())

	flag.StringVar(&name, "name", "Client", "name of the node")
	flag.StringVar(&host, "host", "127.0.0.1", "Listen on interface")
	flag.IntVar(&port, "port", 0, "Listen on port, set to 0 to use random port")
	flag.IntVar(&k, "k", 3, "K parameter for snowball")
	flag.IntVar(&alpha, "alpha", 2, "Alpha parameter for snowball")
	flag.IntVar(&beta, "beta", 10, "Beta parameter for snowball")
	flag.IntVar(&maxStep, "maxStep", 0, "Max running step for snowball")
	flag.IntVar(&chainLen, "chainLen", 4, "Length of chain to sync")
	flag.IntVar(&nChoices, "nChoices", 2, "Number of possible choices")

	flag.Parse()
}

func main() {
	service, err := app.CreateService(app.ServiceConfig{
		ServerConfig: p2p.ServerConfig{
			Name:       serviceName,
			ProtocolID: protocolID,
			Host:       host,
			Port:       port,
		},
		ConsensusConfig: snowball.ConsensusConfig{
			K:       k,
			Alpha:   alpha,
			Beta:    beta,
			MaxStep: maxStep,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < chainLen; i++ {
		service.Add(rand.Intn(nChoices))
	}

	sleepTime := time.Duration(rand.Intn(1000))
	time.Sleep(sleepTime * time.Millisecond)

	service.Sync()

	log.Infof("%s: Sync finished=%v, data=", name, service.Finished)
	for _, block := range service.Blocks {
		log.Infof("%v", block.Data)
	}

	// sleep for make sure other nodes can receive the message
	time.Sleep(20 * time.Second)
	log.Info(name, ": Exit")
}
