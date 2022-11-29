package chain

import (
	"math/rand"
	"snowball/pkg/log"
	"snowball/pkg/snowball"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestConsensusChain_Sync(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	t.Run("200 client", func(t *testing.T) {
		chains := []*SnowballChain{}
		clients := []Client[int]{}
		for i := 0; i < 200; i++ {
			chain := NewConsensusChain(snowball.ConsensusConfig{
				Name:    "Client " + strconv.Itoa(i),
				K:       20,
				Alpha:   10,
				Beta:    10,
				MaxStep: 0,
			})

			for i := 0; i < 3; i++ {
				data := rand.Intn(10)
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

			go func(chain *SnowballChain, i int) {
				<-chain.Finished
				log.Info("Client ", i, " finished")
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
		chains := []*SnowballChain{}
		clients := []Client[int]{}
		for i := 0; i < 200; i++ {
			chain := NewConsensusChain(snowball.ConsensusConfig{
				Name:    "Client " + strconv.Itoa(i),
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

			chains = append(chains, chain)
			clients = append(clients, chain)
		}

		chains150 := chains[:150]
		clients150 := clients[:150]
		chains50 := chains[150:]

		// only first 150 client will sync
		for _, chain := range chains150 {
			chain.SetClients(clients150)
		}

		wg := sync.WaitGroup{}
		for i, chain := range chains150 {
			wg.Add(1)
			log.Info("Client ", i, " started")
			chain.Sync()

			go func(chain *SnowballChain, i int) {
				<-chain.Finished
				log.Info("Client ", i, " finished")
				wg.Done()
			}(chain, i)
		}

		// last 50 client will sync after 5 seconds
		time.Sleep(5 * time.Second)
		for _, chain := range chains {
			chain.SetClients(clients)
		}

		for i, chain := range chains50 {
			wg.Add(1)
			log.Info("Client ", i+150, " started")
			chain.Sync()

			go func(chain *SnowballChain, i int) {
				<-chain.Finished
				log.Info("Client ", i, " finished")
				wg.Done()
			}(chain, i+150)
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
