// Package buffer provides an in-memory ring buffer for log records,
// allowing downstream sinks to absorb bursts without dropping lines.
package buffer

import (
	"errors"
	"sync"
)

// ErrFull is returned by Push when the buffer has reached its capacity
// and the overflow policy is to drop new records.
var ErrFull = errors.New("buffer: capacity exceeded")

// Record is a single log entry held in the ring buffer.
type Record struct {
	Source string
	Fields map[string]string
}

// Buffer is a bounded, thread-safe FIFO queue for log records.
type Buffer struct {
	mu       sync.Mutex
	items    []Record
	cap      int
	dropNew  bool
	dropped  int64
}

// New creates a Buffer with the given capacity.
// When dropNew is true, Push returns ErrFull when full instead of evicting
// the oldest record.
func New(capacity int, dropNew bool) (*Buffer, error) {
	if capacity <= 0 {
		return nil, errors.New("buffer: capacity must be > 0")
	}
	return &Buffer{
		items:   make([]Record, 0, capacity),
		cap:     capacity,
		dropNew: dropNew,
	}, nil
}

// Push adds a record to the buffer. If the buffer is full and dropNew is
// true, ErrFull is returned. If dropNew is false the oldest record is
// silently evicted to make room.
func (b *Buffer) Push(r Record) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.items) >= b.cap {
		if b.dropNew {
			b.dropped++
			return ErrFull
		}
		// evict oldest
		b.items = b.items[1:]
		b.dropped++
	}
	b.items = append(b.items, r)
	return nil
}

// Pop removes and returns the oldest record. The second return value is
// false when the buffer is empty.
func (b *Buffer) Pop() (Record, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.items) == 0 {
		return Record{}, false
	}
	r := b.items[0]
	b.items = b.items[1:]
	return r, true
}

// Len returns the current number of records in the buffer.
func (b *Buffer) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.items)
}

// Dropped returns the cumulative count of records that were dropped due to
// overflow.
func (b *Buffer) Dropped() int64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.dropped
}
