package crawler

import (
	"regexp"
	"strings"
	"time"
)

// retryBackoffNormalizer normalizes retry/backoff strings to a consistent format.
type retryBackoffNormalizer struct{}

// normalize converts the given string to the format "seconds:attempts".
func (n *retryBackoffNormalizer) normalize(s string) (string, error) {
	re := regexp.MustCompile(`^(\d+)s:(\d+)$`)
	match := re.FindStringSubmatch(s)
	if match == nil {
		return "", fmt.Errorf("invalid retry/backoff format: %s", s)
	}
	return s, nil
}

// parse converts the normalized string to a time.Duration and int.
func (n *retryBackoffNormalizer) parse(s string) (time.Duration, int, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid retry/backoff format: %s", s)
	}
	seconds, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}
	attempts, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}
	return time.Duration(seconds) * time.Second, attempts, nil
}
