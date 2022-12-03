package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"snowball/models"
	"snowball/pkg/log"
	"strconv"
	"strings"
	"sync"
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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		req := &models.RegisterNodeReq{}
		err = json.Unmarshal(body, &req)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		remoteAddr := r.RemoteAddr
		s := strings.Split(remoteAddr, ":")
		if len(s) == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		address := s[0] + ":" + strconv.Itoa(req.Port)
		log.Info("New peer coming: ", address)

		// Notify other peers
		newPeerReq, err := json.Marshal(&models.NewNodeHook{
			Address: address,
		})
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for _, peer := range peers {
			if peer == address {
				continue
			}
			go func(peer string) {
				_, err := http.Post("http://"+peer+"/peer", "application/json", bytes.NewBuffer(newPeerReq))
				if err != nil {
					log.Error("Error when send hook to ", peer, ", error=", err)
				}
			}(peer)
		}
		mu.Lock()
		peers = append(peers, address)

		respData, err := json.Marshal(models.ListPeersResp{
			Peers: peers,
		})
		mu.Unlock()

		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(respData)
	})

	http.ListenAndServe(host+":"+strconv.Itoa(port), nil)
}
