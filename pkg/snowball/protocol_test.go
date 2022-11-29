package snowball

import (
	"math/rand"
	"snowball/pkg/utils"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestConsensus(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	type args struct {
		initPreference []int
		config         ConsensusConfig
	}
	tests := []struct {
		name           string
		args           args
		wantFinish     bool
		wantPreference int
	}{
		{
			name: "Normal",
			args: args{
				initPreference: []int{0, 1, 1, 2, 1, 1, 2, 2, 1, 0, 0},
				config: ConsensusConfig{
					K:       6,
					Alpha:   4,
					Beta:    10,
					MaxStep: 200,
				},
			},
			wantFinish:     true,
			wantPreference: 1,
		},
		{
			name: "Too much possible choices",
			args: args{
				initPreference: []int{0, 1, 3, 4, 5, 5, 1, 1, 3, 4, 0},
				config: ConsensusConfig{
					K:       6,
					Alpha:   4,
					Beta:    10,
					MaxStep: 200,
				},
			},
			wantFinish: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			consensuses := []*Consensus[int]{}
			// create clients
			for i, preference := range tt.args.initPreference {
				consensus := NewConsensus[int](ConsensusConfig{
					Name:    "Client " + strconv.Itoa(i),
					K:       tt.args.config.K,
					Alpha:   tt.args.config.Alpha,
					Beta:    tt.args.config.Beta,
					MaxStep: tt.args.config.MaxStep,
				}).SetPreference(preference)

				consensuses = append(consensuses, consensus)
			}

			for _, c := range consensuses {
				c.SetRequestAnswerHandler(func(k int) []int {
					clients := utils.GetRandomSubArray(consensuses, k)

					answers := []int{}
					for _, c := range clients {
						p, err := c.Preference()
						if err != nil {
							continue
						}
						answers = append(answers, p)
					}

					return answers
				})
			}

			// Start consensus
			for _, c := range consensuses {
				c.Sync()
			}

			// Wait all consensus finished
			wg := sync.WaitGroup{}
			allFinish := true
			for _, c := range consensuses {
				wg.Add(1)
				go func(c *Consensus[int]) {
					finished := <-c.Finished
					if !finished {
						allFinish = false
					}
					wg.Done()
				}(c)
			}
			wg.Wait()

			if tt.wantFinish != allFinish {
				t.Errorf("TestConsensus() finish = %v, want %v", allFinish, tt.wantFinish)
				return
			}

			if tt.wantFinish {
				// Make sure all clients have the same preference
				preference, _ := consensuses[0].Preference()
				for _, c := range consensuses {
					p, _ := c.Preference()
					if preference != p {
						t.Errorf("TestConsensus() preference is not same at all client")
						return
					}
				}

				if preference != tt.wantPreference {
					t.Errorf("TestConsensus() = %v, want %v", preference, tt.wantPreference)
					return
				}
			}
		})
	}
}
