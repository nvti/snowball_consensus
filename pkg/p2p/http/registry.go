package http

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"snowball/pkg/log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Registry struct {
	Config          RegistryConfig
	peers           []string
	mu              sync.Mutex
	server          *http.Server
	healthCheckQuit chan struct{}
}

type RegistryConfig struct {
	Host string
	Port int
}

func CreateRegistry(config RegistryConfig) (*Registry, error) {
	r := &Registry{
		Config: config,
		peers:  []string{},
		mu:     sync.Mutex{},
	}

	r.createHttpServer()

	return r, nil
}

func (p *Registry) createHttpServer() {
	srv := &http.Server{Addr: p.Config.Host + ":" + strconv.Itoa(p.Config.Port)}

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// new peer coming
	http.HandleFunc("/", p.handleNewPeer)

	// Handler request get list current peers
	http.HandleFunc("/peers", p.handleGetListPeers)

	p.server = srv
}

// Start the service
func (p *Registry) Start() error {
	// Health check
	ticker := time.NewTicker(2 * time.Second)
	p.healthCheckQuit = make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				p.healthCheckPeers()
				break
			case <-p.healthCheckQuit:
				ticker.Stop()
				return
			}
		}
	}()

	return p.server.ListenAndServe()
}

func (p *Registry) Stop() error {
	p.healthCheckQuit <- struct{}{}
	return p.server.Shutdown(context.TODO())
}

// new peer coming
func (p *Registry) handleNewPeer(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	req := &RegisterNodeReq{}
	err = json.Unmarshal(body, req)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// get client ip
	remoteAddr := r.RemoteAddr
	s := strings.Split(remoteAddr, ":")
	if len(s) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	address := s[0] + ":" + strconv.Itoa(req.Port)
	log.Info("New peer coming: ", address)

	// Return current peers
	respData, err := json.Marshal(ListPeersResp{
		Peers: p.peers,
	})

	// Add new peer to peers
	p.mu.Lock()
	p.peers = append(p.peers, address)
	p.mu.Unlock()

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respData)
}

// Handler request get list current peers
func (p *Registry) handleGetListPeers(w http.ResponseWriter, r *http.Request) {
	respData, err := json.Marshal(ListPeersResp{
		Peers: p.peers,
	})

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respData)
}

// Todo: check node is alive
func (p *Registry) healthCheckPeers() {

}
