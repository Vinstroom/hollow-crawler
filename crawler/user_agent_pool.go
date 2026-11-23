package crawler

import (
	"math/rand"
	"sync"
	"time"
)

// UserAgentPool represents a pool of user agents for rotation.
type UserAgentPool struct {
	agents      []string
	mu          sync.Mutex
	rateLimiter *RateLimiter
}

// NewUserAgentPool creates a new UserAgentPool with the given user agents
// and rate limiter for controlling the frequency of agent rotation.
func NewUserAgentPool(userAgents []string, rateLimiter *RateLimiter) *UserAgentPool {
	return &UserAgentPool{
		agents:      userAgents,
		rateLimiter: rateLimiter,
	}
}

// Get retrieves a user agent from the pool, respecting the rate limit.
func (uap *UserAgentPool) Get() string {
	uap.mu.Lock()
	defer uap.mu.Unlock()

	if !uap.rateLimiter.Allow() {
		return uap.agents[0]
	}

	randomIndex := rand.Intn(len(uap.agents))
	return uap.agents[randomIndex]
}

// Add adds a new user agent to the pool.
func (uap *UserAgentPool) Add(agent string) {
	uap.mu.Lock()
	defer uap.mu.Unlock()

	uap.agents = append(uap.agents, agent)
}

// Remove removes a user agent from the pool.
func (uap *UserAgentPool) Remove(agent string) {
	uap.mu.Lock()
	defer uap.mu.Unlock()

	for i, a := range uap.agents {
		if a == agent {
			uap.agents = append(uap.agents[:i], uap.agents[i+1:]...)
			break
		}
	}
}
