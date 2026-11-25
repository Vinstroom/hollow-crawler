package crawler

import (
	"html"
	"testing"
	"time"
)

func TestTimeStamp(t *testing.T) {
	testCases := []struct {
		name     string
		htmlStr  string
		expected string
	}{
		{
			name:    "no_time_tags",
			htmlStr: "<p>No time tags here</p>",
			expected: "",
		},
		{
			name:    "invalid_dates",
			htmlStr: `<time datetime="invalid">Invalid</time>`,
			expected: "",
		},
		{
			name:    "future_dates",
			htmlStr: `<time datetime="2100-01-01T00:00:00Z">Future</time>`,
			expected: "",
		},
		{
			name:    "multiple_dates",
			htmlStr: `<time datetime="2022-01-01T00:00:00Z">Old</time><time datetime="2020-01-01T00:00:00Z">Newer</time>`,
			expected: "Published on Older on January 1, 2020",
		},
		{
			name:    "multiple_invalid_dates",
			htmlStr: `<time datetime="2020-01-01T00:00:00Z">Old</time><time>Invalid</time>`,
			expected: "Published on Older on January 1, 2020",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := TimeStamp(tc.htmlStr)
			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
	}
}
