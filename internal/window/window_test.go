package window

import (
	"testing"
	"time"
)

func TestNew_InvalidDuration_ReturnsError(t *testing.T) {
	_, err := New(0, 10)
	if err == nil {
		t.Fatal("expected error for zero duration")
	}
}

func TestNew_InvalidBuckets_ReturnsError(t *testing.T) {
	_, err := New(time.Second, 0)
	if err == nil {
		t.Fatal("expected error for zero buckets")
	}
}

func TestNew_ValidArgs_NoError(t *testing.T) {
	c, err := New(time.Second, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil counter")
	}
}

func TestAdd_IncreasesCount(t *testing.T) {
	c, _ := New(time.Second, 10)
	c.Add(3)
	c.Add(7)
	if got := c.Count(); got != 10 {
		t.Fatalf("expected 10, got %d", got)
	}
}

func TestCount_ZeroInitially(t *testing.T) {
	c, _ := New(time.Second, 10)
	if got := c.Count(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestCount_ExpiresAfterWindow(t *testing.T) {
	// Use a very short window so we can expire it quickly.
	c, _ := New(50*time.Millisecond, 5)
	c.Add(42)
	if got := c.Count(); got != 42 {
		t.Fatalf("pre-expiry: expected 42, got %d", got)
	}
	time.Sleep(80 * time.Millisecond)
	if got := c.Count(); got != 0 {
		t.Fatalf("post-expiry: expected 0, got %d", got)
	}
}

func TestAdd_ConcurrentSafe(t *testing.T) {
	c, _ := New(time.Second, 10)
	done := make(chan struct{})
	for i := 0; i < 50; i++ {
		go func() {
			c.Add(1)
			done <- struct{}{}
		}()
	}
	for i := 0; i < 50; i++ {
		<-done
	}
	if got := c.Count(); got != 50 {
		t.Fatalf("expected 50, got %d", got)
	}
}
