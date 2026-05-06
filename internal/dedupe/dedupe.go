// Package dedupe provides log record deduplication based on a configurable
// field fingerprint and a sliding time window.
package dedupe

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"time"
)

// Record is the minimal interface dedupe needs from a log record.
type Record map[string]string

// Deduplicator drops records whose fingerprint has been seen within Window.
type Deduplicator struct {
	fields []string
	window time.Duration

	mu    sync.Mutex
	seen  map[string]time.Time
	nowFn func() time.Time
}

// New creates a Deduplicator that keys on the given fields and suppresses
// duplicates for the duration of window.
func New(fields []string, window time.Duration) (*Deduplicator, error) {
	if len(fields) == 0 {
		return nil, fmt.Errorf("dedupe: at least one field required")
	}
	if window <= 0 {
		return nil, fmt.Errorf("dedupe: window must be positive")
	}
	return &Deduplicator{
		fields: fields,
		window: window,
		seen:   make(map[string]time.Time),
		nowFn:  time.Now,
	}, nil
}

// Allow returns true when the record should be forwarded (not a duplicate).
func (d *Deduplicator) Allow(rec Record) bool {
	fp := d.fingerprint(rec)
	now := d.nowFn()

	d.mu.Lock()
	defer d.mu.Unlock()

	if t, ok := d.seen[fp]; ok && now.Sub(t) < d.window {
		return false
	}
	d.seen[fp] = now
	return true
}

// Purge removes expired entries; call periodically to bound memory usage.
func (d *Deduplicator) Purge() {
	now := d.nowFn()
	d.mu.Lock()
	defer d.mu.Unlock()
	for fp, t := range d.seen {
		if now.Sub(t) >= d.window {
			delete(d.seen, fp)
		}
	}
}

func (d *Deduplicator) fingerprint(rec Record) string {
	h := sha256.New()
	for _, f := range d.fields {
		fmt.Fprintf(h, "%s=%s;", f, rec[f])
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
