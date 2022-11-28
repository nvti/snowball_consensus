package chain

import (
	"math/rand"
	"snowball/pkg/log"
	"snowball/pkg/snowball"
	"strconv"
	"sync"
	"testing"
)

func TestConsensusChain_Sync(t *testing.T) {
	t.Run("", func(t *testing.T) {
		chains := []*ConsensusChain{}
		clients := []Client[int]{}
		for i := 0; i < 10; i++ {
			chain := NewConsensusChain(snowball.ConsensusConfig{
				Name:    "Client " + strconv.Itoa(i),
				K:       6,
				Alpha:   4,
				Beta:    10,
				MaxStep: 100,
			})

			for i := 0; i < 5; i++ {
				data := rand.Intn(2)
				_ = chain.Add(data)
			}

			chains = append(chains, chain)
			clients = append(clients, chain)
		}

		for _, chain := range chains {
			chain.SetClients(clients)
		}

		wg := sync.WaitGroup{}
		for i, chain := range chains {
			wg.Add(1)
			log.Info("Client ", i, " started")
			chain.Sync()

			go func(chain *ConsensusChain, i int) {
				<-chain.Finished
				log.Info("Client ", i, " finished")
				wg.Done()
			}(chain, i)
		}
		wg.Wait()
	})

}
