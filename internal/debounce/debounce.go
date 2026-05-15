// Package debounce suppresses repeated log records for the same key within
// a configurable quiet window. Only the first occurrence is forwarded; all
// subsequent duplicates are dropped until the window expires.
package debounce

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Debouncer holds state for a single debounce rule.
type Debouncer struct {
	field   string
	window  time.Duration
	mu      sync.Mutex
	seen    map[string]time.Time
	nowFunc func() time.Time
}

// New creates a Debouncer that suppresses records where the value of field
// repeats within window. window must be positive and field must be non-empty.
func New(field string, window time.Duration) (*Debouncer, error) {
	if field == "" {
		return nil, errors.New("debounce: field must not be empty")
	}
	if window <= 0 {
		return nil, errors.New("debounce: window must be positive")
	}
	return &Debouncer{
		field:   field,
		window:  window,
		seen:    make(map[string]time.Time),
		nowFunc: time.Now,
	}, nil
}

// Allow returns true if the record should be forwarded (first occurrence or
// window has expired) and false if it should be suppressed.
func (d *Debouncer) Allow(record map[string]any) bool {
	v, _ := record[d.field]
	key := fmt.Sprintf("%v", v)

	now := d.nowFunc()
	d.mu.Lock()
	defer d.mu.Unlock()

	if last, ok := d.seen[key]; ok && now.Sub(last) < d.window {
		return false
	}
	d.seen[key] = now
	return true
}

// Purge removes entries whose window has elapsed, preventing unbounded growth.
func (d *Debouncer) Purge() {
	now := d.nowFunc()
	d.mu.Lock()
	defer d.mu.Unlock()
	for k, t := range d.seen {
		if now.Sub(t) >= d.window {
			delete(d.seen, k)
		}
	}
}
