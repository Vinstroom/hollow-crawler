package crawler

import (
	"time"
)

// TimeUtils provides utility functions for time and date manipulation.
type TimeUtils struct{}

// FormatTime formats a time.Time object in RFC3339 format.
func (tu *TimeUtils) FormatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

// ParseTime parses a string representing time in RFC3339 format.
func (tu *TimeUtils) ParseTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

// DurationInDays returns the number of days in the given duration.
func (tu *TimeUtils) DurationInDays(d time.Duration) float64 {
	return float64(d.Hours()) / 24
}

// TruncateToDays truncates the given time to the nearest day.
func (tu *TimeUtils) TruncateToDays(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// TimeUtilsInstance returns a singleton instance of TimeUtils.
var TimeUtilsInstance TimeUtils

// DurationInDays returns the number of days in the given duration.
func DurationInDays(d time.Duration) float64 {
	return TimeUtilsInstance.DurationInDays(d)
}

// TruncateToDays truncates the given time to the nearest day.
func TruncateToDays(t time.Time) time.Time {
	return TimeUtilsInstance.TruncateToDays(t)
}
