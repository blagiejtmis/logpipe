// Package retry provides configurable retry logic for sink writes.
package retry

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// ErrMaxAttemptsReached is returned when all retry attempts are exhausted.
var ErrMaxAttemptsReached = errors.New("retry: max attempts reached")

// Config holds retry policy parameters.
type Config struct {
	MaxAttempts int           `yaml:"max_attempts"`
	InitialDelay time.Duration `yaml:"initial_delay"`
	MaxDelay     time.Duration `yaml:"max_delay"`
	Multiplier   float64       `yaml:"multiplier"`
}

// Retryer executes a function with exponential back-off.
type Retryer struct {
	cfg Config
}

// New creates a Retryer. Returns an error if the config is invalid.
func New(cfg Config) (*Retryer, error) {
	if cfg.MaxAttempts < 1 {
		return nil, fmt.Errorf("retry: max_attempts must be >= 1, got %d", cfg.MaxAttempts)
	}
	if cfg.InitialDelay < 0 {
		return nil, fmt.Errorf("retry: initial_delay must be >= 0")
	}
	if cfg.Multiplier == 0 {
		cfg.Multiplier = 2.0
	}
	if cfg.MaxDelay == 0 {
		cfg.MaxDelay = 30 * time.Second
	}
	return &Retryer{cfg: cfg}, nil
}

// Do calls fn up to MaxAttempts times, backing off between attempts.
// It stops early if ctx is cancelled.
func (r *Retryer) Do(ctx context.Context, fn func() error) error {
	delay := r.cfg.InitialDelay
	var lastErr error
	for attempt := 1; attempt <= r.cfg.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		if attempt == r.cfg.MaxAttempts {
			break
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
		delay = time.Duration(float64(delay) * r.cfg.Multiplier)
		if delay > r.cfg.MaxDelay {
			delay = r.cfg.MaxDelay
		}
	}
	return fmt.Errorf("%w: %w", ErrMaxAttemptsReached, lastErr)
}

// IsMaxAttemptsReached reports whether err was caused by exhausting all retry
// attempts, as opposed to context cancellation or another sentinel error.
func IsMaxAttemptsReached(err error) bool {
	return errors.Is(err, ErrMaxAttemptsReached)
}
