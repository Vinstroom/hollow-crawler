package crawler

import (
	"sync"
	"time"
)

// RateCalculator calculates the rate of HTML extraction.
type RateCalculator struct {
	startTime time.Time
	extractionCount int
	mutex       sync.Mutex
}

// NewRateCalculator returns a new RateCalculator.
func NewRateCalculator() *RateCalculator {
	return &RateCalculator{
		startTime: time.Now(),
	}
}

// IncrementExtractionCount increases the extraction count and returns the current rate.
func (rc *RateCalculator) IncrementExtractionCount() float64 {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()

	rc.extractionCount++
	elapsedTime := time.Since(rc.startTime).Seconds()
	return float64(rc.extractionCount) / elapsedTime
}
