package main

import (
	"encoding/json"
	"math/rand"
	"snowball/pkg/log"
	"snowball/pkg/p2p"
	"time"
)

const protocolID = "/snowball/1.0.0"
const serviceName = "snowball"

type DataReq struct {
	Index int
}

type DataResp struct {
	Data int
}

func handleStream(reqData []byte) []byte {
	req := DataReq{}
	err := json.Unmarshal(reqData, &req)
	if err != nil {
		log.Error(err)
		return nil
	}

	resp := DataResp{
		Data: 10,
	}
	respData, err := json.Marshal(resp)
	if err != nil {
		log.Error(err)
		return nil
	}

	return respData
}

func main() {
	var err error
	server, err := p2p.CreateService(p2p.ServerConfig{
		Name:       serviceName,
		ProtocolID: protocolID,
	}, handleStream)
	if err != nil {
		log.Panic(err)
	}
	defer server.Close()

	for {
		peers := server.Peers()
		if len(peers) != 0 {
			index := rand.Intn(len(peers))
			log.Info("Send to peer ", peers[index].ID)

			req := DataReq{
				Index: 1,
			}
			reqData, err := json.Marshal(req)
			if err != nil {
				log.Error(err)
				continue
			}

			resp, err := server.Send(peers[index], reqData)
			if err != nil {
				log.Error(err)
				continue
			}

			if resp != nil && len(resp) > 0 {
				log.Info("Resp: ", string(resp))
			}
		}

		time.Sleep(3 * time.Second)
	}

	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt)
	// <-c
}
