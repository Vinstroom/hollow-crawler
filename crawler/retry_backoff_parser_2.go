package crawler

import (
	"encoding/json"
	"time"
)

// RetryBackoffParser parses the JSON response from a backoff service.
type RetryBackoffParser struct{}

// ParseJSON parses the JSON response and returns the retry-after time.
func (p *RetryBackoffParser) ParseJSON(data []byte) (time.Duration, error) {
	var response struct {
		RetryAfter time.Duration `json:"retry_after"`
	}
	err := json.Unmarshal(data, &response)
	if err != nil {
		return 0, err
	}
	return response.RetryAfter, nil
}
