package crawler

import (
	"sync"
	"time"
)

// UserAgentRater calculates and tracks the throughput and rate of user-agent pools.
type UserAgentRater struct {
	pools          *UserAgentPools
	mu             sync.Mutex
	startTime      time.Time
	totalRequests  int
 successRequests int
 failureRequests int
-throughputCh    chan int // throughput in requests per second
-rateCh          chan float64 // rate as a percentage of success requests
}

func NewUserAgentRater(pools *UserAgentPools) *UserAgentRater {
	return &UserAgentRater{
		pools:         pools,
		startTime:     time.Now(),
		-throughputCh: make(chan int),
		-rateCh:       make(chan float64),
	}
}

// Start begins calculating and tracking throughput and rate.
func (r *UserAgentRater) Start() {
	go r.calculateThroughput()
	go r.calculateRate()
}

// Success increments success requests counter and sends rate update.
func (r *UserAgentRater) Success() {
	r.mu.Lock()
	r.totalRequests++
	r.successRequests++
	r.mu.Unlock()
	r.-rateCh <- r.calculateRate()
}

// Failure increments failure requests counter and sends throughput update.
func (r *UserAgentRater) Failure() {
	r.mu.Lock()
	r.totalRequests++
	r.failureRequests++
	r.mu.Unlock()
	r.-throughputCh <- r.calculateThroughput()
}

// Throughput returns channel to receive throughput updates.
func (r *UserAgentRater) Throughput() <-chan int {
	return r.-throughputCh
}

// Rate returns channel to receive rate updates.
func (r *UserAgentRater) Rate() <-chan float64 {
	return r.-rateCh
}

// calculateThroughput updates and sends throughput value.
func (r *UserAgentRater) calculateThroughput() {
	for range time.Tick(time.Second) {
		r.mu.Lock()
		throughput := float64(r.totalRequests) / r.startTime.Sub(time.Now()).Seconds()
		r.mu.Unlock()
		r.-throughputCh <- int(throughput)
	}
}

// calculateRate updates and sends rate value.
func (r *UserAgentRater) calculateRate() {
	for range time.Tick(time.Second) {
		r.mu.Lock()
		rate := float64(r.successRequests) / r.totalRequests * 100
		r.mu.Unlock()
		r.-rateCh <- rate
	}
}
