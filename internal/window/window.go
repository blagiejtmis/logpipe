// Package window provides a sliding-window counter used to track
// event frequency over a rolling time interval.
package window

import (
	"errors"
	"sync"
	"time"
)

// Counter is a thread-safe sliding-window event counter.
type Counter struct {
	mu       sync.Mutex
	buckets  []int64
	size     int
	bucket   time.Duration
	last     time.Time
	total    int64
}

// New creates a Counter that divides [window] into [buckets] equal slots.
// window must be positive and buckets must be >= 1.
func New(window time.Duration, buckets int) (*Counter, error) {
	if window <= 0 {
		return nil, errors.New("window: duration must be positive")
	}
	if buckets < 1 {
		return nil, errors.New("window: buckets must be >= 1")
	}
	return &Counter{
		buckets: make([]int64, buckets),
		size:    buckets,
		bucket:  window / time.Duration(buckets),
		last:    time.Now(),
	}, nil
}

// Add records n events at the current time.
func (c *Counter) Add(n int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.advance(time.Now())
	idx := c.currentIdx(time.Now())
	c.buckets[idx] += n
	c.total += n
}

// Count returns the total number of events recorded across the entire window.
func (c *Counter) Count() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.advance(time.Now())
	var sum int64
	for _, v := range c.buckets {
		sum += v
	}
	return sum
}

// advance expires buckets that are older than the window.
func (c *Counter) advance(now time.Time) {
	elapsed := now.Sub(c.last)
	if elapsed < c.bucket {
		return
	}
	steps := int(elapsed / c.bucket)
	if steps >= c.size {
		for i := range c.buckets {
			c.buckets[i] = 0
		}
		c.total = 0
		c.last = now
		return
	}
	for i := 0; i < steps; i++ {
		oldIdx := c.currentIdx(c.last.Add(time.Duration(i) * c.bucket))
		c.total -= c.buckets[oldIdx]
		c.buckets[oldIdx] = 0
	}
	c.last = c.last.Add(time.Duration(steps) * c.bucket)
}

func (c *Counter) currentIdx(t time.Time) int {
	return int(t.UnixNano()/int64(c.bucket)) % c.size
}
