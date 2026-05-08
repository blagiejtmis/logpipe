package enrich

import (
	"testing"

	"github.com/yourorg/logpipe/internal/config"
)

func makeCfg(global map[string]string, sources map[string]map[string]string) *config.EnrichConfig {
	return &config.EnrichConfig{Global: global, Sources: sources}
}

func TestNewManager_NilConfig_AllowsAll(t *testing.T) {
	m, err := NewManager(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rec := map[string]any{"msg": "hello"}
	m.Apply("src1", rec)
	if len(rec) != 1 {
		t.Errorf("expected record unchanged, got %v", rec)
	}
}

func TestNewManager_GlobalRule(t *testing.T) {
	cfg := makeCfg(map[string]string{"env": "prod"}, nil)
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rec := map[string]any{"msg": "hello"}
	m.Apply("any-source", rec)
	if rec["env"] != "prod" {
		t.Errorf("expected env=prod, got %v", rec["env"])
	}
}

func TestNewManager_SourceSpecific_OverridesGlobal(t *testing.T) {
	cfg := makeCfg(
		map[string]string{"env": "prod"},
		map[string]map[string]string{
			"app-log": {"env": "staging", "service": "app"},
		},
	)
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rec := map[string]any{"msg": "hello"}
	m.Apply("app-log", rec)
	if rec["env"] != "staging" {
		t.Errorf("expected env=staging after source override, got %v", rec["env"])
	}
	if rec["service"] != "app" {
		t.Errorf("expected service=app, got %v", rec["service"])
	}
}

func TestNewManager_GlobalOnly_NoSourceMatch(t *testing.T) {
	cfg := makeCfg(map[string]string{"region": "us-east-1"}, nil)
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rec := map[string]any{}
	m.Apply("unknown-source", rec)
	if rec["region"] != "us-east-1" {
		t.Errorf("expected region from global, got %v", rec["region"])
	}
}

func TestNewManager_InvalidRule_ReturnsError(t *testing.T) {
	cfg := makeCfg(map[string]string{"": "value"}, nil)
	_, err := NewManager(cfg)
	if err == nil {
		t.Fatal("expected error for empty field name, got nil")
	}
}
