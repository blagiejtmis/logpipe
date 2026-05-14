package multiline

import (
	"testing"
)

func TestNewManager_NilConfig_AllowsAll(t *testing.T) {
	m, err := NewManager(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	a, err := m.Assembler("any")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a != nil {
		t.Fatal("expected nil assembler for nil config")
	}
}

func TestNewManager_DefaultRule_Applied(t *testing.T) {
	cfg := &Config{
		Default: &RuleConfig{StartPattern: `^\d{4}-`},
	}
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	a, err := m.Assembler("unknown-source")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil assembler from default rule")
	}
}

func TestNewManager_SourceSpecific_OverridesDefault(t *testing.T) {
	cfg := &Config{
		Default: &RuleConfig{StartPattern: `^START`},
		Sources: map[string]*RuleConfig{
			"app": {ContinuePattern: `^\s`},
		},
	}
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	a, err := m.Assembler("app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil assembler")
	}
	// The assembler should use the continue-pattern, not the start-pattern.
	// Feed a line that would trigger start-pattern but NOT continue-pattern.
	a.Add("first")
	out := a.Add("START next") // does not match `^\s`, so flush should occur
	if out == nil {
		t.Fatal("expected flush when non-continuation line arrives")
	}
}

func TestNewManager_InvalidDefaultRule_ReturnsError(t *testing.T) {
	cfg := &Config{
		Default: &RuleConfig{StartPattern: "[bad"},
	}
	_, err := NewManager(cfg)
	if err == nil {
		t.Fatal("expected error for invalid default pattern")
	}
}

func TestNewManager_InvalidSourceRule_ReturnsError(t *testing.T) {
	cfg := &Config{
		Sources: map[string]*RuleConfig{
			"svc": {ContinuePattern: "[bad"},
		},
	}
	_, err := NewManager(cfg)
	if err == nil {
		t.Fatal("expected error for invalid source pattern")
	}
}

func TestNewManager_TimeoutMs_Propagated(t *testing.T) {
	cfg := &Config{
		Default: &RuleConfig{StartPattern: `^LOG`, TimeoutMs: 50},
	}
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	a, _ := m.Assembler("any")
	if a == nil {
		t.Fatal("expected assembler")
	}
	if a.rule.TimeoutMs != 0 {
		// TimeoutMs lives on RuleConfig, not Rule; check Timeout on Rule instead.
	}
	if a.rule.Timeout.Milliseconds() != 50 {
		t.Errorf("expected 50ms timeout, got %v", a.rule.Timeout)
	}
}
