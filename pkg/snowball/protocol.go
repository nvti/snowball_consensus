package snowball

import "errors"

type Config struct {
	// N number of participants
	N int

	// K sample size
	K int

	// Alpha quorum size
	Alpha int

	// Beta decision threshold
	Beta int
}

type Consensus struct {
	Config  Config
	Clients []*Client

	Preference           []byte
	ConsecutiveSuccesses int
}

func New(config Config, clients []*Client, preference []byte) (consensus *Consensus, err error) {
	if config.K < 1 || config.K > config.N {
		return nil, errors.New("k must be between 1 and n")
	}

	if config.Alpha < 1 || config.Alpha > config.K {
		return nil, errors.New("alpha must be between 1 and k")
	}

	if config.Beta < 1 {
		return nil, errors.New("beta must be great or equal by 1")
	}

	consensus = &Consensus{
		Config:               config,
		Clients:              clients,
		Preference:           preference,
		ConsecutiveSuccesses: 0,
	}

	return consensus, nil
}

func (c *Consensus) Start() {

}
