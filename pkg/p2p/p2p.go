package p2p

import (
	"snowball/pkg/log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
)

func Init(protocolID string, handler network.StreamHandler) (host host.Host, err error) {
	// Setting the TCP port as 0 makes libp2p choose an available port for us.
	host, err = libp2p.New(libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/0"))
	if err != nil {
		return
	}

	log.Info("Addresses:", host.Addrs())
	log.Info("ID:", host.ID())

	// This gets called every time a peer connects and opens a stream to this node.
	host.SetStreamHandler(protocol.ID(protocolID), handler)

	return
}
