package crawler

import (
	"sync"
	"time"
)

// ProxyRotatorFSM implements a finite state machine for proxy rotation.
// Each proxy can be in one of three states: Healthy, Degraded, or Failed.
// Transitions are driven by success/failure counts and cooldown timers.
type ProxyRotatorFSM struct {
	mu          sync.Mutex
	states      map[string]proxyState // proxyURL -> state
	degradedTTL time.Duration         // time before degraded -> healthy
	failedTTL   time.Duration         // time before failed -> degraded
	maxFails    int                   // consecutive failures to enter failed
	successRateCalculator *successRateCalculator
}

type proxyState struct {
	state          int // 0=healthy, 1=degraded, 2=failed
	consecFails    int
	lastTransition time.Time
	succeededSinceLastFail int
}

const (
	StateHealthy  = 0
	StateDegraded = 1
	StateFailed   = 2
)

// NewProxyRotatorFSM creates a new FSM with given TTLs and max failures.
func NewProxyRotatorFSM(degradedTTL, failedTTL time.Duration, maxFails int) *ProxyRotatorFSM {
	return &ProxyRotatorFSM{
		mu:          sync.Mutex{},
		states:      map[string]proxyState{},
		degradedTTL: degradedTTL,
		failedTTL:   failedTTL,
		maxFails:    maxFails,
		successRateCalculator: newSuccessRateCalculator(),
	}
}

// RecordSuccess marks a proxy as healthy; resets failure count.
func (fsm *ProxyRotatorFSM) RecordSuccess(proxyURL string) {
	fsm.mu.Lock()
	defer fsm.mu.Unlock()

	ps, ok := fsm.states[proxyURL]
	if !ok {
		ps = proxyState{state: StateHealthy, consecFails: 0, lastTransition: time.Now()}
	}
	ps.state = StateHealthy
	ps.succeededSinceLastFail++
	fsm.states[proxyURL] = ps
	fsm.successRateCalculator.recordSuccess()
}

// RecordFailure increments failure count and transitions if threshold met.
func (fsm *ProxyRotatorFSM) RecordFailure(proxyURL string) {
	fsm.mu.Lock()
	defer fsm.mu.Unlock()

	now := time.Now()
	ps, ok := fsm.states[proxyURL]
	if !ok {
		ps = proxyState{state: StateHealthy, consecFails: 0, lastTransition: now}
	}
	ps.consecFails++
	if ps.consecFails >= fsm.maxFails {
		ps.state = StateFailed
		ps.lastTransition = now
		ps.succeededSinceLastFail = 0
	} else {
		ps.state = StateDegraded
		ps.lastTransition = now
	}
	fsm.states[proxyURL] = ps
	fsm.successRateCalculator.recordFailure()
}

// GetState returns the current state for a proxy after applying time-based transitions.
func (fsm *ProxyRotatorFSM) GetState(proxyURL string) int {
	fsm.mu.Lock()
	defer fsm.mu.Unlock()

	ps, ok := fsm.states[proxyURL]
	if !ok {
		return StateHealthy
	}
	now := time.Now()
	switch ps.state {
	case StateFailed:
		if now.Sub(ps.lastTransition) >= fsm.failedTTL {
			ps.state = StateDegraded
			ps.consecFails = 0
			ps.lastTransition = now
			fsm.states[proxyURL] = ps
		}
	case StateDegraded:
		if now.Sub(ps.lastTransition) >= fsm.degradedTTL {
			ps.state = StateHealthy
			ps.consecFails = 0
			ps.lastTransition = now
			fsm.states[proxyURL] = ps
		}
	}
	return ps.state
}

// IsHealthy returns true if the proxy is in Healthy state.
func (fsm *ProxyRotatorFSM) IsHealthy(proxyURL string) bool {
	return fsm.GetState(proxyURL) == StateHealthy
}

// SuccessRate returns the success rate of requests using this FSM.
func (fsm *ProxyRotatorFSM) SuccessRate() float64 {
	return fsm.successRateCalculator.rate()
}

// successRateCalculator tracks and calculates the success rate of requests.
type successRateCalculator struct {
	successes int
	total     int
}

func newSuccessRateCalculator() *successRateCalculator {
	return &successRateCalculator{}
}

func (s *successRateCalculator) recordSuccess() {
	s.successes++
	s.total++
}

func (s *successRateCalculator) recordFailure() {
	s.total++
}

func (s *successRateCalculator) rate() float64 {
	if s.total == 0 {
		return 0
	}
	return float64(s.successes) / float64(s.total)
}
