//go:generate go run generate.go

package crawler

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// UserAgentPool represents a pool of user agents.
type UserAgentPool []string

// URLSerializer serializes and deserializes URL frontier and UserAgentPool.
type URLSerializer struct {
	mu sync.Mutex
}

// SerializeUserAgentPool serializes the given UserAgentPool.
func (s *URLSerializer) SerializeUserAgentPool(users []string, destPath string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.Marshal(users)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(destPath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// DeserializeUserAgentPool deserializes a UserAgentPool from the given path.
func (s *URLSerializer) DeserializeUserAgentPool(srcPath string) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := ioutil.ReadFile(srcPath)
	if err != nil {
		return nil, err
	}

	var users []string
	err = json.Unmarshal(data, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// SaveUserAgentPool saves the given UserAgentPool to the specified file path.
func (s *URLSerializer) SaveUserAgentPool(users []string, filePath string) error {
	_, filename := filepath.Split(filePath)
	if filename == "" {
		return ErrInvalidFilePath
	}

	return s.SerializeUserAgentPool(users, filePath)
}

// LoadUserAgentPool loads a UserAgentPool from the specified file path.
func (s *URLSerializer) LoadUserAgentPool(filePath string) ([]string, error) {
	_, filename := filepath.Split(filePath)
	if filename == "" {
		return nil, ErrInvalidFilePath
	}

	return s.DeserializeUserAgentPool(filePath)
}

// ... (rest of the file remains unchanged)
