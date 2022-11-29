package snowball

import (
	"reflect"
	"snowball/pkg/log"
	"snowball/pkg/utils"
)

type PreferenceType comparable

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
type OnRequestAnswerHandler[T PreferenceType] func(int) []T

type Consensus[T PreferenceType] struct {
	config      ConsensusConfig
	onUpdate    OnUpdateHandler[T]
	onReqAnswer OnRequestAnswerHandler[T]
	preference  T
	confidence  int
	Running     bool
	Finished    bool
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
		Running:    false,
		Finished:   false,
	}

	return consensus
}

func (c *Consensus[T]) SetPreference(preference T) *Consensus[T] {
	c.preference = preference
	return c
}

func (c *Consensus[T]) SetUpdateHandler(handler OnUpdateHandler[T]) *Consensus[T] {
	c.onUpdate = handler
	return c
}

func (c *Consensus[T]) SetRequestAnswerHandler(handler OnRequestAnswerHandler[T]) *Consensus[T] {
	c.onReqAnswer = handler
	return c
}

func (c *Consensus[T]) Preference() (T, error) {
	return c.preference, nil
}

func (c *Consensus[T]) StopSync() {
	c.Running = false
}

func (c *Consensus[T]) keepRunning(step int) bool {
	return (c.config.MaxStep == 0 || step < c.config.MaxStep) && c.Running
}

// Start the consensus
func (c *Consensus[T]) Sync() {
	if c.Running {
		return
	}
	c.Running = true
	c.confidence = 1
	i := 0
	for ; c.confidence < c.config.Beta; i++ {
		if !c.keepRunning(i) {
			break
		}

		log.Debug(c.config.Name, ": Step ", i)
		c.step()
	}
	c.Finished = c.keepRunning(i)
	log.Debug(c.config.Name, ": Finnish after ", i, " step, finished = ", c.Finished, ", preference = ", c.preference)
	c.Running = false
}

func (c *Consensus[T]) step() {
	// get k answer
	answers := c.onReqAnswer(c.config.K)
	if len(answers) < c.config.K {
		return
	}
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
