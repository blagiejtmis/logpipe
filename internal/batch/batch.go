// Package batch provides size- and time-based batching of log records.
// Records are accumulated until either a maximum count is reached or a
// flush interval elapses, whichever comes first.
package batch

import (
	"context"
	"sync"
	"time"
)

// Record is a map of string fields representing a single log entry.
type Record = map[string]string

// FlushFunc is called with a slice of records when the batch is flushed.
type FlushFunc func([]Record)

// Batcher accumulates records and flushes them in batches.
type Batcher struct {
	maxSize  int
	interval time.Duration
	flush    FlushFunc

	mu  sync.Mutex
	buf []Record
}

// New creates a Batcher that flushes when buf reaches maxSize records or
// after interval, whichever comes first. maxSize must be >= 1 and interval
// must be > 0.
func New(maxSize int, interval time.Duration, flush FlushFunc) (*Batcher, error) {
	if maxSize < 1 {
		return nil, ErrInvalidMaxSize
	}
	if interval <= 0 {
		return nil, ErrInvalidInterval
	}
	if flush == nil {
		return nil, ErrNilFlushFunc
	}
	return &Batcher{
		maxSize:  maxSize,
		interval: interval,
		flush:    flush,
		buf:      make([]Record, 0, maxSize),
	}, nil
}

// Add appends a record to the current batch. If the batch reaches maxSize
// it is flushed synchronously before Add returns.
func (b *Batcher) Add(r Record) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.buf = append(b.buf, r)
	if len(b.buf) >= b.maxSize {
		b.flushLocked()
	}
}

// Run starts the interval-based flush loop. It blocks until ctx is cancelled,
// at which point any buffered records are flushed and Run returns.
func (b *Batcher) Run(ctx context.Context) {
	ticker := time.NewTicker(b.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			b.mu.Lock()
			b.flushLocked()
			b.mu.Unlock()
		case <-ctx.Done():
			b.mu.Lock()
			b.flushLocked()
			b.mu.Unlock()
			return
		}
	}
}

// flushLocked sends the current buffer to the FlushFunc and resets the buffer.
// Caller must hold b.mu.
func (b *Batcher) flushLocked() {
	if len(b.buf) == 0 {
		return
	}
	out := make([]Record, len(b.buf))
	copy(out, b.buf)
	b.buf = b.buf[:0]
	b.flush(out)
}
