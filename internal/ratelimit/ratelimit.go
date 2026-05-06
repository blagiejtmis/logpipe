// Package ratelimit provides per-source log line rate limiting.
package ratelimit

import (
	"fmt"
	"sync"
	"time"
)

// Limiter enforces a maximum number of log lines per second for a given source.
type Limiter struct {
	mu       sync.Mutex
	buckets  map[string]*bucket
	maxLines int
	window   time.Duration
}

type bucket struct {
	count     int
	windowEnd time.Time
}

// New creates a Limiter that allows at most maxLines lines per window duration
// for each source. Returns an error if maxLines <= 0 or window <= 0.
func New(maxLines int, window time.Duration) (*Limiter, error) {
	if maxLines <= 0 {
		return nil, fmt.Errorf("ratelimit: maxLines must be > 0, got %d", maxLines)
	}
	if window <= 0 {
		return nil, fmt.Errorf("ratelimit: window must be > 0, got %s", window)
	}
	return &Limiter{
		buckets:  make(map[string]*bucket),
		maxLines: maxLines,
		window:   window,
	}, nil
}

// Allow reports whether a log line from source should be allowed through.
// It returns false once the per-window limit is exceeded for that source.
func (l *Limiter) Allow(source string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	b, ok := l.buckets[source]
	if !ok || now.After(b.windowEnd) {
		l.buckets[source] = &bucket{
			count:     1,
			windowEnd: now.Add(l.window),
		}
		return true
	}

	if b.count >= l.maxLines {
		return false
	}
	b.count++
	return true
}

// Reset clears rate-limit state for all sources.
func (l *Limiter) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.buckets = make(map[string]*bucket)
}
