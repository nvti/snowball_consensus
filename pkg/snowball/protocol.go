package snowball

import (
	"errors"
)

type Config struct {
	// K sample size
	K int

	// Alpha quorum size
	Alpha int

	// Beta decision threshold
	Beta int
}

type Consensus struct {
	config               Config
	clients              []*Client
	onDecided            func([]byte)
	preference           []byte
	consecutiveSuccesses int
	running              bool
}

func New(clients []*Client, preference []byte, onDecided func([]byte), config Config) (consensus *Consensus, err error) {
	if config.K < 1 || config.K > len(clients) {
		return nil, errors.New("k must be between 1 and n")
	}

	if config.Alpha < 1 || config.Alpha > config.K {
		return nil, errors.New("alpha must be between 1 and k")
	}

	if config.Beta < 1 {
		return nil, errors.New("beta must be great or equal by 1")
	}

	consensus = &Consensus{
		config:               config,
		clients:              clients,
		onDecided:            onDecided,
		preference:           preference,
		consecutiveSuccesses: 0,
	}

	return consensus, nil
}

func (c *Consensus) Start() {
	if c.running {
		return
	}

	c.running = true
	go func() {
		for c.consecutiveSuccesses >= c.config.Beta {
			// choose random k client from c.clients

			// get k answer

			// check if have more a answer with same response
			// if
			// newPreference := []byte{}

			// if reflect.DeepEqual(newPreference, c.preference) {
			// 	c.consecutiveSuccesses++
			// } else {
			// 	c.consecutiveSuccesses = 1
			// }
		}
		c.onDecided(c.preference)
	}()
}

// func mostFrequent(arr [][]byte) (int, []byte) {
// 	m := map[[32]byte]int{}
// 	var maxCount int
// 	var freq int
// 	for _, a := range arr {
// 		m[a]++
// 		if m[a] > maxCnt {
// 			maxCnt = m[a]
// 			freq = a
// 		}
// 	}

// 	return freq
// }
