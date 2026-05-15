// Package circuitbreaker implements a per-sink circuit breaker that opens
// after a configurable number of consecutive failures and resets after a
// cooldown window, preventing log loss from cascading sink errors.
package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

// ErrOpen is returned by Allow when the circuit is open.
var ErrOpen = errors.New("circuit breaker is open")

// State represents the current circuit state.
type State int

const (
	StateClosed State = iota
	StateOpen
)

// Breaker is a simple consecutive-failure circuit breaker.
type Breaker struct {
	mu          sync.Mutex
	maxFailures int
	cooldown    time.Duration
	failures    int
	openedAt    time.Time
	state       State
}

// New creates a Breaker that opens after maxFailures consecutive failures
// and attempts a reset after cooldown has elapsed.
func New(maxFailures int, cooldown time.Duration) (*Breaker, error) {
	if maxFailures <= 0 {
		return nil, errors.New("circuitbreaker: maxFailures must be > 0")
	}
	if cooldown <= 0 {
		return nil, errors.New("circuitbreaker: cooldown must be > 0")
	}
	return &Breaker{
		maxFailures: maxFailures,
		cooldown:    cooldown,
	}, nil
}

// Allow returns nil if the circuit is closed (or half-open for a probe).
// It returns ErrOpen when the circuit is open and the cooldown has not elapsed.
func (b *Breaker) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.state == StateOpen {
		if time.Since(b.openedAt) >= b.cooldown {
			// Half-open: allow one probe attempt.
			b.state = StateClosed
			b.failures = 0
			return nil
		}
		return ErrOpen
	}
	return nil
}

// RecordSuccess resets the failure counter.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.state = StateClosed
}

// RecordFailure increments the failure counter and opens the circuit if the
// threshold is reached.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	if b.failures >= b.maxFailures {
		b.state = StateOpen
		b.openedAt = time.Now()
	}
}

// CurrentState returns the current State without side-effects.
func (b *Breaker) CurrentState() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}
