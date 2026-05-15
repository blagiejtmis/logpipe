// Package uniq provides field-level uniqueness enforcement across a sliding
// time window. Unlike dedupe (which hashes full records), uniq tracks a single
// key field and drops records whose key value has already been seen within the
// configured window duration.
package uniq

import (
	"errors"
	"sync"
	"time"
)

// Record is the minimal interface uniq needs from a log record.
type Record map[string]any

// Uniq enforces uniqueness on a single key field within a time window.
type Uniq struct {
	field  string
	window time.Duration
	mu     sync.Mutex
	seen   map[string]time.Time
	now    func() time.Time
}

// New creates a Uniq that tracks values of field within window.
// field must be non-empty and window must be positive.
func New(field string, window time.Duration) (*Uniq, error) {
	if field == "" {
		return nil, errors.New("uniq: field must not be empty")
	}
	if window <= 0 {
		return nil, errors.New("uniq: window must be positive")
	}
	return &Uniq{
		field:  field,
		window: window,
		seen:   make(map[string]time.Time),
		now:    time.Now,
	}, nil
}

// Allow returns true if the record's key field has NOT been seen within the
// current window. Records with a missing or non-string key field are always
// allowed through.
func (u *Uniq) Allow(r Record) bool {
	v, ok := r[u.field]
	if !ok {
		return true
	}
	key, ok := v.(string)
	if !ok {
		return true
	}

	now := u.now()
	u.mu.Lock()
	defer u.mu.Unlock()

	// Expire stale entries lazily.
	for k, ts := range u.seen {
		if now.Sub(ts) >= u.window {
			delete(u.seen, k)
		}
	}

	if _, exists := u.seen[key]; exists {
		return false
	}
	u.seen[key] = now
	return true
}
