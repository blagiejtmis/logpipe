package truncate

import (
	"testing"

	"github.com/logpipe/logpipe/internal/config"
)

func makeTruncCfg(defaultRules []config.TruncateRule, sources map[string][]config.TruncateRule) *config.TruncateConfig {
	return &config.TruncateConfig{
		Default: defaultRules,
		Sources: sources,
	}
}

func TestNewManager_NilConfig_AllowsAll(t *testing.T) {
	m, err := NewManager(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.For("any") != nil {
		t.Fatal("expected nil truncator for nil config")
	}
}

func TestNewManager_DefaultRule_Applied(t *testing.T) {
	cfg := makeTruncCfg([]config.TruncateRule{{Field: "msg", MaxLen: 10}}, nil)
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tr := m.For("unknown-source")
	if tr == nil {
		t.Fatal("expected global truncator")
	}
	rec := map[string]any{"msg": "hello world this is long"}
	result := tr.Apply(rec)
	v, _ := result["msg"].(string)
	if len(v) > 10 {
		t.Fatalf("expected msg truncated to 10, got %d chars", len(v))
	}
}

func TestNewManager_SourceSpecific_OverridesDefault(t *testing.T) {
	cfg := makeTruncCfg(
		[]config.TruncateRule{{Field: "msg", MaxLen: 5}},
		map[string][]config.TruncateRule{
			"svc-a": {{Field: "msg", MaxLen: 20}},
		},
	)
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	defaultTr := m.For("other")
	if defaultTr == nil {
		t.Fatal("expected global truncator for unmatched source")
	}

	srcTr := m.For("svc-a")
	if srcTr == nil {
		t.Fatal("expected source-specific truncator")
	}

	rec := map[string]any{"msg": "hello world"}

	defaultResult := defaultTr.Apply(rec)
	if v, _ := defaultResult["msg"].(string); len(v) > 5 {
		t.Fatalf("default: expected max 5, got %d", len(v))
	}

	srcResult := srcTr.Apply(rec)
	if v, _ := srcResult["msg"].(string); v != "hello world" {
		t.Fatalf("source: expected full string, got %q", v)
	}
}

func TestNewManager_InvalidDefaultRule_ReturnsError(t *testing.T) {
	cfg := makeTruncCfg([]config.TruncateRule{{Field: "", MaxLen: 10}}, nil)
	_, err := NewManager(cfg)
	if err == nil {
		t.Fatal("expected error for empty field name")
	}
}

func TestNewManager_InvalidSourceRule_ReturnsError(t *testing.T) {
	cfg := makeTruncCfg(nil, map[string][]config.TruncateRule{
		"svc-b": {{Field: "body", MaxLen: 0}},
	})
	_, err := NewManager(cfg)
	if err == nil {
		t.Fatal("expected error for zero MaxLen")
	}
}
