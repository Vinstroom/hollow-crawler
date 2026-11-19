package crawler

import (
	"regexp"
	"strings"
)

// userAgentNormalizer normalizes user-agent strings and enforces reasonable length constraints.
type userAgentNormalizer struct{}

// Normalize takes a user-agent string and returns a normalized version.
// It enforces a minimum length of 3 characters after normalization and
// trims to a maximum of 100 characters to prevent overly verbose strings.
func (u userAgentNormalizer) Normalize(ua string) string {
	if ua == "" {
		return ""
	}

	// Convert to lowercase and normalize whitespace
	ua = strings.ToLower(ua)
	ua = regexp.MustCompile(`\s+`).ReplaceAllString(ua, " ")

	// Remove problematic characters and tokens
	ua = regexp.MustCompile(`/[^\d\s.]+|;|,`).ReplaceAllString(ua, "")
	ua = strings.TrimSpace(ua)

	// Enforce minimum length (avoid empty strings from over-normalization)
	if len(ua) < 3 {
		return ""
	}

	// Enforce reasonable maximum length
	if len(ua) > 100 {
		ua = ua[:100]
	}

	return ua
}
