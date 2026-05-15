package uniq

import (
	"testing"
	"time"
)

func TestNew_EmptyField_ReturnsError(t *testing.T) {
	_, err := New("", time.Minute)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_ZeroWindow_ReturnsError(t *testing.T) {
	_, err := New("id", 0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNew_ValidArgs_NoError(t *testing.T) {
	_, err := New("id", time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAllow_FirstOccurrence_Passes(t *testing.T) {
	u, _ := New("id", time.Minute)
	r := Record{"id": "abc"}
	if !u.Allow(r) {
		t.Fatal("expected first occurrence to be allowed")
	}
}

func TestAllow_DuplicateWithinWindow_Drops(t *testing.T) {
	u, _ := New("id", time.Minute)
	r := Record{"id": "abc"}
	u.Allow(r)
	if u.Allow(r) {
		t.Fatal("expected duplicate within window to be dropped")
	}
}

func TestAllow_DuplicateAfterWindowExpires_Passes(t *testing.T) {
	now := time.Now()
	u, _ := New("id", 100*time.Millisecond)
	u.now = func() time.Time { return now }

	r := Record{"id": "abc"}
	u.Allow(r)

	// Advance clock past the window.
	u.now = func() time.Time { return now.Add(200 * time.Millisecond) }
	if !u.Allow(r) {
		t.Fatal("expected record to be allowed after window expires")
	}
}

func TestAllow_MissingField_Passes(t *testing.T) {
	u, _ := New("id", time.Minute)
	if !u.Allow(Record{"other": "val"}) {
		t.Fatal("expected record with missing key field to be allowed")
	}
}

func TestAllow_NonStringField_Passes(t *testing.T) {
	u, _ := New("id", time.Minute)
	r := Record{"id": 42}
	if !u.Allow(r) {
		t.Fatal("expected record with non-string key field to be allowed")
	}
}

func TestAllow_DifferentValues_BothPass(t *testing.T) {
	u, _ := New("id", time.Minute)
	if !u.Allow(Record{"id": "a"}) {
		t.Fatal("expected 'a' to be allowed")
	}
	if !u.Allow(Record{"id": "b"}) {
		t.Fatal("expected 'b' to be allowed")
	}
}
