package jitter

import (
	"testing"
	"time"
)

func makeJitterCfg(defMin, defMax string, sources map[string][2]string) *Config {
	cfg := &Config{Sources: make(map[string]*Rule)}
	if defMin != "" || defMax != "" {
		cfg.Default = &Rule{Min: defMin, Max: defMax}
	}
	for src, pair := range sources {
		cfg.Sources[src] = &Rule{Min: pair[0], Max: pair[1]}
	}
	return cfg
}

func TestNewManager_NilConfig_AllowsAll(t *testing.T) {
	m, err := NewManager(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	j := m.For("any-source")
	min, max := j.Range()
	if min != 0 || max != 0 {
		t.Fatalf("expected zero jitter, got min=%s max=%s", min, max)
	}
}

func TestNewManager_DefaultRule_Applied(t *testing.T) {
	cfg := makeJitterCfg("1ms", "5ms", nil)
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	j := m.For("unknown-source")
	min, max := j.Range()
	if min != time.Millisecond || max != 5*time.Millisecond {
		t.Fatalf("expected 1ms–5ms, got %s–%s", min, max)
	}
}

func TestNewManager_SourceSpecific_OverridesDefault(t *testing.T) {
	cfg := makeJitterCfg("1ms", "5ms", map[string][2]string{
		"fast-source": {"0", "1ms"},
	})
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	j := m.For("fast-source")
	_, max := j.Range()
	if max != time.Millisecond {
		t.Fatalf("expected max=1ms for fast-source, got %s", max)
	}
	j2 := m.For("other")
	_, max2 := j2.Range()
	if max2 != 5*time.Millisecond {
		t.Fatalf("expected max=5ms for other, got %s", max2)
	}
}

func TestNewManager_InvalidDefaultMin_ReturnsError(t *testing.T) {
	cfg := makeJitterCfg("bad", "5ms", nil)
	_, err := NewManager(cfg)
	if err == nil {
		t.Fatal("expected error for invalid default min")
	}
}

func TestNewManager_InvalidSourceMax_ReturnsError(t *testing.T) {
	cfg := makeJitterCfg("", "", map[string][2]string{
		"src": {"1ms", "notaduration"},
	})
	_, err := NewManager(cfg)
	if err == nil {
		t.Fatal("expected error for invalid source max")
	}
}
