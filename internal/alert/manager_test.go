package alert

import (
	"testing"

	"github.com/logpipe/logpipe/internal/metrics"
)

func makeAlertCfg(field, pattern string, threshold int, cb func(string, string, int)) *ManagerConfig {
	return &ManagerConfig{
		Global: []AlertRule{
			{Field: field, Pattern: pattern, Threshold: threshold, Callback: cb},
		},
	}
}

func TestNewManager_NilConfig_ReturnsEmptyManager(t *testing.T) {
	m, err := NewManager(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.AlerterFor("any") != nil {
		t.Fatal("expected nil alerter for unconfigured manager")
	}
}

func TestNewManager_GlobalRule_Applied(t *testing.T) {
	fired := false
	cb := func(_, _ string, _ int) { fired = true }
	cfg := makeAlertCfg("level", "error", 1, cb)

	m, err := NewManager(cfg, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	a := m.AlerterFor("some-source")
	if a == nil {
		t.Fatal("expected global alerter to be returned")
	}

	rec := map[string]any{"level": "error"}
	a.Evaluate(rec)
	if !fired {
		t.Fatal("expected callback to fire")
	}
}

func TestNewManager_SourceSpecific_OverridesGlobal(t *testing.T) {
	globalFired, srcFired := false, false
	cfg := &ManagerConfig{
		Global: []AlertRule{
			{Field: "level", Pattern: "error", Threshold: 1, Callback: func(_, _ string, _ int) { globalFired = true }},
		},
		Sources: map[string][]AlertRule{
			"app": {
				{Field: "level", Pattern: "error", Threshold: 1, Callback: func(_, _ string, _ int) { srcFired = true }},
			},
		},
	}

	m, err := NewManager(cfg, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	a := m.AlerterFor("app")
	if a == nil {
		t.Fatal("expected source-specific alerter")
	}
	rec := map[string]any{"level": "error"}
	a.Evaluate(rec)

	if !srcFired {
		t.Fatal("expected source callback to fire")
	}
	if globalFired {
		t.Fatal("global callback should not fire for source-specific alerter")
	}
}

func TestNewManager_InvalidRule_ReturnsError(t *testing.T) {
	cfg := &ManagerConfig{
		Global: []AlertRule{
			{Field: "level", Pattern: "[", Threshold: 1, Callback: func(_, _ string, _ int) {}},
		},
	}
	_, err := NewManager(cfg, nil)
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestNewManager_WithRegistry(t *testing.T) {
	reg := metrics.NewRegistry()
	cb := func(_, _ string, _ int) {}
	cfg := makeAlertCfg("level", "error", 1, cb)

	_, err := NewManager(cfg, reg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
