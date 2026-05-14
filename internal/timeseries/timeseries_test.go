package timeseries

import (
	"testing"
	"time"
)

func TestNew_EmptyField_ReturnsError(t *testing.T) {
	_, err := New("", time.Minute, 6)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_ZeroWindow_ReturnsError(t *testing.T) {
	_, err := New("level", 0, 6)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNew_ZeroBuckets_ReturnsError(t *testing.T) {
	_, err := New("level", time.Minute, 0)
	if err == nil {
		t.Fatal("expected error for zero buckets")
	}
}

func TestNew_ValidArgs_NoError(t *testing.T) {
	_, err := New("level", time.Minute, 6)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRecord_CountsValues(t *testing.T) {
	s, _ := New("level", time.Minute, 6)
	s.Record(map[string]any{"level": "info"})
	s.Record(map[string]any{"level": "info"})
	s.Record(map[string]any{"level": "warn"})

	counts := s.Counts()
	if counts["info"] != 2 {
		t.Errorf("expected info=2, got %d", counts["info"])
	}
	if counts["warn"] != 1 {
		t.Errorf("expected warn=1, got %d", counts["warn"])
	}
}

func TestRecord_MissingField_CountsEmpty(t *testing.T) {
	s, _ := New("level", time.Minute, 6)
	s.Record(map[string]any{"msg": "hello"})

	counts := s.Counts()
	if counts[""] != 1 {
		t.Errorf("expected empty-key count=1, got %d", counts[""])
	}
}

func TestCounts_EmptySeries_ReturnsEmptyMap(t *testing.T) {
	s, _ := New("level", time.Minute, 6)
	counts := s.Counts()
	if len(counts) != 0 {
		t.Errorf("expected empty map, got %v", counts)
	}
}

func TestRotate_ExpiredBuckets_Cleared(t *testing.T) {
	// Use a very short window so we can force rotation.
	s, _ := New("level", 20*time.Millisecond, 2)
	s.Record(map[string]any{"level": "info"})

	time.Sleep(30 * time.Millisecond)

	counts := s.Counts()
	if counts["info"] != 0 {
		t.Errorf("expected info to be expired, got %d", counts["info"])
	}
}
