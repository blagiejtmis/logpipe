package throttle

import (
	"testing"
	"time"
)

func TestNew_InvalidRate_ReturnsError(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for rate=0")
	}
}

func TestNew_ValidRate(t *testing.T) {
	th, err := New(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if th == nil {
		t.Fatal("expected non-nil throttler")
	}
}

func TestWait_UnderBudget_DoesNotSleep(t *testing.T) {
	th, _ := New(5)
	slept := false
	th.sleepFn = func(d time.Duration) { slept = true }

	for i := 0; i < 5; i++ {
		th.Wait()
	}
	if slept {
		t.Error("expected no sleep when under budget")
	}
}

func TestWait_ExceedsBudget_Sleeps(t *testing.T) {
	th, _ := New(3)
	var sleptFor time.Duration
	th.sleepFn = func(d time.Duration) { sleptFor = d }

	for i := 0; i < 4; i++ {
		th.Wait()
	}
	if sleptFor <= 0 {
		t.Error("expected sleep when budget exceeded")
	}
}

func TestWait_WindowReset_AllowsNewBudget(t *testing.T) {
	th, _ := New(2)
	th.sleepFn = func(d time.Duration) {}

	// Exhaust budget
	th.Wait()
	th.Wait()

	// Simulate window expiry by back-dating windowAt
	th.mu.Lock()
	th.windowAt = time.Now().Add(-2 * time.Second)
	th.mu.Unlock()

	slept := false
	th.sleepFn = func(d time.Duration) { slept = true }

	th.Wait()
	if slept {
		t.Error("expected no sleep after window reset")
	}
}
