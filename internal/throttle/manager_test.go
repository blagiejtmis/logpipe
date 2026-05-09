package throttle

import (
	"testing"
)

func makeThrottleCfg(def int, sources map[string]int) *Config {
	return &Config{
		DefaultRatePerSec: def,
		Sources:           sources,
	}
}

func TestNewManager_NilConfig_AllowsAll(t *testing.T) {
	m, err := NewManager(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	w, err := m.For("any-source")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// noop — should not panic
	w.Wait()
}

func TestNewManager_DefaultRate_Applied(t *testing.T) {
	cfg := makeThrottleCfg(100, nil)
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	w, err := m.For("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := w.(*Throttler); !ok {
		t.Errorf("expected *Throttler, got %T", w)
	}
}

func TestNewManager_SourceSpecific_OverridesDefault(t *testing.T) {
	cfg := makeThrottleCfg(10, map[string]int{"fast": 500})
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	w, err := m.For("fast")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	th, ok := w.(*Throttler)
	if !ok {
		t.Fatalf("expected *Throttler, got %T", w)
	}
	if th.rate != 500 {
		t.Errorf("expected rate 500, got %d", th.rate)
	}
}

func TestNewManager_ZeroRate_ReturnsNoop(t *testing.T) {
	cfg := makeThrottleCfg(0, nil)
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	w, err := m.For("src")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := w.(*noopThrottler); !ok {
		t.Errorf("expected *noopThrottler for zero rate, got %T", w)
	}
}

func TestNewManager_CachesThrottler(t *testing.T) {
	cfg := makeThrottleCfg(50, nil)
	m, _ := NewManager(cfg)

	w1, _ := m.For("src")
	w2, _ := m.For("src")
	if w1 != w2 {
		t.Error("expected same Throttler instance on repeated calls")
	}
}
