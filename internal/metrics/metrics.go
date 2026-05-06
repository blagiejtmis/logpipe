// Package metrics provides lightweight in-process counters for logpipe.
// Each component increments counters that can be queried or logged on shutdown.
package metrics

import (
	"sync"
	"sync/atomic"
)

// Counter is a monotonically increasing uint64 counter.
type Counter struct {
	value uint64
}

// Inc increments the counter by 1.
func (c *Counter) Inc() { atomic.AddUint64(&c.value, 1) }

// Add increments the counter by n.
func (c *Counter) Add(n uint64) { atomic.AddUint64(&c.value, n) }

// Value returns the current counter value.
func (c *Counter) Value() uint64 { return atomic.LoadUint64(&c.value) }

// Registry holds a named set of counters.
type Registry struct {
	mu       sync.RWMutex
	counters map[string]*Counter
}

// NewRegistry creates an empty Registry.
func NewRegistry() *Registry {
	return &Registry{counters: make(map[string]*Counter)}
}

// Counter returns the counter for name, creating it if it does not exist.
func (r *Registry) Counter(name string) *Counter {
	r.mu.RLock()
	c, ok := r.counters[name]
	r.mu.RUnlock()
	if ok {
		return c
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	// double-checked
	if c, ok = r.counters[name]; ok {
		return c
	}
	c = &Counter{}
	r.counters[name] = c
	return c
}

// Snapshot returns a copy of all counter values keyed by name.
func (r *Registry) Snapshot() map[string]uint64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string]uint64, len(r.counters))
	for name, c := range r.counters {
		out[name] = c.Value()
	}
	return out
}
