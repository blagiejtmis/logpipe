package masking

import (
	"testing"
)

func base() map[string]any {
	return map[string]any{
		"user":     "alice",
		"email":    "alice@example.com",
		"password": "s3cr3t",
		"level":    "info",
	}
}

func TestNew_EmptyField_ReturnsError(t *testing.T) {
	_, err := New([]Rule{{Field: ""}})
	if err == nil {
		t.Fatal("expected error for empty field name")
	}
}

func TestNew_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := New([]Rule{{Field: "email", Pattern: "(unclosed"}})
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestApply_NoRules_PassesThrough(t *testing.T) {
	m, err := New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := m.Apply(base())
	if out["password"] != "s3cr3t" {
		t.Errorf("expected value unchanged, got %v", out["password"])
	}
}

func TestApply_FullMask_DefaultPlaceholder(t *testing.T) {
	m, err := New([]Rule{{Field: "password"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := m.Apply(base())
	if out["password"] != "***" {
		t.Errorf("expected ***, got %v", out["password"])
	}
	if out["user"] != "alice" {
		t.Error("unrelated field should be unchanged")
	}
}

func TestApply_FullMask_CustomPlaceholder(t *testing.T) {
	m, err := New([]Rule{{Field: "password", Placeholder: "[REDACTED]"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := m.Apply(base())
	if out["password"] != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %v", out["password"])
	}
}

func TestApply_PatternMask_PartialReplacement(t *testing.T) {
	m, err := New([]Rule{{Field: "email", Pattern: `[^@]+$`}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := m.Apply(base())
	if out["email"] != "alice@***" {
		t.Errorf("unexpected email mask: %v", out["email"])
	}
}

func TestApply_MissingField_NoError(t *testing.T) {
	m, err := New([]Rule{{Field: "token"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := m.Apply(base())
	if _, ok := out["token"]; ok {
		t.Error("absent field should not be injected")
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	m, _ := New([]Rule{{Field: "password"}})
	in := base()
	m.Apply(in)
	if in["password"] != "s3cr3t" {
		t.Error("Apply must not mutate the input record")
	}
}
