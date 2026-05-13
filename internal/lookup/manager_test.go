package lookup

import (
	"testing"

	"github.com/yourorg/logpipe/internal/config"
)

func makeLookupCfg(global []config.LookupRule, sources map[string][]config.LookupRule) *config.LookupConfig {
	return &config.LookupConfig{Global: global, Sources: sources}
}

func TestNewManager_NilConfig_AllowsAll(t *testing.T) {
	m, err := NewManager(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rec := map[string]any{"msg": "hello"}
	out, err := m.Apply("src", rec)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["msg"] != "hello" {
		t.Errorf("expected record unchanged")
	}
}

func TestNewManager_GlobalRule(t *testing.T) {
	table := map[string]map[string]any{
		"ERR": {"label": "error"},
	}
	cfg := makeLookupCfg([]config.LookupRule{
		{KeyField: "level", Table: table, DestField: "label", OnMiss: "keep"},
	}, nil)
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rec := map[string]any{"level": "ERR"}
	out, err := m.Apply("anysource", rec)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["label"] != "error" {
		t.Errorf("expected label=error, got %v", out["label"])
	}
}

func TestNewManager_SourceSpecific_OverridesGlobal(t *testing.T) {
	globalTable := map[string]map[string]any{"x": {"region": "global"}}
	srcTable := map[string]map[string]any{"x": {"region": "src-specific"}}

	cfg := makeLookupCfg(
		[]config.LookupRule{{KeyField: "k", Table: globalTable, DestField: "region", OnMiss: "keep"}},
		map[string][]config.LookupRule{
			"myapp": {{KeyField: "k", Table: srcTable, DestField: "region", OnMiss: "keep"}},
		},
	)
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	rec := map[string]any{"k": "x"}
	out, err := m.Apply("myapp", rec)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["region"] != "src-specific" {
		t.Errorf("expected src-specific, got %v", out["region"])
	}

	rec2 := map[string]any{"k": "x"}
	out2, err := m.Apply("other", rec2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out2["region"] != "global" {
		t.Errorf("expected global, got %v", out2["region"])
	}
}

func TestNewManager_InvalidGlobalRule_ReturnsError(t *testing.T) {
	cfg := makeLookupCfg([]config.LookupRule{
		{KeyField: "", Table: map[string]map[string]any{}, DestField: "x", OnMiss: "keep"},
	}, nil)
	_, err := NewManager(cfg)
	if err == nil {
		t.Fatal("expected error for empty KeyField")
	}
}

func TestNewManager_InvalidSourceRule_ReturnsError(t *testing.T) {
	cfg := makeLookupCfg(nil, map[string][]config.LookupRule{
		"bad": {{KeyField: "k", Table: map[string]map[string]any{}, DestField: "", OnMiss: "keep"}},
	})
	_, err := NewManager(cfg)
	if err == nil {
		t.Fatal("expected error for empty DestField")
	}
}
