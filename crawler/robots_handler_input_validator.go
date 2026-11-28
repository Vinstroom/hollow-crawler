package crawler

import (
	"net/http"
	"strings"
)

// RobotsHandlerInput represents the required input for the RobotsHandler.
type RobotsHandlerInput struct {
	URL          string
	UserAgent    string
	Disallowed   []string
	Allowed      []string
	HasRobotsTXT bool
}

// Validate checks if the RobotsHandlerInput is valid.
func (i *RobotsHandlerInput) Validate() error {
	if i.URL == "" {
		return ErrMissingURL
	}
	if i.UserAgent == "" {
		return ErrMissingUserAgent
	}
	if len(i.Disallowed) == 0 && len(i.Allowed) == 0 && !i.HasRobotsTXT {
		return ErrNoRobotsInfo
	}
	if i.HasRobotsTXT && !strings.HasPrefix(i.URL, "http") {
		return ErrInvalidURLForRobotsTXT
	}
	return nil
}

// ErrMissingURL is returned when the URL is missing.
var ErrMissingURL = NewError("missing URL")

// ErrMissingUserAgent is returned when the User-Agent is missing.
var ErrMissingUserAgent = NewError("missing User-Agent")

// ErrNoRobotsInfo is returned when there's no robots info (disallowed, allowed or robots.txt).
var ErrNoRobotsInfo = NewError("no robots info")

// ErrInvalidURLForRobotsTXT is returned when the URL is not valid for retrieving robots.txt.
var ErrInvalidURLForRobotsTXT = NewError("invalid URL for retrieving robots.txt")
