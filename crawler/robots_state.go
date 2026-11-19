package crawler

import (
	"sync"
	"time"
)

// RobotsState tracks per-host politeness state for robots.txt compliance.
type RobotsState struct {
	mu            sync.Mutex
	lastFetches   map[string][]time.Time // Use slice to store last N fetches for better average calculation
	crawlDelays   map[string]time.Duration
}

func NewRobotsState() *RobotsState {
	return &RobotsState{
		lastFetches:  make(map[string][]time.Time),
		crawlDelays:  make(map[string]time.Duration),
	}
}

// SetCrawlDelay sets the crawl delay for a given host.
func (rs *RobotsState) SetCrawlDelay(host string, delay time.Duration) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.crawlDelays[host] = delay
}

// CrawlDelay returns the crawl delay for a host, defaulting to zero.
func (rs *RobotsState) CrawlDelay(host string) time.Duration {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	return rs.crawlDelays[host]
}

// RecordFetch records that a fetch happened for host at time t.
func (rs *RobotsState) RecordFetch(host string, t time.Time) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	if len(rs.lastFetches[host]) >= 5 {
		rs.lastFetches[host] = rs.lastFetches[host][1:]
	}
	rs.lastFetches[host] = append(rs.lastFetches[host], t)
}

// LastFetch returns the last fetch times for host, or an empty slice.
func (rs *RobotsState) LastFetch(host string) []time.Time {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	return rs.lastFetches[host]
}

// WaitDuration calculates how long to wait before the next allowed fetch for host.
// Uses the average delay from the last N fetches (default 5) to smooth out short bursts.
// Returns zero if immediate fetch is allowed.
func (rs *RobotsState) WaitDuration(host string) time.Duration {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	delays, ok := rs.crawlDelays[host]
	if !ok {
		return 0
	}

	fetches, ok := rs.lastFetches[host]
	if !ok || len(fetches) < 2 {
		return delays
	}

	totalDelay := time.Duration(0)
	for i := 1; i < len(fetches); i++ {
		totalDelay += fetches[i].Sub(fetches[i-1])
	}

	avgDelay := totalDelay / time.Duration(len(fetches)-1)
	if time.Since(fetches[len(fetches)-1]) < avgDelay {
		return avgDelay - time.Since(fetches[len(fetches)-1])
	}

	return 0
}
