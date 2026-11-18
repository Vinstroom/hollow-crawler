package crawler

import (
	"encoding/json"
	"fmt"
	"time"
)

// RetryState represents the state of a retry operation.
type RetryState struct {
	Attempts     int
	Delay         time.Duration
	MaxAttempts   int
	ShouldRetry   bool
	LastAttemptAt time.Time
}

// Serialize converts the retry state to a JSON string.
func (rs *RetryState) Serialize() (string, error) {
	bs, err := json.Marshal(rs)
	if err != nil {
		return "", fmt.Errorf("failed to marshal retry state: %w", err)
	}
	return string(bs), nil
}

// Deserialize converts a JSON string back to a retry state.
func DeserializeRetryState(jsonStr string) (*RetryState, error) {
	var rs RetryState
	err := json.Unmarshal([]byte(jsonStr), &rs)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal retry state: %w", err)
	}
	return &rs, nil
}
