package crawler

import (
	"testing"
)

func TestValidateInput(t *testing.T) {
	testCases := []struct {
		name          string
		u            string
		userAgent     string
		proxyURL      string
		wantErr      bool
		wantErrMsg   string
	}{
		{
			name:      "valid input",
			u:         "https://example.com",
			userAgent: "Mozilla/5.0 (compatible; Crawler/1.0)",
			wantErr:   false,
		},
		{
			name:      "invalid URL",
			u:         "invalid_url",
			userAgent: "Mozilla/5.0 (compatible; Crawler/1.0)",
			wantErr:   true,
			wantErrMsg: ErrInvalidURL.Error(),
		},
		{
			name:          "empty user agent",
			u:             "https://example.com",
			userAgent:     "",
			wantErr:       true,
			wantErrMsg:    "empty user agent",
		},
		{
			name:          "user agent with non-ASCII characters",
			u:             "https://example.com",
			userAgent:     "User-Agent problém",
			wantErr:       true,
			wantErrMsg:    "user agent contains non-ASCII characters",
		},
		{
			name:          "user agent with only dots",
			u:             "https://example.com",
			userAgent:     "........",
			wantErr:       true,
			wantErrMsg:    "user agent contains only dots",
		},
		{
			name:          "invalid proxy URL",
			u:             "https://example.com",
			userAgent:     "Mozilla/5.0 (compatible; Crawler/1.0)",
			proxyURL:      "invalid_proxy_url",
			wantErr:       true,
			wantErrMsg:    "invalid proxy URL",
		},
		{
			name:          "proxy URL without 'http://' or 'https://' prefix",
			u:             "https://example.com",
			userAgent:     "Mozilla/5.0 (compatible; Crawler/1.0)",
			proxyURL:      "proxyserver.com",
			wantErr:       true,
			wantErrMsg:    "proxy URL does not start with 'http://' or 'https://'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := ValidateInput(tc.u, tc.userAgent, tc.proxyURL); (err != nil) != tc.wantErr {
				t.Errorf("ValidateInput() error = %v, wantErr %v", err, tc.wantErr)
			}
			if err := err; err != nil && err.Error() != tc.wantErrMsg {
				t.Errorf("ValidateInput() error message = %v, want %v", err.Error(), tc.wantErrMsg)
			}
		})
	}
}
