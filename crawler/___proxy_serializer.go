// crawler/proxy_serializer.go
package crawler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
)

// ProxySerializer implements serialization for Proxy with config validation.
type ProxySerializer struct{}

// Serialize converts a Proxy to a JSON string with validation.
func (s ProxySerializer) Serialize(p *Proxy) (string, error) {
	if p == nil {
		return "", errors.New("cannot serialize nil proxy")
	}

	// Validate proxy address format before serialization
	if err := p.validate(); err != nil {
		return "", fmt.Errorf("proxy validation failed: %v", err)
	}

	data, err := json.Marshal(p)
	if err != nil {
		return "", fmt.Errorf("json marshal failed: %v", err)
	}
	return string(data), nil
}

// Deserialize converts a JSON string to a Proxy with validation.
func (s ProxySerializer) Deserialize(data string) (*Proxy, error) {
	if strings.TrimSpace(data) == "" {
		return nil, errors.New("empty proxy data")
	}

	var p Proxy
	if err := json.Unmarshal([]byte(data), &p); err != nil {
		return nil, fmt.Errorf("json unmarshal failed: %v", err)
	}

	// Validate proxy after deserialization
	if err := p.validate(); err != nil {
		return nil, fmt.Errorf("proxy validation failed after deserialization: %v", err)
	}

	// Ensure password is empty in the struct (even if present in JSON)
	p.Password = ""
	return &p, nil
}

// LoadProxyConfig loads proxy configuration from JSON string with schema validation.
func (s ProxySerializer) LoadProxyConfig(data string) (*Proxy, error) {
	if strings.TrimSpace(data) == "" {
		return nil, errors.New("empty proxy configuration")
	}

	var config struct {
		Proxy Proxy `json:"proxy"`
		// Additional config fields could be added here
	}

	if err := json.Unmarshal([]byte(data), &config); err != nil {
		return nil, fmt.Errorf("config unmarshal failed: %v", err)
	}

	// Validate the proxy configuration
	if err := config.Proxy.validate(); err != nil {
		return nil, fmt.Errorf("proxy configuration validation failed: %v", err)
	}

	// Create a new proxy instance with sanitized data
	proxy := config.Proxy
	proxy.Password = ""
	return &proxy, nil
}

// Proxy represents a proxy with its address and properties.
type Proxy struct {
	IP       string `json:"ip" validate:"required"` // Use string for IP to make it human-readable
	Port     int    `json:"port" validate:"required,min=1,max=65535"`
	Scheme   string `json:"scheme" validate:"required,oneof=http https socks4 socks5"`
	Type     string `json:"type" validate:"required,oneof=transparent authenticated"`
	Username string `json:"username,omitempty"`
	Password string `json:"-"` // Exclude password from JSON output
}

// validate checks the proxy's address format and required fields.
func (p *Proxy) validate() error {
	if p == nil {
		return errors.New("proxy is nil")
	}

	// Validate IP address format
	if p.IP == "" {
		return errors.New("ip is required")
	}

	if !isValidIP(p.IP) {
		return fmt.Errorf("invalid IP address format: %s", p.IP)
	}

	// Validate port range
	if p.Port < 1 || p.Port > 65535 {
		return fmt.Errorf("invalid port: %d (must be between 1 and 65535)", p.Port)
	}

	// Validate scheme
	if p.Scheme == "" {
		return errors.New("scheme is required")
	}

	validSchemes := map[string]bool{
		"http":  true,
		"https": true,
		"socks4": true,
		"socks5": true,
	}

	if !validSchemes[strings.ToLower(p.Scheme)] {
		return fmt.Errorf("invalid scheme: %s", p.Scheme)
	}

	// Validate type
	if p.Type == "" {
		return errors.New("type is required")
	}

	validTypes := map[string]bool{
		"transparent": true,
		"authenticated": true,
	}

	if !validTypes[strings.ToLower(p.Type)] {
		return fmt.Errorf("invalid type: %s", p.Type)
	}

	return nil
}

// isValidIP checks if the IP address is valid.
func isValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

func (p *Proxy) MarshalJSON() ([]byte, error) {
	type ProxyAlias Proxy
	_ = (*ProxyAlias)(p)
	return json.Marshal(ProxyAlias(p))
}

func (p *Proxy) UnmarshalJSON(data []byte) error {
	type ProxyAlias Proxy
	alias := &ProxyAlias{}
	err := json.Unmarshal(data, alias)
	if err != nil {
		return err
	}
	*p = Proxy(alias)
	return nil
}
