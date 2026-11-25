package crawler

import (
	"math"
	"time"
)

// retryBackoffConfig defines the configuration for retry/backoff.
type retryBackoffConfig struct {
	initialDelay time.Duration
	maxDelay     time.Duration
	retryFactor  float64
	maxRetries   int
}

// retryBackoff controls the retry/backoff behavior for a given URL.
type retryBackoff struct {
	config   *retryBackoffConfig
	currentDelay time.Duration
	retries    int
}

// newRetryBackoff creates a new retryBackoff instance with the given config.
func newRetryBackoff(config *retryBackoffConfig) *retryBackoff {
	return &retryBackoff{
		config:   config,
		currentDelay: config.initialDelay,
		retries:    0,
	}
}

// RetryBackoff implements the RetryBackoff interface.
func (rb *retryBackoff) RetryBackoff() {
	rb.retries++
	rb.currentDelay = min(rb.config.maxDelay, math.Ceil(rb.currentDelay*rb.config.retryFactor))
	time.Sleep(rb.currentDelay)
}

// Reset resets the retryBackoff instance, readying it for the next use.
func (rb *retryBackoff) Reset() {
	rb.currentDelay = rb.config.initialDelay
	rb.retries = 0
}

// IsExhausted checks if the retryBackoff instance has exhausted its retries.
func (rb *retryBackoff) IsExhausted() bool {
	return rb.retries >= rb.config.maxRetries
}

func min(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}
