package chain

import (
	"snowball/pkg/log"
	"snowball/pkg/snowball"
)

type wrapperClient struct {
	client    Client[int]
	SyncIndex int
}

func (c *wrapperClient) Preference() (int, error) {
	return c.client.Preference(c.SyncIndex)
}

// Chain a Linear chain implementation
type ConsensusChain struct {
	SimpleLinearChain[int]
	Consensus *snowball.Consensus[int]
	clients   []*wrapperClient
	syncIndex int
	Finished  chan bool
}

func NewConsensusChain(config snowball.ConsensusConfig) *ConsensusChain {
	return &ConsensusChain{
		Consensus: snowball.NewConsensus[int](config),
		Finished:  make(chan bool),
	}
}

func (c *ConsensusChain) Preference(index int) (int, error) {
	block, err := c.Get(index)
	if err != nil {
		return 0, err
	}

	return block.Data, nil
}

func (c *ConsensusChain) SetClients(clients []Client[int]) {
	consensusClients := []snowball.Client[int]{}
	c.clients = []*wrapperClient{}
	for _, client := range clients {
		cl := &wrapperClient{
			client:    client,
			SyncIndex: c.syncIndex,
		}
		consensusClients = append(consensusClients, cl)
		c.clients = append(c.clients, cl)
	}
	c.Consensus.SetClients(consensusClients)
}

func (c *ConsensusChain) Sync() {
	go func() {
		for i := 0; i < c.Length(); i++ {
			c.syncIndex = i
			for _, client := range c.clients {
				client.SyncIndex = i
			}
			block := c.Blocks[i]
			c.Consensus.SetPreference(block.Data).SetUpdateHandler(func(preference int) {
				block.Data = preference
			}).Sync()
			finished := <-c.Consensus.Finished
			log.Info("Block ", i, ": finished=", finished, " preference=", block.Data)
		}
		c.Finished <- true
	}()
}
