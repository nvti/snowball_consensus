package chain

import (
	"math/rand"
	"snowball/pkg/log"
	"snowball/pkg/snowball"
	"snowball/pkg/utils"
	"strconv"
	"sync"
	"testing"
	"time"
)

func createTestConsensusChain(name string, chains []*SnowballChain) *SnowballChain {
	chain := NewConsensusChain(snowball.ConsensusConfig{
		Name:    name,
		K:       20,
		Alpha:   10,
		Beta:    10,
		MaxStep: 0,
	})

	// have 10 possible choices
	for i := 0; i < 3; i++ {
		data := rand.Intn(10)
		_ = chain.Add(data)
	}

	chain.SetRequestAnswerHandler(func(index int, k int) []int {
		clients := utils.GetRandomSubArray(chains, k)

		answers := []int{}
		for _, c := range clients {
			if c != nil {
				p, err := c.Preference(index)
				if err != nil {
					continue
				}
				answers = append(answers, p)
			}
		}

		return answers
	})

	return chain
}

func TestConsensusChain_Sync(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	t.Run("200 client", func(t *testing.T) {
		chains := make([]*SnowballChain, 200)
		for i := 0; i < 200; i++ {
			chain := createTestConsensusChain("Client "+strconv.Itoa(i), chains)
			chains[i] = chain
		}

		wg := sync.WaitGroup{}
		for i, chain := range chains {
			wg.Add(1)
			log.Info("Client ", i, " started")

			go func(chain *SnowballChain, i int) {
				chain.Sync()
				log.Info("Client ", i, " finished=", chain.Finished)
				wg.Done()
			}(chain, i)
		}
		wg.Wait()

		// check if all client have the same data
		chain := chains[0]
		for i := 1; i < len(chains); i++ {
			chain2 := chains[i]
			if chain.Length() != chain2.Length() {
				t.Fatal("Chain length not equal")
			}

			for k := 0; k < chain.Length(); k++ {
				block, _ := chain.Get(k)
				block2, _ := chain2.Get(k)

				if block.Data != block2.Data {
					t.Fatal("Data not equal")
				}
			}
		}
	})
}

func TestConsensusChain_Sync2(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	t.Run("150 client + 50 client", func(t *testing.T) {
		chains := make([]*SnowballChain, 200)
		for i := 0; i < 150; i++ {
			chain := createTestConsensusChain("Client "+strconv.Itoa(i), chains)
			chains[i] = chain
		}

		wg := sync.WaitGroup{}
		for i := 0; i < 150; i++ {
			wg.Add(1)
			log.Info("Client ", i, " started")
			chain := chains[i]

			go func(chain *SnowballChain, i int) {
				chain.Sync()
				log.Info("Client ", i, " finished=", chain.Finished)
				wg.Done()
			}(chain, i)
		}

		// last 50 client will sync after 5 seconds
		time.Sleep(5 * time.Second)
		for i := 150; i < 200; i++ {
			chain := createTestConsensusChain("Client "+strconv.Itoa(i), chains)
			chains[i] = chain

			wg.Add(1)
			log.Info("Client ", i, " started")

			go func(chain *SnowballChain, i int) {
				chain.Sync()
				log.Info("Client ", i, " finished=", chain.Finished)
				wg.Done()
			}(chain, i)
		}

		wg.Wait()

		// check if all client have the same data
		chain := chains[0]
		for i := 1; i < len(chains); i++ {
			chain2 := chains[i]
			if chain.Length() != chain2.Length() {
				t.Fatal("Chain length not equal")
			}

			for k := 0; k < chain.Length(); k++ {
				block, _ := chain.Get(k)
				block2, _ := chain2.Get(k)

				if block.Data != block2.Data {
					t.Fatal("Data not equal")
				}
			}
		}
	})
}
