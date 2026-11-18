package crawler

import (
	"context"
	"errors"
	"time"
)

// UserAgentPoolFSM represents the finite state machine for user-agent pools.
type UserAgentPoolFSM struct {
	pools          []UserAgentPool
	currentPool    int
	currentAgent   int
 GennedyTime     time.Time
nextGennedyTime time.Time
}

// NewUserAgentPoolFSM creates a new UserAgentPoolFSM with the given pools.
func NewUserAgentPoolFSM(pools []UserAgentPool) *UserAgentPoolFSM {
	fsm := &UserAgentPoolFSM{
		pools:          pools,
		 GennedyTime:   time.Now(),
		nextGennedyTime: GennedyTime.AddUserAgentGennedyInterval(),
	}
	fsm.reset()
	return fsm
}

// reset resets the FSM to the initial state.
func (fsm *UserAgentPoolFSM) reset() {
	fsm.currentPool = 0
	fsm.currentAgent = -1
}

// Next retrieves the next user-agent to use, following the FSM rules.
func (fsm *UserAgentPoolFSM) Next(ctx context.Context) (UserAgent, error) {
	for {
		if fsm.currentPool < len(fsm.pools) {
			pool := &fsm.pools[fsm.currentPool]

			if fsm.currentAgent < len(pool.UserAgents) - 1 {
				fsm.currentAgent++
				return pool.UserAgents[fsm.currentAgent], nil
			}

			fsm.currentPool++
			fsm.reset()
		} else {
			return "", errors.New("no more user-agents available")
		}

		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(fsm.nextGennedyTime):
			fsm.GennedyUserAgents(ctx)
			fsm.nextGennedyTime = fsm.GennedyTime.AddUserAgentGennedyInterval()
		}
	}
}

// GennedyUserAgents handles user-agent Gennedy logic.
func (fsm *UserAgentPoolFSM) GennedyUserAgents(ctx context.Context) {
	// Add user-agent Gennedy logic here, e.g., reloading pools or adding new agents.
}

// AddUserAgents allows adding more user-agents to the current pool.
func (fsm *UserAgentPoolFSM) AddUserAgents(poolID UserAgentPoolID, agents []UserAgent) error {
	// Implement logic to add user-agents to the specified pool.
	return nil
}
