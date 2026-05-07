package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

var errTemp = errors.New("temporary error")

func fastCfg(attempts int) Config {
	return Config{
		MaxAttempts:  attempts,
		InitialDelay: time.Millisecond,
		MaxDelay:     10 * time.Millisecond,
		Multiplier:   2.0,
	}
}

func TestNew_InvalidMaxAttempts(t *testing.T) {
	_, err := New(Config{MaxAttempts: 0})
	if err == nil {
		t.Fatal("expected error for max_attempts=0")
	}
}

func TestNew_InvalidInitialDelay(t *testing.T) {
	_, err := New(Config{MaxAttempts: 1, InitialDelay: -1})
	if err == nil {
		t.Fatal("expected error for negative initial_delay")
	}
}

func TestDo_SucceedsOnFirstAttempt(t *testing.T) {
	r, _ := New(fastCfg(3))
	calls := 0
	err := r.Do(context.Background(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesAndSucceeds(t *testing.T) {
	r, _ := New(fastCfg(3))
	calls := 0
	err := r.Do(context.Background(), func() error {
		calls++
		if calls < 3 {
			return errTemp
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ExhaustsAttempts(t *testing.T) {
	r, _ := New(fastCfg(2))
	err := r.Do(context.Background(), func() error { return errTemp })
	if !errors.Is(err, ErrMaxAttemptsReached) {
		t.Fatalf("expected ErrMaxAttemptsReached, got %v", err)
	}
	if !errors.Is(err, errTemp) {
		t.Fatalf("expected wrapped errTemp, got %v", err)
	}
}

func TestDo_ContextCancelled(t *testing.T) {
	r, _ := New(fastCfg(5))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := r.Do(ctx, func() error { return nil })
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestDo_DefaultMultiplierApplied(t *testing.T) {
	cfg := Config{MaxAttempts: 1}
	r, err := New(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.cfg.Multiplier != 2.0 {
		t.Fatalf("expected default multiplier 2.0, got %v", r.cfg.Multiplier)
	}
}
