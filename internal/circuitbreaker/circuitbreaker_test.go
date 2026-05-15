package circuitbreaker

import (
	"testing"
	"time"
)

func TestNew_InvalidMaxFailures(t *testing.T) {
	_, err := New(0, time.Second)
	if err == nil {
		t.Fatal("expected error for zero maxFailures")
	}
}

func TestNew_InvalidCooldown(t *testing.T) {
	_, err := New(3, 0)
	if err == nil {
		t.Fatal("expected error for zero cooldown")
	}
}

func TestNew_ValidArgs_NoError(t *testing.T) {
	_, err := New(3, time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAllow_ClosedCircuit_Passes(t *testing.T) {
	b, _ := New(3, time.Second)
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestAllow_OpensAfterMaxFailures(t *testing.T) {
	b, _ := New(3, time.Second)
	for i := 0; i < 3; i++ {
		b.RecordFailure()
	}
	if err := b.Allow(); err != ErrOpen {
		t.Fatalf("expected ErrOpen, got %v", err)
	}
}

func TestAllow_DoesNotOpenBeforeThreshold(t *testing.T) {
	b, _ := New(3, time.Second)
	b.RecordFailure()
	b.RecordFailure()
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil before threshold, got %v", err)
	}
}

func TestRecordSuccess_ResetsClosed(t *testing.T) {
	b, _ := New(2, time.Second)
	b.RecordFailure()
	b.RecordFailure() // opens
	// simulate cooldown elapsed by manipulating openedAt
	b.mu.Lock()
	b.openedAt = time.Now().Add(-2 * time.Second)
	b.mu.Unlock()
	// probe allowed
	if err := b.Allow(); err != nil {
		t.Fatalf("expected probe to be allowed, got %v", err)
	}
	b.RecordSuccess()
	if b.CurrentState() != StateClosed {
		t.Fatal("expected state to be closed after success")
	}
}

func TestAllow_HalfOpen_ResetOnCooldown(t *testing.T) {
	b, _ := New(1, 50*time.Millisecond)
	b.RecordFailure()
	if err := b.Allow(); err != ErrOpen {
		t.Fatalf("expected ErrOpen immediately after open, got %v", err)
	}
	time.Sleep(60 * time.Millisecond)
	if err := b.Allow(); err != nil {
		t.Fatalf("expected probe allowed after cooldown, got %v", err)
	}
}

func TestCurrentState_ReflectsState(t *testing.T) {
	b, _ := New(2, time.Second)
	if b.CurrentState() != StateClosed {
		t.Fatal("expected initial state to be closed")
	}
	b.RecordFailure()
	b.RecordFailure()
	if b.CurrentState() != StateOpen {
		t.Fatal("expected state to be open after threshold")
	}
}
