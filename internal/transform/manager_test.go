package transform

import (
	"testing"

	"github.com/user/logpipe/internal/config"
)

func cfgWithTransforms(transforms []config.TransformConfig) *config.Config {
	return &config.Config{
		Sources: []config.SourceConfig{{Path: "/tmp/a.log"}},
		Sinks:   []config.SinkConfig{{Type: "stdout"}},
		Transforms: transforms,
	}
}

func TestNewManager_NoTransforms(t *testing.T) {
	m, err := NewManager(cfgWithTransforms(nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rec := map[string]string{"level": "info"}
	out := m.Apply("src", rec)
	if out["level"] != "info" {
		t.Errorf("expected passthrough, got %q", out["level"])
	}
}

func TestNewManager_GlobalRule(t *testing.T) {
	cfg := cfgWithTransforms([]config.TransformConfig{
		{Source: "", Rules: []config.TransformRule{{Op: "add", Field: "env", Value: "prod"}}},
	})
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := m.Apply("any-source", map[string]string{"level": "info"})
	if out["env"] != "prod" {
		t.Errorf("expected env=prod from global rule, got %q", out["env"])
	}
}

func TestNewManager_SourceSpecificRule(t *testing.T) {
	cfg := cfgWithTransforms([]config.TransformConfig{
		{Source: "app", Rules: []config.TransformRule{{Op: "add", Field: "svc", Value: "myapp"}}},
	})
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	outApp := m.Apply("app", map[string]string{"level": "warn"})
	if outApp["svc"] != "myapp" {
		t.Errorf("expected svc=myapp for app source")
	}
	outOther := m.Apply("other", map[string]string{"level": "warn"})
	if _, ok := outOther["svc"]; ok {
		t.Error("svc should not appear for unrelated source")
	}
}

func TestNewManager_GlobalAndSourceCombined(t *testing.T) {
	cfg := cfgWithTransforms([]config.TransformConfig{
		{Source: "", Rules: []config.TransformRule{{Op: "add", Field: "env", Value: "prod"}}},
		{Source: "app", Rules: []config.TransformRule{{Op: "rename", Field: "env", NewName: "environment"}}},
	})
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := m.Apply("app", map[string]string{"level": "info"})
	if out["environment"] != "prod" {
		t.Errorf("expected environment=prod after combined transforms, got %q", out["environment"])
	}
	if _, ok := out["env"]; ok {
		t.Error("env should have been renamed to environment")
	}
}

func TestNewManager_InvalidRule(t *testing.T) {
	cfg := cfgWithTransforms([]config.TransformConfig{
		{Source: "", Rules: []config.TransformRule{{Op: "bad-op", Field: "x"}}},
	})
	_, err := NewManager(cfg)
	if err == nil {
		t.Error("expected error for invalid op in global rules")
	}
}
