package crawler

import (
	"math/rand"
	"time"
)

// RetryBackoffPolicy defines a policy for retrying requests with exponential backoff
type RetryBackoffPolicy struct {
	maxRetries      int
	maxBackoffTime  time.Duration
	initialBackoff  time.Duration
}

// NewRetryBackoffPolicy creates a new RetryBackoffPolicy with the given parameters
func NewRetryBackoffPolicy(maxRetries int, maxBackoffTime, initialBackoff time.Duration) *RetryBackoffPolicy {
	return &RetryBackoffPolicy{
		maxRetries:      maxRetries,
		maxBackoffTime:  maxBackoffTime,
		initialBackoff:  initialBackoff,
	}
}

// ShouldRetry checks if the request failed with a retryable error and if the number of retries is less than the max retries
func (p *RetryBackoffPolicy) ShouldRetry(attempt, maxRetries int, err error) bool {
	// Implement logic to check if the error is retryable and if the number of retries is less than the max retries
}

// NextBackoff returns the duration to wait before retrying the request, based on the current number of retries
func (p *RetryBackoffPolicy) NextBackoff(attempt int) time.Duration {
	// Implement exponential backoff logic with jitter
}

// Reset resets the policy to its initial state
func (p *RetryBackoffPolicy) Reset() {
	// Implement logic to reset the policy
}
