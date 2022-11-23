package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"snowball/pkg/log"
	"snowball/pkg/p2p"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

const protocolID = "/snowball/1.0.0"
const serviceName = "snowball"

func handleStream(stream network.Stream) {
	log.Info("Got a new stream!")

	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	go readData(rw)
	go writeData(rw, "handleStream")

	// 'stream' will stay open until you close it (or the other side closes it).
}

func readData(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			log.Error("Error reading from buffer: ", err)
			break
		}

		if str == "" {
			return
		}
		if str != "\n" {
			fmt.Printf("%s\n> ", str)
		}

	}
}

func writeData(rw *bufio.ReadWriter, name string) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(name, " > ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			log.Error("Error reading from stdin: ", err)
			break
		}

		_, err = rw.WriteString(fmt.Sprintf("%s\n", sendData))
		if err != nil {
			log.Error("Error writing to buffer: ", err)
			break
		}
		err = rw.Flush()
		if err != nil {
			log.Error("Error flushing buffer: ", err)
			break
		}
	}
}

func newPeerHandler(host host.Host, peer peer.AddrInfo) {
	log.Info("Found peer:", peer, ", connecting")
	ctx := context.Background()

	if err := host.Connect(ctx, peer); err != nil {
		log.Error("Connection failed:", err)
	}

	// open a stream, this stream will be handled by handleStream other end
	stream, err := host.NewStream(ctx, peer.ID, protocol.ID(protocolID))

	if err != nil {
		log.Error("Stream open failed", err)
	} else {
		rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

		go writeData(rw, "NewStream")
		go readData(rw)
		log.Info("Connected to:", peer)
	}
}

func main() {
	host, err := p2p.Init(protocolID, handleStream)
	if err != nil {
		log.Panic(err)
	}
	defer host.Close()

	// Discovery other node
	err = p2p.InitDiscovery(host, serviceName, newPeerHandler)
	if err != nil {
		log.Panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
