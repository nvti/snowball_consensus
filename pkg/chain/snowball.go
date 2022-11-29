package chain

import (
	"snowball/pkg/log"
	"snowball/pkg/snowball"
)

type OnRequestAnswerHandler func(index int, k int) []int

// Chain a Linear chain implementation
type SnowballChain struct {
	SimpleLinearChain[int]
	Consensus   *snowball.Consensus[int]
	onReqAnswer OnRequestAnswerHandler
	syncIndex   int
	Running     bool
	Finished    bool
}

func NewConsensusChain(config snowball.ConsensusConfig) *SnowballChain {
	return &SnowballChain{
		Consensus: snowball.NewConsensus[int](config),
		Finished:  false,
	}
}

func (c *SnowballChain) Preference(index int) (int, error) {
	block, err := c.Get(index)
	if err != nil {
		return 0, err
	}

	return block.Data, nil
}

func (c *SnowballChain) SetRequestAnswerHandler(handler OnRequestAnswerHandler) *SnowballChain {
	c.onReqAnswer = handler
	return c
}

func (c *SnowballChain) Sync() {
	if c.Running {
		return
	}

	c.Running = true
	finished := true
	for i := 0; i < c.Length(); i++ {
		c.syncIndex = i

		block := c.Blocks[i]
		c.Consensus.SetPreference(block.Data).SetUpdateHandler(func(preference int) {
			block.Data = preference
		}).SetRequestAnswerHandler(func(k int) []int {
			return c.onReqAnswer(i, k)
		}).Sync()
		log.Debug("Block ", i, ": finished=", c.Consensus.Finished, " preference=", block.Data)
		if !c.Consensus.Finished {
			finished = false
			break
		}
	}
	c.Finished = finished
	c.Running = false
}
