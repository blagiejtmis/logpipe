package timeseries

import (
	"testing"
	"time"
)

func TestNewManager_NilConfig_AllowsAll(t *testing.T) {
	m, err := NewManager(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m.Record("src", map[string]any{"level": "info"})
	counts := m.Counts("src")
	if len(counts) != 0 {
		t.Errorf("expected empty counts with nil config, got %v", counts)
	}
}

func TestNewManager_GlobalField_Applied(t *testing.T) {
	cfg := &Config{Field: "level", Window: time.Minute, Buckets: 6}
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m.Record("any", map[string]any{"level": "error"})
	m.Record("any", map[string]any{"level": "error"})
	counts := m.Counts("any")
	if counts["error"] != 2 {
		t.Errorf("expected error=2, got %d", counts["error"])
	}
}

func TestNewManager_SourceSpecific_OverridesGlobal(t *testing.T) {
	cfg := &Config{
		Field:   "level",
		Window:  time.Minute,
		Buckets: 6,
		SourceOverride: map[string]string{
			"app": "severity",
		},
	}
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m.Record("app", map[string]any{"severity": "warn"})
	counts := m.Counts("app")
	if counts["warn"] != 1 {
		t.Errorf("expected warn=1, got %d", counts["warn"])
	}
	// global field "level" should not be counted for "app"
	if counts["level"] != 0 {
		t.Errorf("unexpected level count for app source")
	}
}

func TestNewManager_DefaultWindowAndBuckets_Applied(t *testing.T) {
	cfg := &Config{Field: "level"} // zero Window/Buckets → defaults
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m.Record("x", map[string]any{"level": "debug"})
	if c := m.Counts("x")["debug"]; c != 1 {
		t.Errorf("expected 1, got %d", c)
	}
}

func TestNewManager_InvalidSourceField_ReturnsError(t *testing.T) {
	cfg := &Config{
		Field:          "level",
		Window:         time.Minute,
		Buckets:        6,
		SourceOverride: map[string]string{"app": ""},
	}
	_, err := NewManager(cfg)
	if err == nil {
		t.Fatal("expected error for empty source field override")
	}
}
