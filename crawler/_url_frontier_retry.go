// _url_frontier_retry.go

package crawler

import (
	"sync"
	"time"
)

// urlFrontierRetryer wraps urlFrontier to add retry logic with backoff.
type urlFrontierRetryer struct {
	urlFrontier
	backoff BackoffStrategy
	mu      sync.Mutex
}

// NewURLFrontierRetryer returns a new urlFrontierRetryer.
func NewURLFrontierRetryer(uf urlFrontier, b BackoffStrategy) *urlFrontierRetryer {
	return &urlFrontierRetryer{
		urlFrontier: uf,
		backoff:     b,
	}
}

// Get fetches a URL from the frontier with added retry and backoff logic.
func (r *urlFrontierRetryer) Get() (string, error) {
	for {
		url, err := r.urlFrontier.Get()
		if err == nil {
			return url, nil
		}

		r.mu.Lock()
		delay := r.backoff.NextBackoff()
		r.mu.Unlock()

		time.Sleep(delay)
		if !r.backoff.CanRetry(err) {
			return "", err
		}
	}
}
