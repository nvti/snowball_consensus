package chain

import (
	"snowball/pkg/snowball"
)

// Chain a Linear chain implementation
type ConsensusChain struct {
	SimpleLinearChain[int]
	Consensus *snowball.Consensus[int]
}

func NewConsensusChain(config snowball.ConsensusConfig) *ConsensusChain {
	return &ConsensusChain{
		Consensus: snowball.NewConsensus[int](config),
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
	c.Consensus.SetClients(clients)
}

func (c *ConsensusChain) Sync() {
	for i := 0; i < c.Length(); i++ {
		block := c.Blocks[i]
		c.Consensus.SetPreference(block.Data).SetUpdateHandler(func(preference int) {
			block.Data = preference
		}).Sync()
	}
}
