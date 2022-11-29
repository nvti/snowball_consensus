package app

import (
	"encoding/json"
	"snowball/pkg/chain"
	"snowball/pkg/log"
	"snowball/pkg/p2p"
	"snowball/pkg/snowball"
)

const protocolID = "/snowball/1.0.0"
const serviceName = "snowball"

type DataReq struct {
	Index int
}

type DataResp struct {
	Data int
}

type Service struct {
	*chain.SnowballChain
	service *p2p.Service
}

func CreateService() (*Service, error) {
	s := &Service{}
	service, err := p2p.CreateService(p2p.ServerConfig{
		Name:       serviceName,
		ProtocolID: protocolID,
	}, s.handleRequest)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	s.SnowballChain = chain.NewConsensusChain(snowball.ConsensusConfig{
		K:       20,
		Alpha:   10,
		Beta:    10,
		MaxStep: 0,
	})

	s.service = service
	return s, nil
}

func (s *Service) handleRequest(reqData []byte) []byte {
	req := DataReq{}
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

	resp := DataResp{
		Data: block.Data,
	}
	respData, err := json.Marshal(resp)
	if err != nil {
		log.Error(err)
		return nil
	}

	return respData
}
