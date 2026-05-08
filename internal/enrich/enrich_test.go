package enrich

import (
	"testing"
)

func base() map[string]any {
	return map[string]any{"message": "hello", "level": "info"}
}

func TestNew_EmptyFieldName_ReturnsError(t *testing.T) {
	_, err := New([]Rule{{Field: "", Value: "x"}})
	if err == nil {
		t.Fatal("expected error for empty field name")
	}
}

func TestNew_ValidRules_NoError(t *testing.T) {
	_, err := New([]Rule{{Field: "env", Value: "prod"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApply_AddsField(t *testing.T) {
	e, _ := New([]Rule{{Field: "env", Value: "staging"}})
	out := e.Apply(base())
	if out["env"] != "staging" {
		t.Errorf("expected env=staging, got %v", out["env"])
	}
}

func TestApply_DoesNotOverwriteExisting(t *testing.T) {
	e, _ := New([]Rule{{Field: "level", Value: "debug"}})
	out := e.Apply(base())
	if out["level"] != "info" {
		t.Errorf("expected level=info (unchanged), got %v", out["level"])
	}
}

func TestApply_MultipleRules(t *testing.T) {
	e, _ := New([]Rule{
		{Field: "env", Value: "prod"},
		{Field: "region", Value: "us-east-1"},
	})
	out := e.Apply(base())
	if out["env"] != "prod" {
		t.Errorf("expected env=prod, got %v", out["env"])
	}
	if out["region"] != "us-east-1" {
		t.Errorf("expected region=us-east-1, got %v", out["region"])
	}
}

func TestApply_OriginalUnmodified(t *testing.T) {
	e, _ := New([]Rule{{Field: "env", Value: "prod"}})
	orig := base()
	e.Apply(orig)
	if _, ok := orig["env"]; ok {
		t.Error("original record should not be modified")
	}
}

func TestApply_NoRules_ReturnsCopy(t *testing.T) {
	e, _ := New(nil)
	out := e.Apply(base())
	if out["message"] != "hello" {
		t.Errorf("expected message=hello, got %v", out["message"])
	}
	if len(out) != len(base()) {
		t.Errorf("expected same length, got %d", len(out))
	}
}
