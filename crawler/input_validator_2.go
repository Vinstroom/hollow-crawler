package crawler

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// URLValidator validates the input URL.
type URLValidator struct {
	validURLRegex *regexp.Regexp
}

// NewURLValidator creates a new URLValidator.
func NewURLValidator() *URLValidator {
	return &URLValidator{
		validURLRegex: regexp.MustCompile(`^(?:http|https):\/\/\S+$`),
	}
}

// Validate checks if the given URL is valid.
func (v *URLValidator) Validate(u string) error {
	if !v.validURLRegex.MatchString(u) {
		return fmt.Errorf("invalid URL: %s", u)
	}
	_, err := url.ParseRequestURI(u)
	if err != nil {
		return fmt.Errorf("invalid URL: %s", u)
	}
	return nil
}

// ProxyValidator validates the input proxy.
type ProxyValidator struct {
	validProxyRegex *regexp.Regexp
}

// NewProxyValidator creates a new ProxyValidator.
func NewProxyValidator() *ProxyValidator {
	return &ProxyValidator{
		validProxyRegex: regexp.MustCompile(`^(?:http|https|socks)://\S+$`),
	}
}

// Validate checks if the given proxy is valid.
func (v *ProxyValidator) Validate(p string) error {
	if !v.validProxyRegex.MatchString(p) {
		return fmt.Errorf("invalid proxy: %s", p)
	}
	_, err := url.Parse(p)
	if err != nil {
		return fmt.Errorf("invalid proxy: %s", p)
	}
	if strings.Split(p, "://")[0] == "socks" && strings.Split(p, ":")[1] != "4444" {
		return fmt.Errorf("invalid SOCKS port: %s", p)
	}
	return nil
}
