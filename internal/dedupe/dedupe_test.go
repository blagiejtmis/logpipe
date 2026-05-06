package dedupe

import (
	"testing"
	"time"
)

func TestNew_InvalidArgs(t *testing.T) {
	_, err := New(nil, time.Second)
	if err == nil {
		t.Fatal("expected error for empty fields")
	}
	_, err = New([]string{"msg"}, 0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestAllow_FirstOccurrence(t *testing.T) {
	d, _ := New([]string{"msg"}, time.Minute)
	rec := Record{"msg": "hello"}
	if !d.Allow(rec) {
		t.Fatal("first occurrence should be allowed")
	}
}

func TestAllow_DuplicateWithinWindow(t *testing.T) {
	now := time.Now()
	d, _ := New([]string{"msg"}, time.Minute)
	d.nowFn = func() time.Time { return now }

	rec := Record{"msg": "dup"}
	d.Allow(rec)
	if d.Allow(rec) {
		t.Fatal("duplicate within window should be rejected")
	}
}

func TestAllow_DuplicateAfterWindowExpires(t *testing.T) {
	now := time.Now()
	d, _ := New([]string{"msg"}, 5*time.Second)
	d.nowFn = func() time.Time { return now }

	rec := Record{"msg": "expire"}
	d.Allow(rec)

	// advance past the window
	d.nowFn = func() time.Time { return now.Add(6 * time.Second) }
	if !d.Allow(rec) {
		t.Fatal("record after window expiry should be allowed")
	}
}

func TestAllow_DifferentFieldValues(t *testing.T) {
	d, _ := New([]string{"msg"}, time.Minute)
	if !d.Allow(Record{"msg": "a"}) {
		t.Fatal("first record should be allowed")
	}
	if !d.Allow(Record{"msg": "b"}) {
		t.Fatal("different value should be allowed")
	}
}

func TestPurge_RemovesExpired(t *testing.T) {
	now := time.Now()
	d, _ := New([]string{"msg"}, 5*time.Second)
	d.nowFn = func() time.Time { return now }

	d.Allow(Record{"msg": "old"})

	d.nowFn = func() time.Time { return now.Add(10 * time.Second) }
	d.Purge()

	d.mu.Lock()
	n := len(d.seen)
	d.mu.Unlock()
	if n != 0 {
		t.Fatalf("expected 0 entries after purge, got %d", n)
	}
}
