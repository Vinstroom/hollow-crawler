package crawler

import (
	"strings"
	"unicode"
)

// normalizeUserAgent takes a UserAgent string and returns a normalized version of it.
// This function is case insensitive, removes whitespace, and ensures consistent string format.
func normalizeUserAgent(userAgent string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return unicode.ToLower(r)
	}, userAgent)
}
