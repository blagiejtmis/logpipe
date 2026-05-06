package buffer

import (
	"testing"

	"github.com/logpipe/logpipe/internal/config"
)

func makeBufCfg(cap int, policy string, sources map[string]config.BufferSourceConfig) config.BufferConfig {
	return config.BufferConfig{
		Capacity: cap,
		Policy:   policy,
		Sources:  sources,
	}
}

func TestNewManager_NoConfig_ReturnsNilBuffer(t *testing.T) {
	m, err := NewManager(makeBufCfg(0, "", nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.BufferFor("any") != nil {
		t.Fatal("expected nil buffer when no config provided")
	}
}

func TestNewManager_DefaultBuffer_Applied(t *testing.T) {
	m, err := NewManager(makeBufCfg(32, DropNew, nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	buf := m.BufferFor("unknown-source")
	if buf == nil {
		t.Fatal("expected default buffer, got nil")
	}
	if buf.Cap() != 32 {
		t.Fatalf("expected cap 32, got %d", buf.Cap())
	}
}

func TestNewManager_SourceSpecific_OverridesDefault(t *testing.T) {
	sources := map[string]config.BufferSourceConfig{
		"/var/log/app.log": {Capacity: 8, Policy: DropOld},
	}
	m, err := NewManager(makeBufCfg(64, DropNew, sources))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	src := m.BufferFor("/var/log/app.log")
	if src == nil {
		t.Fatal("expected source-specific buffer")
	}
	if src.Cap() != 8 {
		t.Fatalf("expected cap 8, got %d", src.Cap())
	}

	def := m.BufferFor("other")
	if def == nil || def.Cap() != 64 {
		t.Fatal("expected default buffer with cap 64 for unknown source")
	}
}

func TestNewManager_InvalidDefaultPolicy_ReturnsError(t *testing.T) {
	_, err := NewManager(makeBufCfg(16, "bad-policy", nil))
	if err == nil {
		t.Fatal("expected error for invalid policy")
	}
}

func TestNewManager_InvalidSourcePolicy_ReturnsError(t *testing.T) {
	sources := map[string]config.BufferSourceConfig{
		"src": {Capacity: 4, Policy: "invalid"},
	}
	_, err := NewManager(makeBufCfg(0, "", sources))
	if err == nil {
		t.Fatal("expected error for invalid source policy")
	}
}
