package redact_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/redact"
)

func baseRecord() map[string]string {
	return map[string]string{
		"user":     "alice",
		"password": "s3cr3t",
		"email":    "alice@example.com",
		"msg":      "login attempt",
	}
}

func TestNew_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := redact.New([]redact.Rule{
		{Field: "email", Pattern: "["},
	})
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestApply_NoRules_PassesThrough(t *testing.T) {
	r, err := redact.New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rec := baseRecord()
	out := r.Apply(rec)
	if out["password"] != "s3cr3t" {
		t.Errorf("expected password unchanged, got %q", out["password"])
	}
}

func TestApply_FullMask_DefaultPlaceholder(t *testing.T) {
	r, err := redact.New([]redact.Rule{
		{Field: "password"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := r.Apply(baseRecord())
	if out["password"] != "***" {
		t.Errorf("expected *** got %q", out["password"])
	}
	if out["user"] != "alice" {
		t.Errorf("unrelated field mutated")
	}
}

func TestApply_FullMask_CustomPlaceholder(t *testing.T) {
	r, err := redact.New([]redact.Rule{
		{Field: "password", Placeholder: "[REDACTED]"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := r.Apply(baseRecord())
	if out["password"] != "[REDACTED]" {
		t.Errorf("expected [REDACTED] got %q", out["password"])
	}
}

func TestApply_PatternMask_OnlyMatchingValue(t *testing.T) {
	r, err := redact.New([]redact.Rule{
		{Field: "email", Pattern: `[^@]+@[^@]+`},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := r.Apply(baseRecord())
	if out["email"] == "alice@example.com" {
		t.Error("expected email to be masked")
	}
}

func TestApply_MissingField_NoOp(t *testing.T) {
	r, err := redact.New([]redact.Rule{
		{Field: "token"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := r.Apply(baseRecord())
	if _, ok := out["token"]; ok {
		t.Error("missing field should not be introduced")
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	r, _ := redact.New([]redact.Rule{{Field: "password"}})
	orig := baseRecord()
	r.Apply(orig)
	if orig["password"] != "s3cr3t" {
		t.Error("original record was mutated")
	}
}
