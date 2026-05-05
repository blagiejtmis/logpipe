package tail

import (
	"context"
	"fmt"
	"sync"
)

// Manager supervises multiple Tailers and fans their output into one channel.
type Manager struct {
	output chan Line
	wg     sync.WaitGroup
	mu     sync.Mutex
	errs   []error
}

// NewManager creates a Manager with a buffered output channel.
func NewManager(bufSize int) *Manager {
	return &Manager{
		output: make(chan Line, bufSize),
	}
}

// Add registers a file path to be tailed.
func (m *Manager) Add(ctx context.Context, path string) {
	t := New(path, m.output)
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		if err := t.Run(ctx); err != nil && ctx.Err() == nil {
			m.mu.Lock()
			m.errs = append(m.errs, fmt.Errorf("tailer %s: %w", path, err))
			m.mu.Unlock()
		}
	}()
}

// Output returns the read-only channel of tailed lines.
func (m *Manager) Output() <-chan Line {
	return m.output
}

// Wait blocks until all tailers have stopped and then closes the output channel.
func (m *Manager) Wait() []error {
	m.wg.Wait()
	close(m.output)
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.errs
}
