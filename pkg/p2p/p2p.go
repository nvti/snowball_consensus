package p2p

import (
	"crypto/rand"
	"fmt"
	"snowball/pkg/log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
)

func InitP2P(protocolID string, host string, port int, handler network.StreamHandler) (server host.Host, err error) {
	// Setting the TCP port as 0 makes libp2p choose an available port for us.
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
	server, err = libp2p.New(libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/%s/tcp/%d", host, port)), libp2p.Identity(priv))
	if err != nil {
		return
	}

	log.Debug("Addresses:", server.Addrs())
	log.Debug("ID:", server.ID())

	// This gets called every time a peer connects and opens a stream to this node.
	server.SetStreamHandler(protocol.ID(protocolID), handler)

	return
}
