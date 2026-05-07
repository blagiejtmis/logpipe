package checkpoint

import (
	"fmt"
	"sync"
)

// Manager holds a single shared checkpoint store and provides
// per-source offset access with safe concurrent use.
type Manager struct {
	mu    sync.RWMutex
	store *Store
}

// NewManager opens (or creates) the checkpoint file at path and returns
// a Manager ready for use. Returns an error if the store cannot be opened.
func NewManager(path string) (*Manager, error) {
	if path == "" {
		return nil, fmt.Errorf("checkpoint: manager path must not be empty")
	}
	s, err := New(path)
	if err != nil {
		return nil, fmt.Errorf("checkpoint: open store: %w", err)
	}
	return &Manager{store: s}, nil
}

// Get returns the last committed offset for source. If no offset has been
// recorded, 0 is returned along with a false ok value.
func (m *Manager) Get(source string) (offset int64, ok bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	offset, ok = m.store.Get(source)
	return
}

// Set persists offset for source. It is safe to call from multiple goroutines.
func (m *Manager) Set(source string, offset int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if err := m.store.Set(source, offset); err != nil {
		return fmt.Errorf("checkpoint: set %q: %w", source, err)
	}
	return nil
}

// Close flushes and closes the underlying store.
func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.store.Close()
}
