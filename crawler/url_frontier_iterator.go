package crawler

import (
	"context"
	"sync"
)

// URLFrontierIterator is a concurrent-safe iterator over URLs in the frontier.
// It allows sequential processing of URLs with deduplication checks and
// optional scoring-based prioritization. The iterator blocks until URLs are
// available or the context is canceled.
type URLFrontierIterator struct {
	frontier   *URLFrontier
	mu         sync.Mutex
	urls       []URL
	current    int
	scored     bool
	scorer     Scorer
	deduplicator Deduplicator
}

// NewURLFrontierIterator creates a new URLFrontierIterator for the given frontier.
// The iterator can optionally use a scorer for prioritization and a deduplicator
// for duplicate URL checks.
func NewURLFrontierIterator(frontier *URLFrontier, scorer Scorer, deduplicator Deduplicator) *URLFrontierIterator {
	return &URLFrontierIterator{
		frontier:   frontier,
		scored:     scorer != nil,
		scorer:     scorer,
		deduplicator: deduplicator,
	}
}

// Next returns the next URL in the frontier, blocking until one is available.
// If the iterator is scored, URLs are prioritized by their score.
// If deduplication is enabled, it checks for duplicates before returning.
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
		i.current++
		return url, true
	}
}

// fetchNextBatch loads a batch of URLs from the frontier into the iterator's buffer.
// Returns false if no more URLs are available.
func (i *URLFrontierIterator) fetchNextBatch(ctx context.Context) bool {
	if !i.scored {
		// Unscored: fetch all available URLs
		urls, err := i.frontier.GetAll(ctx)
		if err != nil {
			return false
		}
		i.urls = urls
		return len(i.urls) > 0
	} else {
		// Scored: fetch URLs in order of priority
		urls := make([]URL, 0, 100)
		for {
			// Fetch a batch of URLs (size depends on frontier implementation)
			batch, err := i.frontier.GetNext(ctx, 100)
			if err != nil || len(batch) == 0 {
				break
			}
			urls = append(urls, batch...)
		}

		// Sort URLs by score (descending)
		if len(urls) > 1 {
			scoreURLs := make([]scoreURL, len(urls))
			for j, url := range urls {
				scoreURLs[j] = scoreURL{url: url, score: i.scorer.Score(url)}
			}
			// Sort in descending order of score
			for i := 0; i < len(scoreURLs); i++ {
				for j := i + 1; j < len(scoreURLs); j++ {
					if scoreURLs[i].score < scoreURLs[j].score {
						scoreURLs[i], scoreURLs[j] = scoreURLs[j], scoreURLs[i]
					}
				}
			}
			for _, su := range scoreURLs {
				urls = append(urls, su.url)
			}
		}

		i.urls = urls
		return len(i.urls) > 0
	}
}

// Reset resets the iterator to its initial state, allowing it to be reused.
func (i *URLFrontierIterator) Reset() {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.urls = nil
	i.current = 0
}

// scoreURL is a helper struct for sorting URLs by score.
type scoreURL struct {
	url  URL
	score float64
}
