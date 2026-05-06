package dedupe

import (
	"testing"

	"github.com/your-org/logpipe/internal/config"
)

func makeCfg(globalFields []string, globalWindow int, srcRules map[string]config.DedupeRule) config.DedupeConfig {
	cfg := config.DedupeConfig{Sources: srcRules}
	if globalFields != nil {
		cfg.Global = &config.DedupeRule{Fields: globalFields, WindowSecs: globalWindow}
	}
	return cfg
}

func TestNewManager_NoConfig_AllowsAll(t *testing.T) {
	m, err := NewManager(config.DedupeConfig{Sources: map[string]config.DedupeRule{}})
	if err != nil {
		t.Fatal(err)
	}
	if !m.Allow("src", Record{"msg": "x"}) {
		t.Fatal("expected allow when no rules configured")
	}
}

func TestNewManager_GlobalRule(t *testing.T) {
	m, err := NewManager(makeCfg([]string{"msg"}, 60, nil))
	if err != nil {
		t.Fatal(err)
	}
	rec := Record{"msg": "hello"}
	if !m.Allow("any", rec) {
		t.Fatal("first occurrence should be allowed")
	}
	if m.Allow("any", rec) {
		t.Fatal("duplicate should be rejected by global rule")
	}
}

func TestNewManager_SourceSpecific_OverridesGlobal(t *testing.T) {
	srcRules := map[string]config.DedupeRule{
		"app": {Fields: []string{"level"}, WindowSecs: 60},
	}
	m, err := NewManager(makeCfg([]string{"msg"}, 60, srcRules))
	if err != nil {
		t.Fatal(err)
	}
	rec := Record{"msg": "same", "level": "error"}
	// source "app" uses "level" field — same level should be deduped
	m.Allow("app", rec)
	if m.Allow("app", rec) {
		t.Fatal("source-specific rule should deduplicate")
	}
	// source "other" falls back to global rule keyed on "msg"
	if !m.Allow("other", rec) {
		t.Fatal("first occurrence on other source should be allowed")
	}
}

func TestNewManager_InvalidRule_ReturnsError(t *testing.T) {
	srcRules := map[string]config.DedupeRule{
		"bad": {Fields: []string{}, WindowSecs: 60},
	}
	_, err := NewManager(makeCfg(nil, 0, srcRules))
	if err == nil {
		t.Fatal("expected error for empty fields")
	}
}
