package crawler

import (
	"fmt"
	"net/url"
	"strings"
)

type URLFrontierValidator struct {
	MaxURLLength int
}

func (v *URLFrontierValidator) Validate(u string) error {
	if strings.TrimSpace(u) == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	u, err := url.Parse(strings.TrimSpace(u))
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	if len(u.String()) > v.MaxURLLength {
		return fmt.Errorf("URL is too long (max: %d)", v.MaxURLLength)
	}

	if u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("URL must contain a scheme and a host")
	}

	return nil
}
