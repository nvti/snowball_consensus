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
)

var peers = []string{}
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
		ip := s[0]

		peers = append(peers, ip+":"+strconv.Itoa(req.Port))
		w.WriteHeader(http.StatusOK)

	})
}
