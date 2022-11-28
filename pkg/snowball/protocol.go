package snowball

import (
	"math/rand"
	"reflect"
	"snowball/pkg/log"
	"snowball/pkg/utils"
	"sync"
	"time"
)

type ConsensusConfig struct {
	// Name of consensus
	Name string

	// K sample size
	K int

	// Alpha quorum size
	Alpha int

	// Beta decision threshold
	Beta int

	// MaxStep Max number of consensus step
	// Set to 0 for running until consensus is reached
	MaxStep int
}

type OnUpdateHandler[T PreferenceType] func(T)

type Consensus[T PreferenceType] struct {
	config     ConsensusConfig
	clients    []Client[T]
	onUpdate   func(T)
	preference T
	confidence int
	running    bool
	Finished   chan bool
}

func NewConsensus[T PreferenceType](config ConsensusConfig) (consensus *Consensus[T]) {
	if config.Alpha < 1 || config.Alpha > config.K {
		log.Error("Alpha must be in range [1, K]")
		return nil
	}

	if config.Beta < 1 {
		log.Error("Beta must be greater than 1")
		return nil
	}

	consensus = &Consensus[T]{
		config:     config,
		onUpdate:   func(t T) {},
		confidence: 0,
		running:    false,
		Finished:   make(chan bool),
	}

	return consensus
}

func (c *Consensus[T]) SetClients(clients []Client[T]) *Consensus[T] {
	c.clients = clients
	return c
}

func (c *Consensus[T]) SetPreference(preference T) *Consensus[T] {
	c.preference = preference
	return c
}

func (c *Consensus[T]) SetUpdateHandler(handler OnUpdateHandler[T]) *Consensus[T] {
	c.onUpdate = handler
	return c
}

func (c *Consensus[T]) Preference() (T, error) {
	return c.preference, nil
}

func (c *Consensus[T]) StopSync() {
	c.running = false
}

func (c *Consensus[T]) keepRunning(step int) bool {
	return (c.config.MaxStep == 0 || step < c.config.MaxStep) && c.running
}

// Start the consensus
func (c *Consensus[T]) Sync() {
	if c.running {
		return
	}

	c.running = true
	go func() {
		i := 0
		for ; c.confidence < c.config.Beta; i++ {
			if !c.keepRunning(i) {
				break
			}

			log.Debug(c.config.Name, ": Step ", i)
			c.step()
		}
		finished := c.keepRunning(i)
		c.Finished <- finished
		log.Info(c.config.Name, ": Finnish after ", i, " step, finished = ", finished, ", preference = ", c.preference)
		c.running = false
	}()
}

func (c *Consensus[T]) step() {
	if len(c.clients) < c.config.K {
		// wait for other client join the network
		sleepTime := time.Duration(rand.Intn(1000))
		time.Sleep(sleepTime * time.Millisecond)
		return
	}

	// choose random k client from c.clients
	clients := utils.GetRandomSubArray(c.clients, c.config.K)

	// get k answer
	answers := []T{}
	// get answer in parallel
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	for _, client := range clients {
		wg.Add(1)
		go func(client Client[T]) {
			mu.Lock()
			preference, err := client.Preference()
			if err == nil {
				answers = append(answers, preference)
			} else {
				log.Error("Error when get preference from client: ", err)
			}
			mu.Unlock()
			wg.Done()
		}(client)
	}
	wg.Wait()

	// check if there is a majority
	count, preference, err := utils.MostFrequent(answers)
	if err != nil {
		log.Error(err)
		return
	}

	if count >= c.config.Alpha {
		oldPreference := c.preference
		c.preference = preference

		c.onUpdate(c.preference)

		if reflect.DeepEqual(oldPreference, c.preference) {
			c.confidence++
			log.Debug(c.config.Name, ": Got same preference, confidence = ", c.confidence)
		} else {
			log.Debug(c.config.Name, ": Got new preference")
			c.confidence = 1
		}
	} else {
		log.Debug(c.config.Name, ": There are no major answer, most answer is ", preference, ", count = ", count)
		c.confidence = 0
	}
}
