package p2p

import (
	"context"
	"io"
	"snowball/pkg/log"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/protocol"
)

type Service struct {
	Config     ServerConfig
	server     host.Host
	peers      []*peer.AddrInfo
	reqHandler NewRequestHandler
}

type ServerConfig struct {
	Name       string
	ProtocolID string
	Host       string
	Port       int
}

type NewRequestHandler func([]byte) []byte

func CreateService(config ServerConfig, handler NewRequestHandler) (*Service, error) {
	var host string
	if config.Host == "" {
		host = "0.0.0.0"
	} else {
		host = config.Host
	}

	service := &Service{
		Config: ServerConfig{
			Name:       config.Name,
			ProtocolID: config.ProtocolID,
			Port:       config.Port,
			Host:       host,
		},
		reqHandler: handler,
		peers:      []*peer.AddrInfo{},
	}

	server, err := InitP2P(service.Config.ProtocolID, service.Config.Host, service.Config.Port, service.newStreamHandler)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// Discovery other node
	err = InitDiscovery(server, service.Config.Name, service.newPeerHandler)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	service.server = server

	return service, nil
}

func (s *Service) Send(peer *peer.AddrInfo, data []byte) ([]byte, error) {
	ctx := context.Background()

	// open a stream, this stream will be handled by handleStream other end
	stream, err := s.server.NewStream(ctx, peer.ID, protocol.ID(s.Config.ProtocolID))
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer stream.Close()

	_, err = stream.Write(data)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	stream.CloseWrite()

	out, err := io.ReadAll(stream)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return out, nil
}

func (s *Service) Peers() []*peer.AddrInfo {
	return s.peers
}

func (s *Service) Close() error {
	return s.server.Close()
}

func (s *Service) newPeerHandler(peer peer.AddrInfo) {
	log.Info("Found peer:", peer, ", connecting")
	s.server.Peerstore().AddAddrs(peer.ID, peer.Addrs, peerstore.PermanentAddrTTL)
	s.peers = append(s.peers, &peer)
}

func (s *Service) newStreamHandler(stream network.Stream) {
	defer stream.Close()
	reqData, err := io.ReadAll(stream)
	if err != nil {
		log.Error(err)
		return
	}

	respData := s.reqHandler(reqData)

	if len(respData) > 0 {
		_, err = stream.Write(respData)
		if err != nil {
			log.Error(err)
			return
		}
	}
}
