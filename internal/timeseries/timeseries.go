// Package timeseries provides a rolling time-series counter that tracks
// per-field value frequencies over a sliding window of fixed-size buckets.
package timeseries

import (
	"errors"
	"sync"
	"time"
)

// Series tracks value counts for a single field over a sliding window.
type Series struct {
	mu      sync.Mutex
	field   string
	buckets []map[string]int64
	size    int
	dur     time.Duration
	start   time.Time
}

// New creates a Series that partitions [window] into [buckets] equal slots.
// field is the log-record key whose values are counted.
func New(field string, window time.Duration, buckets int) (*Series, error) {
	if field == "" {
		return nil, errors.New("timeseries: field must not be empty")
	}
	if window <= 0 {
		return nil, errors.New("timeseries: window must be positive")
	}
	if buckets < 1 {
		return nil, errors.New("timeseries: buckets must be >= 1")
	}
	s := &Series{
		field:   field,
		buckets: make([]map[string]int64, buckets),
		size:    buckets,
		dur:     window / time.Duration(buckets),
		start:   time.Now(),
	}
	for i := range s.buckets {
		s.buckets[i] = make(map[string]int64)
	}
	return s, nil
}

// Record increments the count for the value of s.field found in rec.
func (s *Series) Record(rec map[string]any) {
	val, _ := rec[s.field].(string)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rotate()
	s.buckets[s.current()][val]++
}

// Counts returns a snapshot of total counts per value across all live buckets.
func (s *Series) Counts() map[string]int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rotate()
	out := make(map[string]int64)
	for _, b := range s.buckets {
		for k, v := range b {
			out[k] += v
		}
	}
	return out
}

func (s *Series) current() int {
	elapsed := int(time.Since(s.start) / s.dur)
	return elapsed % s.size
}

// rotate clears buckets that have aged out.
func (s *Series) rotate() {
	now := time.Now()
	slots := int(now.Sub(s.start) / s.dur)
	if slots == 0 {
		return
	}
	clear := slots
	if clear > s.size {
		clear = s.size
	}
	prev := int((slots - 1) % s.size)
	for i := 1; i <= clear; i++ {
		idx := (prev + i) % s.size
		s.buckets[idx] = make(map[string]int64)
	}
	s.start = s.start.Add(time.Duration(slots) * s.dur)
}
