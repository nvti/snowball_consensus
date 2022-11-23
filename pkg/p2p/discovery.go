package p2p

import (
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

// PeerFoundHandler handler when a new peer found
type PeerFoundHandler func(host host.Host, peerInfo peer.AddrInfo)

type discoveryNotifee struct {
	host    host.Host
	handler PeerFoundHandler
}

func (n *discoveryNotifee) HandlePeerFound(peerInfo peer.AddrInfo) {
	n.handler(n.host, peerInfo)
}

// InitDiscovery Initialize Discovery service with mDNS
func InitDiscovery(host host.Host, serviceName string, handler PeerFoundHandler) error {
	notifee := &discoveryNotifee{}
	notifee.handler = handler
	notifee.host = host

	ser := mdns.NewMdnsService(host, serviceName, notifee)
	if err := ser.Start(); err != nil {
		return err
	}

	return nil
}
