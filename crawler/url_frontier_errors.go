package crawler

import (
	"fmt"
	"strings"
)

type (
	URLFrontierError struct {
		cause   error
		frontierID string
	}
)

func NewFrontierError(err error, frontierID string) *URLFrontierError {
	// Wrap the error with additional context
	err = fmt.Errorf("%w: URL frontier (%s) encountered an error", err, frontierID)
	return &URLFrontierError{
		cause:      err,
		frontierID: frontierID,
	}
}

func (e *URLFrontierError) Error() string {
	// Return a formatted error message with the frontier ID
	return fmt.Sprintf("%v", e.cause)
}

func IsURLFrontierError(err error) bool {
	// Check if the error is of type *URLFrontierError
	_, ok := err.(*URLFrontierError)
	return ok
}

func (e *URLFrontierError) Unwrap() error {
	// Return the original cause of the error
	return e.cause
}

func (e *URLFrontierError) Cause() error {
	// Alias for Unwrap to maintain compatibility with the original function
	return e.Unwrap()
}
