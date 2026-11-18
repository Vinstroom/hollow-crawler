package crawler

import (
	"sync"
	"time"
)

// ThroughputRateCalc tracks the rate of processed items per second.
// It provides a sliding window calculation based on item counts over time.
type ThroughputRateCalc struct {
	mu       sync.RWMutex
	window   [2]struct {
		count int64
		ts    int64
	}
	curIdx   int
	total    int64
}

// NewThroughputRateCalc creates a new ThroughputRateCalc with an initial zero state.
func NewThroughputRateCalc() *ThroughputRateCalc {
	trc := new(ThroughputRateCalc)
	trc.mu = sync.NewRWMutex()
	trc.window[0].ts = time.Now().UnixNano()
	trc.window[1].ts = trc.window[0].ts
	return trc
}

// Increment records one processed item at the current time.
func (trc *ThroughputRateCalc) Increment() {
	trc.mu.Lock()
	defer trc.mu.Unlock()

	now := time.Now().UnixNano()
	if now - trc.window[trc.curIdx].ts >= 1_000_000_000 {
		trc.curIdx ^= 1
		trc.window[trc.curIdx].count = 0
		trc.window[trc.curIdx].ts = now
	}
	trc.window[trc.curIdx].count++
	trc.total++
}

// Rate returns the throughput in items per second over the sliding window.
// It returns 0 if no time has elapsed.
func (trc *ThroughputRateCalc) Rate() float64 {
	trc.mu.RLock()
	defer trc.mu.RUnlock()

	now := time.Now().UnixNano()
	elapsed := now - trc.window[trc.curIdx].ts
	if elapsed <= 0 {
		return 0.0
	}
	count := trc.window[trc.curIdx].count
	return (float64(count) * 1_000_000_000.0) / float64(elapsed)
}

// Total returns the total number of items recorded.
func (trc *ThroughputRateCalc) Total() int64 {
	trc.mu.RLock()
	defer trc.mu.RUnlock()
	return trc.total
}
