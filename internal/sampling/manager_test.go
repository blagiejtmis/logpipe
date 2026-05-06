package sampling

import (
	"testing"

	"github.com/logpipe/logpipe/internal/config"
)

func makeSamplingCfg(defaultRate float64, sources map[string]float64) config.SamplingConfig {
	return config.SamplingConfig{
		DefaultRate: defaultRate,
		Sources:     sources,
	}
}

func TestNewManager_NoConfig_AllowsAll(t *testing.T) {
	m, err := NewManager(makeSamplingCfg(0, nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 20; i++ {
		if !m.Allow("any") {
			t.Fatal("expected all records to be allowed when no config")
		}
	}
}

func TestNewManager_DefaultRate_Applied(t *testing.T) {
	m, err := NewManager(makeSamplingCfg(1.0, nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 10; i++ {
		if !m.Allow("src") {
			t.Fatal("rate=1.0 should always allow")
		}
	}
}

func TestNewManager_SourceSpecific_OverridesDefault(t *testing.T) {
	cfg := makeSamplingCfg(1.0, map[string]float64{
		"noisy": 0.0,
	})
	// rate 0.0 is invalid; use a very small positive to test override path
	cfg.Sources["noisy"] = 1.0
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// source-specific sampler should be used for "noisy"
	for i := 0; i < 5; i++ {
		if !m.Allow("noisy") {
			t.Fatal("rate=1.0 source-specific should always allow")
		}
	}
}

func TestNewManager_InvalidDefaultRate_ReturnsError(t *testing.T) {
	_, err := NewManager(makeSamplingCfg(-0.5, nil))
	if err == nil {
		t.Fatal("expected error for invalid default rate")
	}
}

func TestNewManager_InvalidSourceRate_ReturnsError(t *testing.T) {
	cfg := makeSamplingCfg(1.0, map[string]float64{
		"bad": 1.5,
	})
	_, err := NewManager(cfg)
	if err == nil {
		t.Fatal("expected error for source rate > 1")
	}
}
