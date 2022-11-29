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
	Finished    chan bool
}

func NewConsensusChain(config snowball.ConsensusConfig) *SnowballChain {
	return &SnowballChain{
		Consensus: snowball.NewConsensus[int](config),
		Finished:  make(chan bool),
	}
}

func (c *SnowballChain) Preference(index int) (int, error) {
	block, err := c.Get(index)
	if err != nil {
		return 0, err
	}

	return block.Data, nil
}

func (c *SnowballChain) Sync() {
	go func() {
		for i := 0; i < c.Length(); i++ {
			c.syncIndex = i

			block := c.Blocks[i]
			c.Consensus.SetPreference(block.Data).SetUpdateHandler(func(preference int) {
				block.Data = preference
			}).SetRequestAnswerHandler(func(k int) []int {
				return c.onReqAnswer(i, k)
			}).Sync()
			finished := <-c.Consensus.Finished
			log.Info("Block ", i, ": finished=", finished, " preference=", block.Data)
		}
		c.Finished <- true
	}()
}
