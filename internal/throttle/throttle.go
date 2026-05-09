// Package throttle provides per-source write throttling for logpipe sinks.
// Unlike ratelimit (which drops records), throttle introduces backpressure
// by delaying writes when a source exceeds its configured throughput.
package throttle

import (
	"fmt"
	"sync"
	"time"
)

// Throttler enforces a maximum number of records per second for a single source
// by sleeping the caller when the budget is exhausted.
type Throttler struct {
	mu       sync.Mutex
	rate     int           // records per second
	window   time.Duration // typically 1s
	count    int
	windowAt time.Time
	sleepFn  func(time.Duration) // injectable for tests
}

// New creates a Throttler that allows at most ratePerSec records per second.
// ratePerSec must be >= 1.
func New(ratePerSec int) (*Throttler, error) {
	if ratePerSec < 1 {
		return nil, fmt.Errorf("throttle: ratePerSec must be >= 1, got %d", ratePerSec)
	}
	return &Throttler{
		rate:    ratePerSec,
		window:  time.Second,
		sleepFn: time.Sleep,
	}, nil
}

// Wait blocks until the caller is allowed to proceed under the throttle budget.
// It resets the window automatically when the current window expires.
func (t *Throttler) Wait() {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	if now.After(t.windowAt.Add(t.window)) {
		t.windowAt = now
		t.count = 0
	}

	if t.count < t.rate {
		t.count++
		return
	}

	// Budget exhausted — sleep until the current window expires.
	remaining := t.windowAt.Add(t.window).Sub(now)
	if remaining > 0 {
		t.mu.Unlock()
		t.sleepFn(remaining)
		t.mu.Lock()
	}

	t.windowAt = time.Now()
	t.count = 1
}
