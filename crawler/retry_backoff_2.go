package crawler

import (
	"sync/atomic"
	"time"
	"encoding"
	"fmt"
)

// RetryBackoffConfig represents the configuration for retry/backoff.
type RetryBackoffConfig struct {
	MaxRetries     int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
	BackoffFactor  float64
	UserAgentPool  encoding.BinaryMarshaler // Encode and decode user-agent pool
}

// RetryBackoffState holds the state for retry/backoff.
type RetryBackoffState struct {
	retries       uint32
	backoff       time.Duration
	config        *RetryBackoffConfig
	shouldRetry   bool
	attemptsLeft atomic.Value // Use atomic.Value for safe access/modification
	userAgent    string // Store current selected user-agent
}

// NewRetryBackoffState creates a new RetryBackoffState with the given config.
func NewRetryBackoffState(config *RetryBackoffConfig) (*RetryBackoffState, error) {
	pool, err := decodeUserAgentPool(config.UserAgentPool)
	if err != nil {
		return nil, err
	}
	return &RetryBackoffState{
		config:        config,
		backoff:       config.InitialBackoff,
		shouldRetry:   true,
		attemptsLeft:  atomic.Value{},
		userAgent:     pool[0], // Initialize with first user-agent from pool
	}, nil
}

// decodeUserAgentPool decodes the user-agent pool from the provided encoding.BinaryMarshaler.
func decodeUserAgentPool(marshaler encoding.BinaryMarshaler) ([]string, error) {
	data, err := marshaler.MarshalBinary()
	if err != nil {
		return nil, err
	}
	var pool []string
	// Assuming the pool is marshalled as a comma-separated list of strings
	err = encoding.JSON.Unmarshal(data, &pool)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

// SelectNextUserAgent picks the next user-agent from the pool.
func (s *RetryBackoffState) SelectNextUserAgent() {
	pool, err := decodeUserAgentPool(s.config.UserAgentPool)
	if err != nil {
		// Handle error, e.g., log and continue with current user-agent
		fmt.Println("Error decoding user-agent pool: ", err)
		return
	}
	s.userAgent = pool[s.retries%len(pool)]
}

// UserAgent returns the current selected user-agent.
func (s *RetryBackoffState) UserAgent() string {
	return s.userAgent
}

// Retry should be called when a task needs to be retried.
func (s *RetryBackoffState) Retry() {
	s.retries++
	s.backoff = calculateBackoff(s.backoff, s.config.BackoffFactor)
	s.SelectNextUserAgent() // Select the next user-agent from the pool

	// Update attemptsLeft using atomic.Value
	s.attemptsLeft.Store(s.config.MaxRetries - int(s.retries))

	s.shouldRetry = s.attemptsLeft.Load().(int) > 0
}

// AttemptsLeft returns the number of attempts left for retries.
func (s *RetryBackoffState) AttemptsLeft() int {
	return s.attemptsLeft.Load().(int)
}

// WaitingTime returns the time to wait before the next retry.
func (s *RetryBackoffState) WaitingTime() time.Duration {
	if !s.shouldRetry {
		return 0
	}
	return s.backoff
}

// calculateBackoff calculates the next backoff delay.
func calculateBackoff(current, factor float64) time.Duration {
	return time.Duration(min(current*factor, float64(s.config.MaxBackoff)))
}

// min returns the minimum of two values.
func min(x, y float64) float64 {
	if x < y {
		return x
	}
	return y
}
