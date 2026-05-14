// Package ceiling provides a per-source record count ceiling that drops
// records once a configurable maximum has been reached within a time window.
package ceiling

import (
	"errors"
	"sync"
	"time"
)

// Ceiling enforces a maximum number of records accepted within a rolling window.
type Ceiling struct {
	mu      sync.Mutex
	max     int64
	window  time.Duration
	count   int64
	windowStart time.Time
	now     func() time.Time
}

// New creates a Ceiling that allows at most max records per window duration.
func New(max int64, window time.Duration) (*Ceiling, error) {
	if max <= 0 {
		return nil, errors.New("ceiling: max must be greater than zero")
	}
	if window <= 0 {
		return nil, errors.New("ceiling: window must be greater than zero")
	}
	return &Ceiling{
		max:         max,
		window:      window,
		now:         time.Now,
		windowStart: time.Now(),
	}, nil
}

// Allow returns true if the record should be accepted, false if the ceiling has
// been reached for the current window.
func (c *Ceiling) Allow() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := c.now()
	if now.Sub(c.windowStart) >= c.window {
		c.count = 0
		c.windowStart = now
	}

	if c.count >= c.max {
		return false
	}
	c.count++
	return true
}

// Reset resets the counter and window start to the current time.
func (c *Ceiling) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count = 0
	c.windowStart = c.now()
}
