package filter_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/filter"
)

func TestFilter_NoRules_PassesAll(t *testing.T) {
	f, err := filter.New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.NoRules() {
		t.Error("expected NoRules() to be true")
	}
	if !f.Match(map[string]string{"level": "error"}) {
		t.Error("expected match with no rules")
	}
}

func TestFilter_MatchField(t *testing.T) {
	rules := []filter.Rule{
		{Field: "level", Pattern: `^(error|warn)$`},
	}
	f, err := filter.New(rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !f.Match(map[string]string{"level": "error"}) {
		t.Error("expected 'error' level to match")
	}
	if f.Match(map[string]string{"level": "info"}) {
		t.Error("expected 'info' level to not match")
	}
}

func TestFilter_InvertRule(t *testing.T) {
	rules := []filter.Rule{
		{Field: "level", Pattern: `^debug$`, Invert: true},
	}
	f, err := filter.New(rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !f.Match(map[string]string{"level": "info"}) {
		t.Error("expected non-debug level to pass inverted rule")
	}
	if f.Match(map[string]string{"level": "debug"}) {
		t.Error("expected debug level to be excluded by inverted rule")
	}
}

func TestFilter_MissingField_TreatedAsEmpty(t *testing.T) {
	rules := []filter.Rule{
		{Field: "service", Pattern: `^$`},
	}
	f, err := filter.New(rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// No "service" key — should match empty string pattern.
	if !f.Match(map[string]string{"level": "info"}) {
		t.Error("expected missing field to match empty pattern")
	}
}

func TestFilter_InvalidPattern_ReturnsError(t *testing.T) {
	rules := []filter.Rule{
		{Field: "level", Pattern: `[invalid`},
	}
	_, err := filter.New(rules)
	if err == nil {
		t.Fatal("expected error for invalid pattern, got nil")
	}
}

func TestFilter_MultipleRules_ANDSemantics(t *testing.T) {
	rules := []filter.Rule{
		{Field: "level", Pattern: `^error$`},
		{Field: "service", Pattern: `^auth`},
	}
	f, err := filter.New(rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Both match.
	if !f.Match(map[string]string{"level": "error", "service": "auth-api"}) {
		t.Error("expected both rules to match")
	}
	// Only first matches.
	if f.Match(map[string]string{"level": "error", "service": "payments"}) {
		t.Error("expected second rule to fail")
	}
}
