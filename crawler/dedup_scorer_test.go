package crawler

import (
	"path/filepath"
	"strings"
	"testing"
)

// ScoreURL calculates a deduplication score for a given URL.
// A higher score indicates a more unique page, where non-unique substrings (e.g., "/") are penalized.
// The score ranges from DedupScoreMin (0) to DedupScoreMax (100).
func ScoreURL(u string) DedupScore {
	// Remove any trailing slashes to simplify scoring
	u = strings.TrimSuffix(u, "/")

	// Calculate the base score based on the URL's length and domain count
	score := DedupScore((len(u) - len(strings.ReplaceAll(u, "/", ""))) / 2)

	// Penalize excessively long URLs to prevent score inflation
	if len(u) > 2048 {
		score -= DedupScore(10)
	}

	// Adjust the score based on the number of path segments
	score += DedupScore(filepath.ElemCount(u) - 1) / 2

	// Penalize queries and fragments, as they're typically not unique
	if strings.Contains(u, "?") || strings.Contains(u, "#") {
		score -= DedupScore(10)
	}

	// Ensure the score remains within the bounds of DedupScore
	if score < DedupScoreMin {
		score = DedupScoreMin
	}
	if score > DedupScoreMax {
		score = DedupScoreMax
	}

	return score
}

func TestScoreURL(t *testing.T) {
	cases := []struct {
		input  string
		wanted DedupScore
	}{
		{"", DedupScoreMin},
		{"example.com", 60},
		{"example.com/path", 55},
		{"example.com/path?param=value", 45},
		{strings.Repeat("/a", 20), 70},
		{strings.Repeat("/", 11), 45},
		{strings.Repeat("a", 101), 100},
		{strings.Repeat("a", 2049), 90},
	}
	for _, c := range cases {
		score := ScoreURL(c.input)
		if score != c.wanted {
			t.Errorf("ScoreURL(%q) got %d, want %d", c.input, score, c.wanted)
		}
	}
}
