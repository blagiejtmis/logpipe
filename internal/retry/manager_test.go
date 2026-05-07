package retry

import (
	"testing"
	"time"
)

func makeCfg(defaultAttempts int, overrides map[string]int) *SinkConfig {
	sc := &SinkConfig{
		Default: &Config{
			MaxAttempts:  defaultAttempts,
			InitialDelay: time.Millisecond,
			MaxDelay:     10 * time.Millisecond,
			Multiplier:   2.0,
		},
		Overrides: make(map[string]Config),
	}
	for name, attempts := range overrides {
		sc.Overrides[name] = Config{
			MaxAttempts:  attempts,
			InitialDelay: time.Millisecond,
			MaxDelay:     10 * time.Millisecond,
			Multiplier:   2.0,
		}
	}
	return sc
}

func TestNewManager_NilConfig_UsesDefaults(t *testing.T) {
	m, err := NewManager(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := m.For("any-sink")
	if r == nil {
		t.Fatal("expected non-nil retryer")
	}
	if r.cfg.MaxAttempts != defaultConfig.MaxAttempts {
		t.Fatalf("expected %d attempts, got %d", defaultConfig.MaxAttempts, r.cfg.MaxAttempts)
	}
}

func TestNewManager_DefaultApplied(t *testing.T) {
	m, err := NewManager(makeCfg(5, nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := m.For("stdout")
	if r.cfg.MaxAttempts != 5 {
		t.Fatalf("expected 5 attempts, got %d", r.cfg.MaxAttempts)
	}
}

func TestNewManager_OverrideForSink(t *testing.T) {
	m, err := NewManager(makeCfg(3, map[string]int{"file-sink": 7}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.For("file-sink").cfg.MaxAttempts != 7 {
		t.Fatal("expected override to apply")
	}
	if m.For("stdout").cfg.MaxAttempts != 3 {
		t.Fatal("expected default for unoverridden sink")
	}
}

func TestNewManager_InvalidDefault_ReturnsError(t *testing.T) {
	_, err := NewManager(&SinkConfig{
		Default: &Config{MaxAttempts: 0},
	})
	if err == nil {
		t.Fatal("expected error for invalid default config")
	}
}

func TestNewManager_InvalidOverride_ReturnsError(t *testing.T) {
	_, err := NewManager(&SinkConfig{
		Default: &Config{MaxAttempts: 1},
		Overrides: map[string]Config{
			"bad-sink": {MaxAttempts: -1},
		},
	})
	if err == nil {
		t.Fatal("expected error for invalid override config")
	}
}
