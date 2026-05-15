package classify

import (
	"testing"
)

func makeCfg() *Config {
	return &Config{
		Default: []Rule{
			{Field: "level", Pattern: "error", Label: "fault"},
		},
		Sources: map[string][]Rule{
			"app": {
				{Field: "message", Pattern: "login", Label: "auth"},
			},
		},
	}
}

func TestNewManager_NilConfig_AllowsAll(t *testing.T) {
	m, err := NewManager(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rec := map[string]any{"level": "error"}
	out := m.Apply("any", rec)
	if _, exists := out["category"]; exists {
		t.Fatal("expected no classification for nil config")
	}
}

func TestNewManager_DefaultRule_Applied(t *testing.T) {
	m, err := NewManager(makeCfg())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rec := map[string]any{"level": "error", "message": "disk full"}
	out := m.Apply("unknown-source", rec)
	if got := out["category"]; got != "fault" {
		t.Fatalf("expected fault, got %v", got)
	}
}

func TestNewManager_SourceSpecific_OverridesDefault(t *testing.T) {
	m, err := NewManager(makeCfg())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rec := map[string]any{"level": "error", "message": "user login ok"}
	out := m.Apply("app", rec)
	if got := out["category"]; got != "auth" {
		t.Fatalf("expected auth from source rule, got %v", got)
	}
}

func TestNewManager_InvalidDefaultRule_ReturnsError(t *testing.T) {
	cfg := &Config{
		Default: []Rule{
			{Field: "level", Pattern: "[", Label: "bad"},
		},
	}
	_, err := NewManager(cfg)
	if err == nil {
		t.Fatal("expected error for invalid default pattern")
	}
}

func TestNewManager_InvalidSourceRule_ReturnsError(t *testing.T) {
	cfg := &Config{
		Sources: map[string][]Rule{
			"svc": {{Field: "msg", Pattern: "[", Label: "x"}},
		},
	}
	_, err := NewManager(cfg)
	if err == nil {
		t.Fatal("expected error for invalid source pattern")
	}
}

func TestNewManager_NoMatch_RecordUnchanged(t *testing.T) {
	m, err := NewManager(makeCfg())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rec := map[string]any{"level": "info", "message": "heartbeat"}
	out := m.Apply("unknown-source", rec)
	if _, exists := out["category"]; exists {
		t.Fatalf("expected no category, record should be unchanged")
	}
}
