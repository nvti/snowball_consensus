package snowball

import (
	"errors"
	"reflect"
	"snowball/pkg/utils"
	"sync"
)

type Config struct {
	// K sample size
	K int

	// Alpha quorum size
	Alpha int

	// Beta decision threshold
	Beta int
}

type OnDecidedHandler[T PreferenceType] func(T)

type Consensus[T PreferenceType] struct {
	config     Config
	clients    []Client[T]
	onDecided  func(T)
	preference T
	confidence int
	running    bool
}

func New[T PreferenceType](clients []Client[T], preference T, onDecided OnDecidedHandler[T], config Config) (consensus *Consensus[T], err error) {
	if config.K < 1 || config.K > len(clients) {
		return nil, errors.New("k must be between 1 and n")
	}

	if config.Alpha < 1 || config.Alpha > config.K {
		return nil, errors.New("alpha must be between 1 and k")
	}

	if config.Beta < 1 {
		return nil, errors.New("beta must be great or equal by 1")
	}

	consensus = &Consensus[T]{
		config:     config,
		clients:    clients,
		onDecided:  onDecided,
		preference: preference,
		confidence: 0,
	}

	return consensus, nil
}

func (c *Consensus[T]) Start() {
	if c.running {
		return
	}

	c.running = true
	go func() {
		for c.confidence >= c.config.Beta {
			// choose random k client from c.clients
			clients := utils.GetRandomSubArray(c.clients, c.config.K)

			// get k answer
			answers := make([]T, c.config.K)
			// get answer in parallel
			wg := sync.WaitGroup{}
			for i, client := range clients {
				wg.Add(1)
				go func(i int, client Client[T]) {
					answers[i] = client.Preference()
					wg.Done()
				}(i, client)
			}
			wg.Wait()

			// check if have more a answer with same response
			count, preference := utils.MostFrequent(answers)
			if count >= c.config.Alpha {
				oldPreference := c.preference
				c.preference = preference

				if reflect.DeepEqual(oldPreference, c.preference) {
					c.confidence++
				} else {
					c.confidence = 1
				}
			} else {
				c.confidence = 0
			}
		}
		c.onDecided(c.preference)
		c.running = false
	}()
}
