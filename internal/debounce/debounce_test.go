package debounce

import (
	"testing"
	"time"
)

func base() map[string]any {
	return map[string]any{"host": "web-01", "level": "error"}
}

func TestNew_EmptyField_ReturnsError(t *testing.T) {
	_, err := New("", time.Second)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_ZeroWindow_ReturnsError(t *testing.T) {
	_, err := New("host", 0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNew_NegativeWindow_ReturnsError(t *testing.T) {
	_, err := New("host", -time.Second)
	if err == nil {
		t.Fatal("expected error for negative window")
	}
}

func TestNew_ValidArgs_NoError(t *testing.T) {
	_, err := New("host", time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAllow_FirstOccurrence_Passes(t *testing.T) {
	d, _ := New("host", time.Second)
	if !d.Allow(base()) {
		t.Fatal("expected first occurrence to pass")
	}
}

func TestAllow_RepeatWithinWindow_Drops(t *testing.T) {
	now := time.Now()
	d, _ := New("host", time.Second)
	d.nowFunc = func() time.Time { return now }

	if !d.Allow(base()) {
		t.Fatal("expected first to pass")
	}
	if d.Allow(base()) {
		t.Fatal("expected repeat within window to be dropped")
	}
}

func TestAllow_AfterWindowExpires_Passes(t *testing.T) {
	now := time.Now()
	d, _ := New("host", 100*time.Millisecond)
	d.nowFunc = func() time.Time { return now }
	d.Allow(base())

	d.nowFunc = func() time.Time { return now.Add(200 * time.Millisecond) }
	if !d.Allow(base()) {
		t.Fatal("expected record to pass after window expiry")
	}
}

func TestAllow_DifferentFieldValues_IndependentState(t *testing.T) {
	d, _ := New("host", time.Second)
	r1 := map[string]any{"host": "web-01"}
	r2 := map[string]any{"host": "web-02"}

	d.Allow(r1)
	if !d.Allow(r2) {
		t.Fatal("expected different key to pass independently")
	}
}

func TestAllow_MissingField_TreatedAsEmpty(t *testing.T) {
	d, _ := New("host", time.Second)
	r := map[string]any{"level": "info"} // no "host"
	if !d.Allow(r) {
		t.Fatal("expected first missing-field record to pass")
	}
	if d.Allow(r) {
		t.Fatal("expected second missing-field record to be dropped")
	}
}

func TestPurge_RemovesExpiredEntries(t *testing.T) {
	now := time.Now()
	d, _ := New("host", 100*time.Millisecond)
	d.nowFunc = func() time.Time { return now }
	d.Allow(base())

	d.nowFunc = func() time.Time { return now.Add(200 * time.Millisecond) }
	d.Purge()

	if len(d.seen) != 0 {
		t.Fatalf("expected empty seen map after purge, got %d entries", len(d.seen))
	}
}
