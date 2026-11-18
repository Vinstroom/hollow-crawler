// url_frontier_errors.go

package crawler

import "fmt"

// ErrURLFrontier represents errors specific to the URL frontier.
type ErrURLFrontier struct {
	Error string
}

func (e ErrURLFrontier) Error() string {
	return fmt.Sprintf("URL frontier error: %s", e.Error)
}

// ErrURLAlreadyInFrontier indicates that the URL is already in the frontier.
var ErrURLAlreadyInFrontier = ErrURLFrontier{"URL already in frontier"}

// ErrURLNotFound indicates that the URL was not found in the frontier.
var ErrURLNotFound = ErrURLFrontier{"URL not found"}

// ErrFrontierIsEmpty indicates that the URL frontier is empty.
var ErrFrontierIsEmpty = ErrURLFrontier{"URL frontier is empty"}
