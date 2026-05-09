package coalesce

import (
	"testing"
)

func makeCfg(global []RuleConfig, sources map[string][]RuleConfig) *Config {
	return &Config{Global: global, Sources: sources}
}

func TestNewManager_NilConfig_AllowsAll(t *testing.T) {
	m, err := NewManager(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Get("any") != nil {
		t.Error("expected nil coalescer for unconfigured manager")
	}
}

func TestNewManager_GlobalRule(t *testing.T) {
	cfg := makeCfg([]RuleConfig{
		{Sources: []string{"msg", "message"}, Dest: "message"},
	}, nil)
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Get("any") == nil {
		t.Error("expected global coalescer to be returned")
	}
}

func TestNewManager_SourceSpecific_OverridesGlobal(t *testing.T) {
	cfg := makeCfg(
		[]RuleConfig{{Sources: []string{"a", "b"}, Dest: "out"}},
		map[string][]RuleConfig{
			"myapp": {{Sources: []string{"x", "y"}, Dest: "result"}},
		},
	)
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	specific := m.Get("myapp")
	if specific == nil {
		t.Fatal("expected source-specific coalescer")
	}
	rec := map[string]any{"x": "val", "y": "other"}
	specific.Apply(rec)
	if rec["result"] != "val" {
		t.Errorf("expected 'val', got %v", rec["result"])
	}
}

func TestNewManager_InvalidGlobalRule_ReturnsError(t *testing.T) {
	cfg := makeCfg([]RuleConfig{
		{Sources: nil, Dest: "out"},
	}, nil)
	_, err := NewManager(cfg)
	if err == nil {
		t.Fatal("expected error for invalid global rule")
	}
}

func TestNewManager_InvalidSourceRule_ReturnsError(t *testing.T) {
	cfg := makeCfg(nil, map[string][]RuleConfig{
		"svc": {{Sources: []string{"a"}, Dest: ""}},
	})
	_, err := NewManager(cfg)
	if err == nil {
		t.Fatal("expected error for invalid source rule")
	}
}
