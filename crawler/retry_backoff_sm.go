package crawler

import (
	"sync"
	"time"
	"log"
)

// RetryBackoffState represents the states in the retry/backoff state machine.
type RetryBackoffState int

const (
	Initial RetryBackoffState = iota
	Retry   RetryBackoffState = iota
	Backoff RetryBackoffState = iota
	Failed  RetryBackoffState = iota
)

type EventType string

const (
	RetryScheduled EventType = "retry_scheduled"
	BackoffScheduled EventType = "backoff_scheduled"
	Failure          EventType = "failure"
)

type Event struct {
	Time    time.Time
	Type    EventType
	Details string
}

type RetryBackoffSM struct {
	state RetryBackoffState
	cfg   RetryBackoffConfig

	mu        sync.Mutex
	nextRetry time.Time
	eventChan chan Event
}

// NewRetryBackoffSM creates a new retry/backoff state machine instance.
func NewRetryBackoffSM(cfg RetryBackoffConfig, eventChan chan Event) *RetryBackoffSM {
	return &RetryBackoffSM{
		state:      Initial,
		cfg:        cfg,
		eventChan:  eventChan,
	}
}

// CanRetry checks if the state machine can retry.
func (s *RetryBackoffSM) CanRetry() bool {
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
func (s *RetryBackoffSM) Next(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var e Event

	switch s.state {
	case Initial:
		s.state = Retry
		s.nextRetry = time.Now().Add(s.cfg.InitialRetryDelay)
		e = Event{Time: s.nextRetry, Type: RetryScheduled, Details: "Initial retry scheduled"}
	case Retry:
		if err == nil {
			s.state = Initial
			s.nextRetry = time.Time{}
			e = Event{Time: time.Now(), Type: RetryScheduled, Details: "Retry successful, reset state"}
		} else {
			s.state = Backoff
			s.nextRetry = time.Now().Add(s.calculateBackoff())
			e = Event{Time: s.nextRetry, Type: BackoffScheduled, Details: "Backoff scheduled due to error"}
		}
	case Backoff:
		if err == nil {
			s.state = Retry
			s.nextRetry = time.Now().Add(s.cfg.RetryDelayAfterBackoff)
			e = Event{Time: s.nextRetry, Type: RetryScheduled, Details: "Retry scheduled after successful backoff"}
		} else {
			s.state = Failed
			s.nextRetry = time.Time{}
			e = Event{Time: time.Now(), Type: Failure, Details: "Failure, will not retry"}
		}
	default:
		panic("unknown state")
	}

	if e.Type != "" {
		s.eventChan <- e
		log.Printf("Emitted event: %+v", e)
	}
}

// calculateBackoff calculates the next backoff interval.
func (s *RetryBackoffSM) calculateBackoff() time.Duration {
	// Implement exponential backoff logic here.
	// Example: return time.Duration(float64(s.cfg.BackoffFactor) * float64(s.cfg.MaxBackoff))
	return 0
}
