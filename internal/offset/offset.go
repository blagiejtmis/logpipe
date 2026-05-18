// Package offset tracks per-source byte offsets so that log tailing can
// resume from the last known position after a restart.
package offset

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

// Store persists and retrieves byte offsets keyed by source name.
type Store struct {
	mu   sync.RWMutex
	path string
	data map[string]int64
}

// New opens (or creates) the offset store at path.
func New(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("offset: path must not be empty")
	}
	s := &Store{path: path, data: make(map[string]int64)}
	if err := s.load(); err != nil {
		return nil, err
	}
	return s, nil
}

// Get returns the stored offset for source, or 0 if not found.
func (s *Store) Get(source string) int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data[source]
}

// Set updates the offset for source and flushes to disk.
func (s *Store) Set(source string, offset int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[source] = offset
	return s.flush()
}

// Delete removes the stored offset for source and flushes to disk.
func (s *Store) Delete(source string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, source)
	return s.flush()
}

func (s *Store) load() error {
	b, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &s.data)
}

func (s *Store) flush() error {
	b, err := json.Marshal(s.data)
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, b, 0o644)
}
