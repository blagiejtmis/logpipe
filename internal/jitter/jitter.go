// Package jitter applies randomised delay to log records, useful for
// smoothing burst traffic before it reaches downstream sinks.
package jitter

import (
	"fmt"
	"math/rand"
	"time"
)

// Jitter holds configuration for a single jitter rule.
type Jitter struct {
	minDelay time.Duration
	maxDelay time.Duration
	rng      *rand.Rand
}

// New creates a Jitter that sleeps a random duration in [minDelay, maxDelay]
// before returning. Both values must be non-negative and min must be <= max.
func New(minDelay, maxDelay time.Duration) (*Jitter, error) {
	if minDelay < 0 {
		return nil, fmt.Errorf("jitter: minDelay must be non-negative, got %s", minDelay)
	}
	if maxDelay < 0 {
		return nil, fmt.Errorf("jitter: maxDelay must be non-negative, got %s", maxDelay)
	}
	if minDelay > maxDelay {
		return nil, fmt.Errorf("jitter: minDelay %s must be <= maxDelay %s", minDelay, maxDelay)
	}
	return &Jitter{
		minDelay: minDelay,
		maxDelay: maxDelay,
		rng:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}, nil
}

// Wait blocks for a random duration within the configured range.
// If min == max the delay is deterministic.
func (j *Jitter) Wait() {
	span := j.maxDelay - j.minDelay
	var d time.Duration
	if span == 0 {
		d = j.minDelay
	} else {
		d = j.minDelay + time.Duration(j.rng.Int63n(int64(span)+1))
	}
	if d > 0 {
		time.Sleep(d)
	}
}

// Range returns the configured [min, max] delay pair.
func (j *Jitter) Range() (time.Duration, time.Duration) {
	return j.minDelay, j.maxDelay
}
