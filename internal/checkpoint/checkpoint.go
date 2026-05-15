// Package checkpoint persists tail offsets so logpipe can resume from
// where it left off after a restart.
package checkpoint

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

// Store persists and retrieves file offsets keyed by source path.
type Store struct {
	mu      sync.Mutex
	path    string
	offsets map[string]int64
}

// New opens (or creates) the checkpoint file at path and loads any
// previously saved offsets into memory.
func New(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("checkpoint: path must not be empty")
	}
	s := &Store{
		path:    path,
		offsets: make(map[string]int64),
	}
	data, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	if len(data) > 0 {
		if err := json.Unmarshal(data, &s.offsets); err != nil {
			return nil, err
		}
	}
	return s, nil
}

// Get returns the last saved offset for source, or 0 if none exists.
func (s *Store) Get(source string) int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.offsets[source]
}

// Set updates the in-memory offset for source and flushes to disk.
func (s *Store) Set(source string, offset int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.offsets[source] = offset
	return s.flush()
}

// Delete removes the saved offset for source and flushes to disk.
// It is a no-op if source has no recorded offset.
func (s *Store) Delete(source string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.offsets[source]; !ok {
		return nil
	}
	delete(s.offsets, source)
	return s.flush()
}

// flush writes the current offsets map to disk atomically.
func (s *Store) flush() error {
	data, err := json.Marshal(s.offsets)
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}
