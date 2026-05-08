package ratelimit

import (
	"testing"
	"time"
)

func TestNew_InvalidArgs(t *testing.T) {
	_, err := New(0, time.Second)
	if err == nil {
		t.Fatal("expected error for maxLines=0")
	}
	_, err = New(10, 0)
	if err == nil {
		t.Fatal("expected error for window=0")
	}
}

func TestAllow_UnderLimit(t *testing.T) {
	l, err := New(5, time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 5; i++ {
		if !l.Allow("src") {
			t.Fatalf("expected Allow=true on call %d", i+1)
		}
	}
}

func TestAllow_ExceedsLimit(t *testing.T) {
	l, _ := New(3, time.Second)
	for i := 0; i < 3; i++ {
		l.Allow("src")
	}
	if l.Allow("src") {
		t.Fatal("expected Allow=false after limit exceeded")
	}
}

func TestAllow_IndependentSources(t *testing.T) {
	l, _ := New(2, time.Second)
	l.Allow("a")
	l.Allow("a")
	// source "a" is now at limit; "b" should still be allowed
	if !l.Allow("b") {
		t.Fatal("expected Allow=true for independent source b")
	}
	// source "a" should be blocked
	if l.Allow("a") {
		t.Fatal("expected Allow=false for source a at limit")
	}
}

func TestAllow_WindowReset(t *testing.T) {
	l, _ := New(1, 50*time.Millisecond)
	if !l.Allow("src") {
		t.Fatal("first call should be allowed")
	}
	if l.Allow("src") {
		t.Fatal("second call in same window should be blocked")
	}
	time.Sleep(60 * time.Millisecond)
	if !l.Allow("src") {
		t.Fatal("call after window expiry should be allowed")
	}
}

func TestReset_ClearsState(t *testing.T) {
	l, _ := New(1, time.Second)
	l.Allow("src")
	if l.Allow("src") {
		t.Fatal("should be blocked before reset")
	}
	l.Reset()
	if !l.Allow("src") {
		t.Fatal("should be allowed after reset")
	}
}

func TestAllow_MultipleSourcesAfterReset(t *testing.T) {
	l, _ := New(2, time.Second)
	l.Allow("x")
	l.Allow("x")
	l.Allow("y")
	l.Reset()
	// After reset, all sources should be cleared and allowed again.
	if !l.Allow("x") {
		t.Fatal("expected Allow=true for source x after reset")
	}
	if !l.Allow("y") {
		t.Fatal("expected Allow=true for source y after reset")
	}
}
