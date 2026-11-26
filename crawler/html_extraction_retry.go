package crawler

import (
	"errors"
	"time"
)

type htmlExtractionRetryPolicy interface {
	shouldRetry(err error) bool
	delayOnRetry(err error) time.Duration
}

type defaultHtmlExtractionRetryPolicy struct{}

func (d *defaultHtmlExtractionRetryPolicy) shouldRetry(err error) bool {
	return errors.Is(err, ErrHtmlExtractionFailed)
}

func (d *defaultHtmlExtractionRetryPolicy) delayOnRetry(err error) time.Duration {
	return time.Second * 5
}

func RetryHtmlExtraction(f func() error, policy htmlExtractionRetryPolicy, maxAttempts int) error {
	attempts := 0
	for attempts < maxAttempts {
		err := f()
		if err == nil || !policy.shouldRetry(err) {
			return err
		}
		time.Sleep(policy.delayOnRetry(err))
		attempts++
	}
	return errors.New("failed to extract html after multiple attempts")
}

func WithDefaultHtmlExtractionRetryPolicy(f func() error) func() error {
	return func() error {
		return RetryHtmlExtraction(f, &defaultHtmlExtractionRetryPolicy{}, 3)
	}
}
