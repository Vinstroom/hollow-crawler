package crawler

import (
	"encoding/json"
	"io"
	"os"
	"sync"

	"github.com/pkg/errors"
)

// Serializer is an interface for objects that can serialize and deserialize themselves.
type Serializer interface {
	Serialize() ([]byte, error)
	Deserialize(data []byte) error
}

// JSONSerializer wraps an object and provides JSON serialization and deserialization.
type JSONSerializer[T any] struct {
	mu     sync.Mutex
	obj    *T
	loaded bool
}

// NewJSONSerializer creates a new JSONSerializer with a given object.
func NewJSONSerializer[T any](obj T) *JSONSerializer[T] {
	return &JSONSerializer[T]{obj: &obj}
}

// Serialize serializes the wrapped object to JSON.
func (s *JSONSerializer[T]) Serialize() ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return json.Marshal(s.obj)
}

// Deserialize deserializes the given JSON data into the wrapped object.
func (s *JSONSerializer[T]) Deserialize(data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	err := json.Unmarshal(data, s.obj)
	if err != nil {
		return errors.Wrap(err, "failed to deserialize object")
	}
	s.loaded = true
	return nil
}

// SaveToFile saves the serialized object to the given file path.
func (s *JSONSerializer[T]) SaveToFile(filePath string) error {
	data, err := s.Serialize()
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

// LoadFromFile loads the object from the given file path.
func (s *JSONSerializer[T]) LoadFromFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	return s.Deserialize(data)
}

// SaveToWriter saves the serialized object to the given io.Writer.
func (s *JSONSerializer[T]) SaveToWriter(w io.Writer) error {
	data, err := s.Serialize()
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}
