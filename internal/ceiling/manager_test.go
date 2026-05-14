package ceiling

import (
	"testing"
	"time"
)

func makeCeilingCfg(defaultMax int64, defaultWindow time.Duration, sources map[string]*Rule) *Config {
	cfg := &Config{
		Sources: sources,
	}
	if defaultMax > 0 {
		cfg.Default = &Rule{Max: defaultMax, Window: defaultWindow}
	}
	return cfg
}

func TestNewManager_NilConfig_AllowsAll(t *testing.T) {
	m, err := NewManager(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 100; i++ {
		if !m.Allow("any-source") {
			t.Fatal("expected allow for nil config")
		}
	}
}

func TestNewManager_DefaultCeiling_Applied(t *testing.T) {
	cfg := makeCeilingCfg(3, time.Minute, nil)
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 3; i++ {
		if !m.Allow("src") {
			t.Fatalf("expected allow on call %d", i+1)
		}
	}
	if m.Allow("src") {
		t.Fatal("expected deny after ceiling reached")
	}
}

func TestNewManager_SourceSpecific_OverridesDefault(t *testing.T) {
	cfg := makeCeilingCfg(10, time.Minute, map[string]*Rule{
		"special": {Max: 2, Window: time.Minute},
	})
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// "special" source hits ceiling at 2
	if !m.Allow("special") || !m.Allow("special") {
		t.Fatal("expected first two allows for special")
	}
	if m.Allow("special") {
		t.Fatal("expected deny after source ceiling")
	}
	// default source still has 10
	for i := 0; i < 10; i++ {
		if !m.Allow("other") {
			t.Fatalf("expected allow for other on call %d", i+1)
		}
	}
}

func TestNewManager_InvalidDefaultRule_ReturnsError(t *testing.T) {
	cfg := &Config{
		Default: &Rule{Max: 0, Window: time.Minute},
	}
	_, err := NewManager(cfg)
	if err == nil {
		t.Fatal("expected error for zero max")
	}
}

func TestNewManager_InvalidSourceRule_ReturnsError(t *testing.T) {
	cfg := &Config{
		Sources: map[string]*Rule{
			"bad": {Max: 5, Window: 0},
		},
	}
	_, err := NewManager(cfg)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}
