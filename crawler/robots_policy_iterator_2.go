package crawler

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sync"
)

// PolicyIterator generates URLs according to robots.txt policy.
// It implements the iterator pattern for controlled URL generation.
type PolicyIterator struct {
	baseURL    string
	policy     *RobotsPolicy
	seen       map[string]bool
	mu         sync.Mutex
	currentURL string
}

// NewPolicyIterator creates a new PolicyIterator for a given base URL and policy.
// It initializes the iterator with the base URL and prepares to generate allowed URLs.
func NewPolicyIterator(baseURL string, policy *RobotsPolicy) (*PolicyIterator, error) {
	parsedBase, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	// Ensure the base URL is absolute
	if parsedBase.Scheme == "" {
		parsedBase.Scheme = "http"
	}

	return &PolicyIterator{
		baseURL:    parsedBase.String(),
		policy:     policy,
		seen:       make(map[string]bool),
		currentURL: parsedBase.String(),
	}, nil
}

// Next returns the next URL that is allowed by the robots policy.
// It returns nil if no more URLs are allowed or an error occurs.
func (it *PolicyIterator) Next(ctx context.Context) (string, error) {
	// Check if context is canceled
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	// Get the current URL and reset it
	current := it.currentURL
	it.mu.Lock()
	defer it.mu.Unlock()

	// If we've already processed this URL, skip it
	if it.seen[current] {
		return "", errors.New("url already processed")
	}

	// Mark as seen
	it.seen[current] = true

	// Check if the current URL is allowed
	if !it.policy.Allowed(current) {
		return "", errors.New("url disallowed by robots policy")
	}

	// Generate next possible URLs (simplified for demo - in practice would use
	// actual link discovery from pages)
	nextURL := it.generateNextURL(current)
	if nextURL == "" {
		return "", errors.New("no more URLs to process")
	}

	// Update current URL
	it.currentURL = nextURL
	return current, nil
}

// generateNextURL generates the next URL to process (simplified example).
// In a real implementation, this would parse the HTML of the current URL
// and extract all allowed links according to the robots policy.
func (it *PolicyIterator) generateNextURL(current string) string {
	// For demo purposes, we'll just append "/page2" to the base URL
	// In a real crawler, this would parse the HTML and extract links
	if current == it.baseURL {
		return fmt.Sprintf("%s/page2", it.baseURL)
	}
	return ""
}

// Reset resets the iterator to its initial state, allowing it to be reused.
func (it *PolicyIterator) Reset() {
	it.mu.Lock()
	defer it.mu.Unlock()

	it.seen = make(map[string]bool)
	it.currentURL = it.baseURL
}

// Close releases any resources held by the iterator.
func (it *PolicyIterator) Close() error {
	it.mu.Lock()
	defer it.mu.Unlock()

	// No resources to release in this simple implementation
	return nil
}
