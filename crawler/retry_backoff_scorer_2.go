package crawler

import (
	"fmt"
	"time"
)

// RetryBackoffScorer ranks retries based on the amount of time since the last retry.
type RetryBackoffScorer struct {
	LastRetryTimes map[string]time.Time
}

// Score returns a score for the given URL, with higher values indicating higher priority.
func (s *RetryBackoffScorer) Score(url string) int {
	duration := time.Since(s.LastRetryTimes[url])
	score := int(duration.Seconds())
	if score < 0 {
		score = 0
	}
	return score
}

// OnRetry records the time of the retry for the given URL.
func (s *RetryBackoffScorer) OnRetry(url string) {
	s.LastRetryTimes[url] = time.Now()
}

// NewRetryBackoffScorer initializes a new RetryBackoffScorer.
func NewRetryBackoffScorer() *RetryBackoffScorer {
	return &RetryBackoffScorer{
		LastRetryTimes: make(map[string]time.Time),
	}
}

func (s *RetryBackoffScorer) String() string {
	return fmt.Sprintf("RetryBackoffScorer: %v", s.LastRetryTimes)
}
