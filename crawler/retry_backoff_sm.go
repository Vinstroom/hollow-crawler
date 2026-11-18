package crawler

import (
	"sync"
	"time"
)

// RetryBackoffState represents the states in the retry/backoff state machine.
type RetryBackoffState int

const (
	Initial RetryBackoffState = iota
	Retry   RetryBackoffState = iota
	Backoff RetryBackoffState = iota
	Failed  RetryBackoffState = iota
)

type retryBackoffSM struct {
	state RetryBackoffState
	cfg   RetryBackoffConfig

	mu        sync.Mutex
	nextRetry time.Time
}

// NewRetryBackoffSM creates a new retry/backoff state machine instance.
func NewRetryBackoffSM(cfg RetryBackoffConfig) *retryBackoffSM {
	return &retryBackoffSM{
		state: Initial,
		cfg:   cfg,
	}
}

// CanRetry checks if the state machine can retry.
func (s *retryBackoffSM) CanRetry() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch s.state {
	case Initial, Failed:
		return true
	case Retry:
		return time.Now().After(s.nextRetry)
	case Backoff:
		return false
	default:
		panic("unknown state")
	}
}

// Next tries to perform a retry or backoff, updating the state machine.
func (s *retryBackoffSM) Next(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch s.state {
	case Initial:
		s.state = Retry
		s.nextRetry = time.Now().Add(s.cfg.InitialRetryDelay)
	case Retry:
		if err == nil {
			s.state = Initial
			s.nextRetry = time.Time{}
		} else {
			s.state = Backoff
			s.nextRetry = time.Now().Add(s.calculateBackoff())
		}
	case Backoff:
		if err == nil {
			s.state = Retry
			s.nextRetry = time.Now().Add(s.cfg.RetryDelayAfterBackoff)
		} else {
			s.state = Failed
			s.nextRetry = time.Time{}
		}
	default:
		panic("unknown state")
	}
}

// calculateBackoff calculates the next backoff interval.
func (s *retryBackoffSM) calculateBackoff() time.Duration {
	// Implement exponential backoff logic here.
	// Example: return time.Duration(float64(s.cfg.BackoffFactor) * float64(s.cfg.MaxBackoff))
	return 0
}
