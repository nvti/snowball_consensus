package http

import (
	"errors"
	"snowball/pkg/log"
	"strconv"

	"github.com/go-resty/resty/v2"
)

type NewRequestHandler func([]byte) []byte

type Service struct {
	Config     ServerConfig
	peers      []string
	reqHandler NewRequestHandler
}

type ServerConfig struct {
	Name       string
	ProtocolID string
	Host       string
	Port       int
}

func CreateService(config ServerConfig, reqHandler NewRequestHandler) (*Service, error) {
	service := &Service{
		Config:     config,
		peers:      []string{},
		reqHandler: reqHandler,
	}

	app := CreateHttpServer(config.Host, config.Port, service.newPeerHandler)

	go func() {
		app.Listen(config.Name + ":" + strconv.Itoa(config.Port))
	}()

	return service, nil
}

func (s *Service) newPeerHandler(peerAddress string) {
	log.Debug("Found peer:", peerAddress, ", connecting")
	s.peers = append(s.peers, peerAddress)
}

func (s *Service) Peers() []string {
	return s.peers
}

func (s *Service) Send(peer string, data interface{}) ([]byte, error) {
	req := resty.New().SetCloseConnection(true).R().SetBody(data)
	resp, err := req.Post(peer + "/" + s.Config.ProtocolID)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return nil, errors.New("status code " + resp.Status())
	}
	return resp.Body(), nil
}
