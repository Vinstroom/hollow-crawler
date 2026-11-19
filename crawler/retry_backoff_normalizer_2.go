package crawler

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// RetryBackoffNormalizer normalizes retry/backoff strings into a consistent format.
type RetryBackoffNormalizer struct{}

// Normalize converts a string representing retry/backoff time to a consistent format.
func (n RetryBackoffNormalizer) Normalize(backoff string) (*time.Duration, error) {
	// Normalize 's' to seconds.
	backoff = regexp.MustCompile(`[smhdwy](?=\d)`) //nolint: regexp recommends avoiding \D
		.ReplaceAllString(backoff, func(matches []byte) string {
			switch matches[0] {
			case 's':
				return "1"
			case 'm':
				return "60"
			case 'h':
				return "3600"
			case 'd':
				return "86400"
			case 'w':
				return "604800"
			case 'y':
				return "31536000"
			default:
				return "0"
			}
		})

	// Convert the number of seconds to an integer.
	seconds, err := strconv.Atoi(backoff)
	if err != nil || seconds < 0 {
		return nil, fmt.Errorf("invalid retry/backoff value: %v", err)
	}

	// Parse the integer as a duration.
	duration := time.Duration(seconds) * time.Second
	return &duration, nil
}
