package schema

import (
	"testing"

	"github.com/logpipe/logpipe/internal/config"
)

func makeSchemaCfg(global []config.FieldRule, sources map[string][]config.FieldRule) *config.SchemaConfig {
	return &config.SchemaConfig{
		Global:  global,
		Sources: sources,
	}
}

func TestNewManager_NilConfig_AllowsAll(t *testing.T) {
	m, err := NewManager(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rec := map[string]any{"msg": "hello"}
	if err := m.Validate("any-source", rec); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestNewManager_GlobalRule_ValidRecord(t *testing.T) {
	cfg := makeSchemaCfg([]config.FieldRule{
		{Field: "level", Type: "string", Required: true},
	}, nil)
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rec := map[string]any{"level": "info"}
	if err := m.Validate("src", rec); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestNewManager_GlobalRule_InvalidRecord(t *testing.T) {
	cfg := makeSchemaCfg([]config.FieldRule{
		{Field: "level", Type: "string", Required: true},
	}, nil)
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := m.Validate("src", map[string]any{}); err == nil {
		t.Fatal("expected error for missing required field")
	}
}

func TestNewManager_SourceSpecific_OverridesGlobal(t *testing.T) {
	cfg := makeSchemaCfg(
		[]config.FieldRule{{Field: "level", Type: "string", Required: true}},
		map[string][]config.FieldRule{
			"app": {{Field: "code", Type: "float", Required: true}},
		},
	)
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// source "app" uses its own rule — "level" not required
	rec := map[string]any{"code": float64(200)}
	if err := m.Validate("app", rec); err != nil {
		t.Fatalf("expected nil for source-specific schema, got %v", err)
	}
	// source "other" falls back to global — "level" required
	if err := m.Validate("other", rec); err == nil {
		t.Fatal("expected error from global schema for missing level")
	}
}

func TestNewManager_InvalidGlobalRule_ReturnsError(t *testing.T) {
	cfg := makeSchemaCfg([]config.FieldRule{
		{Field: "x", Type: "badtype", Required: true},
	}, nil)
	if _, err := NewManager(cfg); err == nil {
		t.Fatal("expected error for invalid type")
	}
}

func TestNewManager_InvalidSourceRule_ReturnsError(t *testing.T) {
	cfg := makeSchemaCfg(nil, map[string][]config.FieldRule{
		"svc": {{Field: "ts", Type: "unknown", Required: false}},
	})
	if _, err := NewManager(cfg); err == nil {
		t.Fatal("expected error for invalid source rule type")
	}
}
