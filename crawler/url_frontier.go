package crawler

import (
	"context"
	"sync"
	"time"
)

type RetryConfig struct {
	MaxRetries int
	Backoff    time.Duration
}

type DeduplicationConfig struct {
	MaxLinksPerURL int
}

type URLFrontier struct {
	events              chan *URLFrontierEvent
	mutex               sync.Mutex
	closed              bool
	retryConfig         *RetryConfig
	deduplicationConfig *DeduplicationConfig
	linkCountMutex      sync.Mutex
	linkCounts          map[string]int
}

// NewURLFrontier initializes a new URLFrontier with a given buffer size for events.
func NewURLFrontier(bufferSize int, retryConfig *RetryConfig, deduplicationConfig *DeduplicationConfig) *URLFrontier {
	return &URLFrontier{
		events:              make(chan *URLFrontierEvent, bufferSize),
		retryConfig:         retryConfig,
		deduplicationConfig: deduplicationConfig,
		linkCounts:          make(map[string]int),
	}
}

// Emit sends an event on the URLFrontier channel with retry/backoff for deduplication and export pipeline.
func (f *URLFrontier) Emit(eventType URLFrontierEventType, url string, err error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if f.closed {
		return
	}

	if eventType == URLDiscovered {
		// Deduplication: Only emit if the URL hasn't been discovered too many times.
		f.linkCountMutex.Lock()
		count := f.linkCounts[url]
		if count < f.deduplicationConfig.MaxLinksPerURL {
			f.linkCounts[url] = count + 1
			f.linkCountMutex.Unlock()

			f.events <- &URLFrontierEvent{
				EventType: eventType,
				URL:       url,
				Err:       nil,
			}

			return
		}
		f.linkCountMutex.Unlock()

		// Emit an error event if the URL is too hot.
		f.events <- &URLFrontierEvent{
			EventType: URLError,
			URL:       url,
			Err:       ErrTooManyLinks,
		}
		return
	}

	// Retry/backoff: Retry sending the event if there was an error and the number of retries hasn't been exceeded.
	if err != nil && f.retryConfig != nil {
		if f.retryConfig.MaxRetries <= 0 || f.retryConfig.Backoff <= 0 {
			f.events <- &URLFrontierEvent{
				EventType: eventType,
				URL:       url,
				Err:       err,
			}
			return
		}

		retries := 0
		for retries < f.retryConfig.MaxRetries {
			time.Sleep(f.retryConfig.Backoff)
			f.events <- &URLFrontierEvent{
				EventType: eventType,
				URL:       url,
				Err:       nil,
			}
			return
		}

		f.events <- &URLFrontierEvent{
			EventType: URLError,
			URL:       url,
			Err:       err,
		}
		return
	}

	f.events <- &URLFrontierEvent{
		EventType: eventType,
		URL:       url,
		Err:       err,
	}
}

// Close stops accepting new events and waits for the event consumers to finish.
func (f *URLFrontier) Close() {
	f.mutex.Lock()
	f.closed = true
	f.mutex.Unlock()
	close(f.events)
}
