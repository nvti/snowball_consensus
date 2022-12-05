package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"snowball/pkg/log"
	"snowball/pkg/utils"
	"strconv"
	"time"
)

type NewRequestHandler func([]byte) ([]byte, error)
type peerFoundHandler func(peerAddress string)
type Service struct {
	Config         ServerConfig
	peers          []string
	reqHandler     NewRequestHandler
	server         *http.Server
	updatePeerQuit chan struct{}
}

type ServerConfig struct {
	Name       string
	ProtocolID string
	Host       string
	Port       int
	Registry   string
}

func CreateService(config ServerConfig, reqHandler NewRequestHandler) (*Service, error) {
	if config.Host == "" {
		config.Host = "127.0.0.1"
	}

	if config.Port == 0 {
		port, err := utils.GetFreePort(config.Host)
		if err != nil {
			return nil, err
		}

		config.Port = port
	}

	service := &Service{
		Config:     config,
		peers:      []string{},
		reqHandler: reqHandler,
	}

	service.createHttpServer()

	return service, nil
}

// Start the service
func (s *Service) Start() error {
	go func() {
		s.server.ListenAndServe()
	}()

	// Register to registry
	err := s.callRegistry()
	if err != nil {
		s.server.Shutdown(context.TODO())
		return err
	}

	// Update list peers from registry every 500ms
	ticker := time.NewTicker(500 * time.Millisecond)
	s.updatePeerQuit = make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				peers := s.getListPeers()
				if peers != nil && len(peers) > 0 {
					s.peers = peers
				}
				break
			case <-s.updatePeerQuit:
				ticker.Stop()
				return
			}
		}
	}()

	return nil
}

func (s *Service) Stop() error {
	s.updatePeerQuit <- struct{}{}
	return s.server.Shutdown(context.TODO())
}

func (s *Service) newPeerHandler(peerAddress ...string) {
	log.Debug("Found peer:", peerAddress, ", connecting")
	s.peers = append(s.peers, peerAddress...)
}

func (s *Service) Peers() []string {
	return s.peers
}

func (s *Service) Send(peer string, data []byte) ([]byte, error) {
	resp, err := http.Post("http://"+peer+s.Config.ProtocolID, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New("status code " + resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (s *Service) callRegistry() error {
	data, err := json.Marshal(RegisterNodeReq{
		Port: s.Config.Port,
	})
	if err != nil {
		return err
	}

	resp, err := http.Post("http://"+s.Config.Registry, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Error(err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("fail register")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return err
	}

	respJson := &ListPeersResp{}
	err = json.Unmarshal(body, respJson)
	if err != nil {
		log.Error(err)
		return err
	}
	s.newPeerHandler(respJson.Peers...)

	return nil
}

func (s *Service) createHttpServer() {
	srv := &http.Server{Addr: s.Config.Host + ":" + strconv.Itoa(s.Config.Port)}
	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// New peer hook from registry
	http.HandleFunc("/peer", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		req := &NewNodeHook{}
		err = json.Unmarshal(body, req)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		s.newPeerHandler(req.Address)
		w.WriteHeader(http.StatusOK)
	})

	// Handle service
	http.HandleFunc(s.Config.ProtocolID, func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		respData, err := s.reqHandler(body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(respData)
	})

	s.server = srv
}

func (s *Service) getListPeers() []string {
	resp, err := http.Get("http://" + s.Config.Registry + "/peers")
	if err != nil {
		log.Error(err)
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return nil
	}
	respJson := &ListPeersResp{}
	err = json.Unmarshal(body, respJson)
	if err != nil {
		log.Error(err)
		return nil
	}

	return respJson.Peers
}
