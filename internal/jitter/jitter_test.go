package jitter

import (
	"testing"
	"time"
)

func TestNew_NegativeMin_ReturnsError(t *testing.T) {
	_, err := New(-time.Millisecond, time.Millisecond)
	if err == nil {
		t.Fatal("expected error for negative minDelay")
	}
}

func TestNew_NegativeMax_ReturnsError(t *testing.T) {
	_, err := New(0, -time.Millisecond)
	if err == nil {
		t.Fatal("expected error for negative maxDelay")
	}
}

func TestNew_MinGreaterThanMax_ReturnsError(t *testing.T) {
	_, err := New(10*time.Millisecond, 5*time.Millisecond)
	if err == nil {
		t.Fatal("expected error when min > max")
	}
}

func TestNew_ValidArgs_NoError(t *testing.T) {
	j, err := New(0, 10*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if j == nil {
		t.Fatal("expected non-nil Jitter")
	}
}

func TestNew_EqualMinMax_NoError(t *testing.T) {
	j, err := New(5*time.Millisecond, 5*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	min, max := j.Range()
	if min != max {
		t.Fatalf("expected equal range, got %s %s", min, max)
	}
}

func TestWait_ZeroDelay_DoesNotBlock(t *testing.T) {
	j, _ := New(0, 0)
	done := make(chan struct{})
	go func() {
		j.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Wait blocked longer than expected for zero delay")
	}
}

func TestWait_SmallRange_CompletesQuickly(t *testing.T) {
	j, _ := New(0, 5*time.Millisecond)
	start := time.Now()
	j.Wait()
	if elapsed := time.Since(start); elapsed > 200*time.Millisecond {
		t.Fatalf("Wait took too long: %s", elapsed)
	}
}

func TestRange_ReturnsConfiguredValues(t *testing.T) {
	min := 2 * time.Millisecond
	max := 8 * time.Millisecond
	j, _ := New(min, max)
	gotMin, gotMax := j.Range()
	if gotMin != min || gotMax != max {
		t.Fatalf("Range() = (%s, %s), want (%s, %s)", gotMin, gotMax, min, max)
	}
}
