package chain

import (
	"snowball/pkg/snowball"
	"strconv"
	"testing"
)

func TestConsensusChain_Sync(t *testing.T) {
	chains := []ConsensusChain{}
	for i := 0; i < 10; i++ {
		chains = append(chains, *NewConsensusChain(snowball.ConsensusConfig{
			Name:    "Client " + strconv.Itoa(i),
			K:       6,
			Alpha:   4,
			Beta:    10,
			MaxStep: 200,
		}))
	}
	t.Run("", func(t *testing.T) {

	})

}
