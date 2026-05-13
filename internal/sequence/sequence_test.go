package sequence

import (
	"testing"

	"logpipe/internal/config"
)

func base() map[string]any {
	return map[string]any{"message": "hello", "level": "info"}
}

func TestNew_EmptyField_ReturnsError(t *testing.T) {
	_, err := New("", 1)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_ValidArgs_NoError(t *testing.T) {
	s, err := New("seq", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil sequencer")
	}
}

func TestApply_IncrementsSequence(t *testing.T) {
	s, _ := New("seq", 1)
	r1 := s.Apply(base())
	r2 := s.Apply(base())

	if r1["seq"] != uint64(1) {
		t.Errorf("expected seq=1, got %v", r1["seq"])
	}
	if r2["seq"] != uint64(2) {
		t.Errorf("expected seq=2, got %v", r2["seq"])
	}
}

func TestApply_CustomStart(t *testing.T) {
	s, _ := New("n", 100)
	r := s.Apply(base())
	if r["n"] != uint64(100) {
		t.Errorf("expected n=100, got %v", r["n"])
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	s, _ := New("seq", 1)
	in := base()
	out := s.Apply(in)
	if _, ok := in["seq"]; ok {
		t.Error("Apply mutated the input record")
	}
	if _, ok := out["seq"]; !ok {
		t.Error("output record missing seq field")
	}
}

func TestReset_ResetsCounter(t *testing.T) {
	s, _ := New("seq", 1)
	s.Apply(base())
	s.Apply(base())
	s.Reset()
	r := s.Apply(base())
	if r["seq"] != uint64(1) {
		t.Errorf("expected seq=1 after reset, got %v", r["seq"])
	}
}

func TestNewManager_NilConfig_AllowsAll(t *testing.T) {
	m, err := NewManager(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rec := m.Apply("src", base())
	if _, ok := rec["seq"]; ok {
		t.Error("expected no seq field when manager has no config")
	}
}

func TestNewManager_DefaultApplied(t *testing.T) {
	cfg := &config.SequenceConfig{
		Default: &Config{Field: "seq", Start: 1},
	}
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := m.Apply("any-source", base())
	if r["seq"] != uint64(1) {
		t.Errorf("expected seq=1, got %v", r["seq"])
	}
}

func TestNewManager_SourceSpecific_OverridesDefault(t *testing.T) {
	cfg := &config.SequenceConfig{
		Default: &Config{Field: "seq", Start: 1},
		Sources: map[string]*Config{
			"app": {Field: "app_seq", Start: 500},
		},
	}
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	generic := m.Apply("other", base())
	if _, ok := generic["seq"]; !ok {
		t.Error("expected global seq field on non-specific source")
	}

	specific := m.Apply("app", base())
	if specific["app_seq"] != uint64(500) {
		t.Errorf("expected app_seq=500, got %v", specific["app_seq"])
	}
	if _, ok := specific["seq"]; ok {
		t.Error("source-specific sequencer should not add global seq field")
	}
}

func TestNewManager_InvalidSourceRule_ReturnsError(t *testing.T) {
	cfg := &config.SequenceConfig{
		Sources: map[string]*Config{
			"bad": {Field: "", Start: 1},
		},
	}
	_, err := NewManager(cfg)
	if err == nil {
		t.Fatal("expected error for empty field in source rule")
	}
}
