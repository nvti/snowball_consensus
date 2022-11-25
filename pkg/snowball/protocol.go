package snowball

import (
	"errors"
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
	// Set to 0 for run until got finished
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

func NewConsensus[T PreferenceType](preference T, config ConsensusConfig) (consensus *Consensus[T], err error) {
	if config.Alpha < 1 || config.Alpha > config.K {
		return nil, errors.New("alpha must be between 1 and k")
	}

	if config.Beta < 1 {
		return nil, errors.New("beta must be great or equal by 1")
	}

	consensus = &Consensus[T]{
		config:     config,
		onUpdate:   func(t T) {},
		preference: preference,
		confidence: 0,
		running:    false,
		Finished:   make(chan bool),
	}

	return consensus, nil
}

func (c *Consensus[T]) SetClients(clients []Client[T]) {
	c.clients = clients
}

func (c *Consensus[T]) SetUpdateHandler(handler OnUpdateHandler[T]) {
	c.onUpdate = handler
}

func (c *Consensus[T]) Preference() T {
	return c.preference
}

func (c *Consensus[T]) Start() {
	if c.running {
		return
	}

	c.running = true
	go func() {
		i := 0
		for ; c.confidence < c.config.Beta; i++ {
			if c.config.MaxStep > 0 && i >= c.config.MaxStep {
				break
			}

			log.Debug(c.config.Name, ": Step ", i)
			c.Step()
		}
		finished := c.config.MaxStep == 0 || i < c.config.MaxStep
		c.Finished <- finished
		log.Info(c.config.Name, ": Finnish after ", i, " step, finished = ", finished, ", preference = ", c.preference)
		c.running = false
	}()
}

func (c *Consensus[T]) Step() {
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
			answers = append(answers, client.Preference())
			mu.Unlock()
			wg.Done()
		}(client)
	}
	wg.Wait()

	// check if have more a answer with same response
	count, preference := utils.MostFrequent(answers)
	if count >= c.config.Alpha {
		oldPreference := c.preference
		c.preference = preference

		c.onUpdate(c.preference)

		if reflect.DeepEqual(oldPreference, c.preference) {
			c.confidence++
			log.Debug(c.config.Name, ": Got same preference, confidence = ", c.confidence)
		} else {
			log.Debug(c.config.Name, ": Got difference preference")
			c.confidence = 1
		}
	} else {
		log.Debug(c.config.Name, ": There are no major answer, most answer is ", preference, ", count = ", count)
		c.confidence = 0
	}
}
