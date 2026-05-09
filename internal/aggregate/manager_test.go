package aggregate

import (
	"testing"
)

func makeCfg(defaultRules []RuleConfig, sources map[string][]RuleConfig) *Config {
	return &Config{Default: defaultRules, Sources: sources}
}

func TestNewManager_NilConfig_AllowsAll(t *testing.T) {
	m, err := NewManager(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should not panic on Add or Snapshot.
	m.Add("src", map[string]any{"bytes": float64(1)})
	snap := m.Snapshot()
	if len(snap) != 0 {
		t.Fatalf("expected empty snapshot, got %v", snap)
	}
}

func TestNewManager_DefaultRule_Applied(t *testing.T) {
	cfg := makeCfg([]RuleConfig{
		{Field: "n", Op: "count", Window: "1m"},
	}, nil)
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m.Add("src1", map[string]any{"n": 1})
	m.Add("src1", map[string]any{"n": 1})
	snap := m.Snapshot()
	if snap["default"]["src1"]["n"] != 2 {
		t.Fatalf("expected count 2, got %v", snap)
	}
}

func TestNewManager_SourceSpecific_OverridesDefault(t *testing.T) {
	cfg := makeCfg(
		[]RuleConfig{{Field: "n", Op: "count", Window: "1m"}},
		map[string][]RuleConfig{
			"app": {{Field: "bytes", Op: "sum", Window: "1m"}},
		},
	)
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m.Add("app", map[string]any{"bytes": float64(512)})
	m.Add("other", map[string]any{"n": 1})
	snap := m.Snapshot()
	if snap["app"]["app"]["bytes"] != 512 {
		t.Fatalf("expected bytes sum 512, got %v", snap["app"])
	}
	if snap["default"]["other"]["n"] != 1 {
		t.Fatalf("expected default count 1, got %v", snap["default"])
	}
}

func TestNewManager_InvalidDefaultRule_ReturnsError(t *testing.T) {
	cfg := makeCfg([]RuleConfig{
		{Field: "n", Op: "count", Window: "bad"},
	}, nil)
	_, err := NewManager(cfg)
	if err == nil {
		t.Fatal("expected error for invalid window")
	}
}

func TestNewManager_InvalidSourceRule_ReturnsError(t *testing.T) {
	cfg := makeCfg(nil, map[string][]RuleConfig{
		"svc": {{Field: "lat", Op: "unknown", Window: "1m"}},
	})
	_, err := NewManager(cfg)
	if err == nil {
		t.Fatal("expected error for unknown op")
	}
}
