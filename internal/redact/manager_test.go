package redact

import (
	"testing"

	"github.com/logpipe/logpipe/internal/config"
)

func makeCfg(global []config.RedactRule, sources map[string][]config.RedactRule) *config.RedactConfig {
	return &config.RedactConfig{
		Global:  global,
		Sources: sources,
	}
}

func TestNewManager_NilConfig_AllowsAll(t *testing.T) {
	m, err := NewManager(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rec := map[string]any{"password": "secret"}
	out := m.Apply("any", rec)
	if out["password"] != "secret" {
		t.Errorf("expected passthrough, got %v", out["password"])
	}
}

func TestNewManager_GlobalRule(t *testing.T) {
	cfg := makeCfg([]config.RedactRule{
		{Field: "token", Pattern: ".*", Placeholder: "[REDACTED]"},
	}, nil)
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rec := map[string]any{"token": "abc123"}
	out := m.Apply("src-a", rec)
	if out["token"] != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %v", out["token"])
	}
}

func TestNewManager_SourceSpecific_OverridesGlobal(t *testing.T) {
	cfg := makeCfg(
		[]config.RedactRule{{Field: "token", Pattern: ".*", Placeholder: "[GLOBAL]"}},
		map[string][]config.RedactRule{
			"src-b": {{Field: "token", Pattern: ".*", Placeholder: "[SRC]"}},
		},
	)
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	recA := map[string]any{"token": "x"}
	if out := m.Apply("src-a", recA); out["token"] != "[GLOBAL]" {
		t.Errorf("src-a: expected [GLOBAL], got %v", out["token"])
	}

	recB := map[string]any{"token": "x"}
	if out := m.Apply("src-b", recB); out["token"] != "[SRC]" {
		t.Errorf("src-b: expected [SRC], got %v", out["token"])
	}
}

func TestNewManager_InvalidGlobalRule_ReturnsError(t *testing.T) {
	cfg := makeCfg([]config.RedactRule{
		{Field: "f", Pattern: "[", Placeholder: "x"},
	}, nil)
	_, err := NewManager(cfg)
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestNewManager_InvalidSourceRule_ReturnsError(t *testing.T) {
	cfg := makeCfg(nil, map[string][]config.RedactRule{
		"src-c": {{Field: "f", Pattern: "[", Placeholder: "x"}},
	})
	_, err := NewManager(cfg)
	if err == nil {
		t.Fatal("expected error for invalid source pattern")
	}
}
