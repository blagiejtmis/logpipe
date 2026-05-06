package ratelimit

import (
	"testing"

	"github.com/logpipe/logpipe/internal/config"
)

func makeRLConfig(defaultRate int, sources []config.SourceRateLimit) *config.Config {
	return &config.Config{
		RateLimit: config.RateLimitConfig{
			Default: config.RateLimitEntry{Rate: defaultRate, WindowSecs: 1},
			Sources: sources,
		},
	}
}

func TestNewManager_NoLimits_AllowsAll(t *testing.T) {
	cfg := makeRLConfig(0, nil)
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 100; i++ {
		if !m.Allow("/var/log/app.log") {
			t.Fatal("expected all lines to be allowed when no limits configured")
		}
	}
}

func TestNewManager_DefaultLimit_Applied(t *testing.T) {
	cfg := makeRLConfig(5, nil)
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	allowed := 0
	for i := 0; i < 20; i++ {
		if m.Allow("any-source") {
			allowed++
		}
	}
	if allowed > 5 {
		t.Fatalf("expected at most 5 allowed, got %d", allowed)
	}
}

func TestNewManager_SourceSpecific_OverridesDefault(t *testing.T) {
	cfg := makeRLConfig(100, []config.SourceRateLimit{
		{Source: "/var/log/noisy.log", Rate: 2, WindowSecs: 1},
	})
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	allowed := 0
	for i := 0; i < 20; i++ {
		if m.Allow("/var/log/noisy.log") {
			allowed++
		}
	}
	if allowed > 2 {
		t.Fatalf("source-specific limit should cap at 2, got %d", allowed)
	}
}

func TestNewManager_InvalidSourceRate_ReturnsError(t *testing.T) {
	cfg := makeRLConfig(0, []config.SourceRateLimit{
		{Source: "/var/log/app.log", Rate: -1, WindowSecs: 1},
	})
	_, err := NewManager(cfg)
	if err == nil {
		t.Fatal("expected error for negative rate")
	}
}

func TestNewManager_InvalidDefaultRate_ReturnsError(t *testing.T) {
	cfg := makeRLConfig(-5, nil)
	_, err := NewManager(cfg)
	if err == nil {
		t.Fatal("expected error for negative default rate")
	}
}
