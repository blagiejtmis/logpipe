package unwrap_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/unwrap"
)

func makeCfg(global []unwrap.Rule, sources map[string][]unwrap.Rule) *unwrap.Config {
	return &unwrap.Config{Global: global, Sources: sources}
}

func TestNewManager_NilConfig_AllowsAll(t *testing.T) {
	m, err := unwrap.NewManager(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.For("any") != nil {
		t.Error("expected nil unwrapper for nil config")
	}
}

func TestNewManager_GlobalRule(t *testing.T) {
	cfg := makeCfg([]unwrap.Rule{{Field: "meta", Remove: true}}, nil)
	m, err := unwrap.NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.For("any") == nil {
		t.Error("expected global unwrapper")
	}
}

func TestNewManager_SourceSpecific_OverridesGlobal(t *testing.T) {
	cfg := makeCfg(
		[]unwrap.Rule{{Field: "meta"}},
		map[string][]unwrap.Rule{
			"svc": {{Field: "details", Prefix: "d_"}},
		},
	)
	m, err := unwrap.NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	svc := m.For("svc")
	if svc == nil {
		t.Fatal("expected source-specific unwrapper")
	}
	// global should still be returned for unknown sources
	if m.For("other") == nil {
		t.Error("expected global unwrapper for unknown source")
	}
}

func TestNewManager_InvalidGlobalRule_ReturnsError(t *testing.T) {
	cfg := makeCfg([]unwrap.Rule{{Field: ""}}, nil)
	_, err := unwrap.NewManager(cfg)
	if err == nil {
		t.Fatal("expected error for invalid global rule")
	}
}

func TestNewManager_InvalidSourceRule_ReturnsError(t *testing.T) {
	cfg := makeCfg(nil, map[string][]unwrap.Rule{
		"bad": {{Field: ""}},
	})
	_, err := unwrap.NewManager(cfg)
	if err == nil {
		t.Fatal("expected error for invalid source rule")
	}
}
