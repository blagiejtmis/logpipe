package debounce

import (
	"testing"
	"time"
)

func makeCfg(field string, window time.Duration) *Config {
	return &Config{
		Default: &Rule{Field: field, Window: window},
	}
}

func TestNewManager_NilConfig_AllowsAll(t *testing.T) {
	m, err := NewManager(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !m.Allow("src", map[string]any{"host": "x"}) {
		t.Fatal("expected nil config to allow all records")
	}
}

func TestNewManager_DefaultRule_Applied(t *testing.T) {
	m, err := NewManager(makeCfg("host", time.Second))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := map[string]any{"host": "web-01"}
	if !m.Allow("any", r) {
		t.Fatal("expected first occurrence to pass")
	}
	if m.Allow("any", r) {
		t.Fatal("expected duplicate within window to be dropped")
	}
}

func TestNewManager_SourceSpecific_OverridesDefault(t *testing.T) {
	cfg := &Config{
		Default: &Rule{Field: "host", Window: time.Hour},
		Sources: map[string]*Rule{
			"svc-a": {Field: "host", Window: 50 * time.Millisecond},
		},
	}
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r := map[string]any{"host": "web-01"}
	m.Allow("svc-a", r)

	// Advance time past the source-specific window via the debouncer's clock.
	if d, ok := m.sources["svc-a"]; ok {
		now := time.Now().Add(200 * time.Millisecond)
		d.nowFunc = func() time.Time { return now }
	}

	if !m.Allow("svc-a", r) {
		t.Fatal("expected record to pass after source-specific window expired")
	}
}

func TestNewManager_InvalidDefaultRule_ReturnsError(t *testing.T) {
	_, err := NewManager(makeCfg("", time.Second))
	if err == nil {
		t.Fatal("expected error for invalid default rule")
	}
}

func TestNewManager_InvalidSourceRule_ReturnsError(t *testing.T) {
	cfg := &Config{
		Sources: map[string]*Rule{
			"bad": {Field: "", Window: time.Second},
		},
	}
	_, err := NewManager(cfg)
	if err == nil {
		t.Fatal("expected error for invalid source rule")
	}
}
