package app

import (
	"encoding/json"
	"errors"
	"math/rand"
	"snowball/pkg/chain"
	"snowball/pkg/log"
	"snowball/pkg/p2p/libp2p"
	"snowball/pkg/snowball"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
)

type dataReq struct {
	Index int
}

type dataResp struct {
	Data int
}

type Service struct {
	*chain.SnowballChain
	service *libp2p.Service
}

type ServiceConfig struct {
	libp2p.ServerConfig
	snowball.ConsensusConfig
}

func CreateService(config ServiceConfig) (*Service, error) {
	s := &Service{}
	service, err := libp2p.CreateService(config.ServerConfig, s.handleRequest)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	s.SnowballChain = chain.NewConsensusChain(config.ConsensusConfig).
		SetRequestAnswerHandler(s.onRequestAnswerHandler)

	s.service = service
	return s, nil
}

func (s *Service) onRequestAnswerHandler(index int, k int) []int {
	peers := s.service.Peers()
	indexArr := rand.Perm(len(peers))
	answers := []int{}
	for i, j := 0, 0; i < k && j < len(peers); j++ {
		p := peers[indexArr[j]]
		if p != nil {
			resp, err := s.SendRequest(p, index)
			if err != nil || resp == nil {
				continue
			}
			answers = append(answers, resp.Data)
			i++
		}
	}

	if len(answers) < k {
		// Wait for other node connect
		sleepTime := time.Duration(rand.Intn(1000))
		time.Sleep(sleepTime * time.Millisecond)
	}
	return answers
}
func (s *Service) SendRequest(peer *peer.AddrInfo, index int) (*dataResp, error) {
	req := dataReq{
		Index: index,
	}
	reqData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	respData, err := s.service.Send(peer, reqData)
	if err != nil {
		return nil, err
	}
	if respData == nil {
		return nil, errors.New("no response")
	}

	resp := &dataResp{}
	err = json.Unmarshal(respData, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Service) handleRequest(reqData []byte) []byte {
	req := dataReq{}
	err := json.Unmarshal(reqData, &req)
	if err != nil {
		log.Error(err)
		return nil
	}

	block, err := s.Get(req.Index)
	if err != nil {
		log.Error(err)
		return nil
	}

	resp := dataResp{
		Data: block.Data,
	}
	respData, err := json.Marshal(resp)
	if err != nil {
		log.Error(err)
		return nil
	}

	return respData
}
