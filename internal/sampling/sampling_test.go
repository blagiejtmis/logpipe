package sampling

import (
	"testing"
)

func TestNew_InvalidRate_Zero(t *testing.T) {
	_, err := New([]Config{{Rate: 0.0}})
	if err == nil {
		t.Fatal("expected error for rate 0.0")
	}
}

func TestNew_InvalidRate_AboveOne(t *testing.T) {
	_, err := New([]Config{{Rate: 1.1}})
	if err == nil {
		t.Fatal("expected error for rate 1.1")
	}
}

func TestNew_ValidRate(t *testing.T) {
	_, err := New([]Config{{Rate: 0.5}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAllow_RateOne_AlwaysAllows(t *testing.T) {
	s, err := New([]Config{{Rate: 1.0}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 100; i++ {
		if !s.Allow("any") {
			t.Fatal("rate=1.0 should always allow")
		}
	}
}

func TestAllow_NoRules_DefaultsToOne(t *testing.T) {
	s, err := New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 50; i++ {
		if !s.Allow("src") {
			t.Fatal("default rate=1.0 should always allow")
		}
	}
}

func TestAllow_LowRate_DropsRecords(t *testing.T) {
	// Rate of 0.01 should drop the vast majority over 10000 trials.
	s, err := New([]Config{{Rate: 0.01}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	allowed := 0
	const trials = 10_000
	for i := 0; i < trials; i++ {
		if s.Allow("src") {
			allowed++
		}
	}
	// Expect roughly 1% ± generous margin.
	if allowed > 300 {
		t.Errorf("expected ~100 allowed, got %d", allowed)
	}
}

func TestAllow_SourceSpecific_OverridesDefault(t *testing.T) {
	s, err := New([]Config{
		{Rate: 1.0},
		{Rate: 0.01, Source: "noisy"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// "other" should always pass (default rate 1.0)
	for i := 0; i < 50; i++ {
		if !s.Allow("other") {
			t.Fatal("default source should always be allowed")
		}
	}
	// "noisy" should rarely pass
	allowed := 0
	for i := 0; i < 5000; i++ {
		if s.Allow("noisy") {
			allowed++
		}
	}
	if allowed > 150 {
		t.Errorf("noisy source: expected ~50 allowed, got %d", allowed)
	}
}
