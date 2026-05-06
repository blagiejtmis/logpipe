package router

import (
	"testing"

	"github.com/user/logpipe/internal/config"
)

func makeConfig(rules []config.RoutingRule, defaults []string, sinkNames []string) *config.Config {
	cfg := &config.Config{}
	cfg.Routing.Rules = rules
	cfg.Routing.DefaultSinks = defaults
	for _, n := range sinkNames {
		cfg.Sinks = append(cfg.Sinks, config.SinkConfig{Name: n})
	}
	return cfg
}

func TestNewFromConfig_DefaultsFromSinks(t *testing.T) {
	cfg := makeConfig(nil, nil, []string{"stdout", "file1"})
	r, err := NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sinks := r.Resolve("any", nil)
	if len(sinks) != 2 {
		t.Fatalf("expected 2 default sinks, got %v", sinks)
	}
}

func TestNewFromConfig_ExplicitDefaults(t *testing.T) {
	cfg := makeConfig(nil, []string{"stdout"}, []string{"stdout", "file1"})
	r, err := NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sinks := r.Resolve("any", nil)
	if len(sinks) != 1 || sinks[0] != "stdout" {
		t.Fatalf("expected [stdout], got %v", sinks)
	}
}

func TestNewFromConfig_RuleWithNoSinks_ReturnsError(t *testing.T) {
	cfg := makeConfig([]config.RoutingRule{{Source: "app"}}, nil, nil)
	_, err := NewFromConfig(cfg)
	if err == nil {
		t.Fatal("expected error for rule with no sinks")
	}
}

func TestNewFromConfig_RoutingRule(t *testing.T) {
	rules := []config.RoutingRule{
		{Source: "svc\\.log", Sinks: []string{"file-svc"}},
	}
	cfg := makeConfig(rules, []string{"stdout"}, nil)
	r, err := NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sinks := r.Resolve("svc.log", nil); len(sinks) != 1 || sinks[0] != "file-svc" {
		t.Fatalf("expected [file-svc], got %v", sinks)
	}
}
