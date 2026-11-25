package crawler

import (
	"context"
	"sync"

	"github.com/hollow-crawler/pkg/events"
)

// URLFrontierIterator is a concurrent-safe iterator over URLs in the frontier.
// It allows sequential processing of URLs with deduplication checks, optional
// scoring-based prioritization, and now also emits events for proxy rotations.
type URLFrontierIterator struct {
	frontier   *URLFrontier
	mu         sync.Mutex
	urls       []URL
	current    int
	scored     bool
	scorer     Scorer
	deduplicator Deduplicator
	proxyRotator ProxyRotator
	eventEmitter *events.EventEmitter
}

// NewURLFrontierIterator creates a new URLFrontierIterator for the given frontier.
// The iterator can optionally use a scorer for prioritization, a deduplicator for
// duplicate URL checks, and a proxyRotator for rotating proxies. An event emitter
// is also provided for emitting proxy rotation events.
func NewURLFrontierIterator(
	frontier *URLFrontier,
	scorer Scorer,
	deduplicator Deduplicator,
	proxyRotator ProxyRotator,
	eventEmitter *events.EventEmitter,
) *URLFrontierIterator {
	return &URLFrontierIterator{
		frontier:      frontier,
		scored:        scorer != nil,
		scorer:        scorer,
		deduplicator:  deduplicator,
		proxyRotator:  proxyRotator,
		eventEmitter:  eventEmitter,
	}
}

// Next returns the next URL in the frontier, blocking until one is available.
// If the iterator is scored, URLs are prioritized by their score.
// If deduplication is enabled, it checks for duplicates before returning.
// If proxy rotation is enabled, it rotates the proxy before returning the URL.
// Returns false if the context is canceled or no more URLs are available.
func (i *URLFrontierIterator) Next(ctx context.Context) (URL, bool) {
	select {
	case <-ctx.Done():
		return URL{}, false
	default:
		i.mu.Lock()
		defer i.mu.Unlock()

		// If no URLs are currently buffered, fetch more from the frontier
		if i.current >= len(i.urls) {
			if !i.fetchNextBatch(ctx) {
				return URL{}, false
			}
		}

		// Check for duplicates if deduplicator is provided
		if i.deduplicator != nil {
			url := i.urls[i.current]
			if !i.deduplicator.IsUnique(ctx, url) {
				i.current++
				return i.Next(ctx)
			}
		}

		url := i.urls[i.current]

		// Rotate proxy if proxyRotator is provided
		if i.proxyRotator != nil {
			if proxy, err := i.proxyRotator.Rotate(ctx); err == nil {
				url.Proxy = proxy
				i.eventEmitter.Emit(events.ProxyRotated, proxy)
			}
		}

		i.current++
		return url, true
	}
}

// ... (Other methods remain the same)
