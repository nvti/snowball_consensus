package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"snowball/models"
	"snowball/pkg/log"
	"strconv"
	"strings"
	"sync"
	"time"
)

var peers = []string{}
var mu sync.Mutex
var host string
var port int

func init() {
	flag.StringVar(&host, "host", "127.0.0.1", "Listen on interface")
	flag.IntVar(&port, "port", 5001, "Listen on port")

	flag.Parse()
}

func main() {
	// new peer coming
	http.HandleFunc("/", handleNewPeer)

	// Handler request get list current peers
	http.HandleFunc("/peers", handleGetListPeers)

	// Health check
	ticker := time.NewTicker(2 * time.Second)
	healthCheckQuit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				healthCheckPeers()
				break
			case <-healthCheckQuit:
				ticker.Stop()
				return
			}
		}
	}()

	http.ListenAndServe(host+":"+strconv.Itoa(port), nil)
}

func handleNewPeer(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	req := &models.RegisterNodeReq{}
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
	respData, err := json.Marshal(models.ListPeersResp{
		Peers: peers,
	})

	// Add new peer to peers
	mu.Lock()
	peers = append(peers, address)
	mu.Unlock()

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respData)
}

func handleGetListPeers(w http.ResponseWriter, r *http.Request) {
	respData, err := json.Marshal(models.ListPeersResp{
		Peers: peers,
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
func healthCheckPeers() {

}
