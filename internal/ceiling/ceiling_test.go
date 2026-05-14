package ceiling

import (
	"testing"
	"time"
)

func TestNew_ZeroMax_ReturnsError(t *testing.T) {
	_, err := New(0, time.Second)
	if err == nil {
		t.Fatal("expected error for zero max")
	}
}

func TestNew_ZeroWindow_ReturnsError(t *testing.T) {
	_, err := New(5, 0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNew_ValidArgs_NoError(t *testing.T) {
	_, err := New(10, time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAllow_UnderCeiling_Accepts(t *testing.T) {
	c, _ := New(3, time.Minute)
	for i := 0; i < 3; i++ {
		if !c.Allow() {
			t.Fatalf("expected Allow()=true on call %d", i+1)
		}
	}
}

func TestAllow_AtCeiling_Drops(t *testing.T) {
	c, _ := New(2, time.Minute)
	c.Allow()
	c.Allow()
	if c.Allow() {
		t.Fatal("expected Allow()=false after ceiling reached")
	}
}

func TestAllow_WindowReset_AllowsAgain(t *testing.T) {
	now := time.Now()
	c, _ := New(1, 100*time.Millisecond)
	c.now = func() time.Time { return now }

	if !c.Allow() {
		t.Fatal("first call should be allowed")
	}
	if c.Allow() {
		t.Fatal("second call within window should be dropped")
	}

	// Advance clock past the window.
	c.now = func() time.Time { return now.Add(200 * time.Millisecond) }
	if !c.Allow() {
		t.Fatal("first call after window reset should be allowed")
	}
}

func TestReset_ClearsCount(t *testing.T) {
	c, _ := New(1, time.Minute)
	c.Allow()
	if c.Allow() {
		t.Fatal("should be dropped before reset")
	}
	c.Reset()
	if !c.Allow() {
		t.Fatal("should be allowed after reset")
	}
}
