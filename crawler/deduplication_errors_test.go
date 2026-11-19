package crawler

import (
	"errors"
	"testing"
)

func TestErrorStrings(t *testing.T) {
	tests := []struct {
		err     error
		want    string
		fp      string
	}{
		{
			err:  ErrDeduplicationFailed{Err: errors.New("test")},
			want: "failed to deduplicate URL: - test",
			fp:   "test",
		},
		{
			err:     ErrExportItemFailed{Item: "item", Err: errors.New("test"), Fingerprint: "fp"},
			want:    "failed to export item: item - test",
			fp:      "fp",
		},
		{
			err: ErrAlreadyProcessed{URL: "https://example.com"},
			want: "URL already processed: https://example.com",
		},
		{
			err:             ErrDuplicateFingerprint{Fingerprint: "fp"},
			want:            "duplicate fingerprint encountered: fp",
		},
		{
			err:     ErrExportFailed{Err: errors.New("test"), Fingerprint: "fp"},
			want:    "failed to export URLs with fingerprint: fp - test",
			fp:      "fp",
		},
	}

	for _, tt := range tests {
		err := tt.err
		if got := err.Error(); got != tt.want {
			t.Errorf("Error() = %v, want %v", got, tt.want)
		}
	}
}
