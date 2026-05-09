package normalize

import (
	"testing"
)

func base() map[string]any {
	return map[string]any{
		"level":   "INFO",
		"message": "  hello world  ",
		"source":  "myapp",
		"count":   42,
	}
}

func TestNew_EmptyField_ReturnsError(t *testing.T) {
	_, err := New([]Rule{{Field: "", Op: OpLowercase}})
	if err == nil {
		t.Fatal("expected error for empty field name")
	}
}

func TestNew_UnknownOp_ReturnsError(t *testing.T) {
	_, err := New([]Rule{{Field: "level", Op: "explode"}})
	if err == nil {
		t.Fatal("expected error for unknown op")
	}
}

func TestNew_ValidRules_NoError(t *testing.T) {
	_, err := New([]Rule{
		{Field: "level", Op: OpLowercase},
		{Field: "source", Op: OpUppercase},
		{Field: "message", Op: OpTrimSpace},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApply_Lowercase(t *testing.T) {
	n, _ := New([]Rule{{Field: "level", Op: OpLowercase}})
	rec := base()
	out := n.Apply(rec)
	if out["level"] != "info" {
		t.Errorf("expected 'info', got %q", out["level"])
	}
}

func TestApply_Uppercase(t *testing.T) {
	n, _ := New([]Rule{{Field: "source", Op: OpUppercase}})
	rec := base()
	out := n.Apply(rec)
	if out["source"] != "MYAPP" {
		t.Errorf("expected 'MYAPP', got %q", out["source"])
	}
}

func TestApply_TrimSpace(t *testing.T) {
	n, _ := New([]Rule{{Field: "message", Op: OpTrimSpace}})
	rec := base()
	out := n.Apply(rec)
	if out["message"] != "hello world" {
		t.Errorf("expected 'hello world', got %q", out["message"])
	}
}

func TestApply_MissingField_Skipped(t *testing.T) {
	n, _ := New([]Rule{{Field: "nonexistent", Op: OpLowercase}})
	rec := base()
	out := n.Apply(rec)
	if _, ok := out["nonexistent"]; ok {
		t.Error("expected missing field to remain absent")
	}
}

func TestApply_NonStringField_Skipped(t *testing.T) {
	n, _ := New([]Rule{{Field: "count", Op: OpLowercase}})
	rec := base()
	out := n.Apply(rec)
	if out["count"] != 42 {
		t.Errorf("expected count to remain 42, got %v", out["count"])
	}
}

func TestApply_MultipleRules(t *testing.T) {
	n, _ := New([]Rule{
		{Field: "level", Op: OpLowercase},
		{Field: "message", Op: OpTrimSpace},
	})
	rec := base()
	out := n.Apply(rec)
	if out["level"] != "info" {
		t.Errorf("level: expected 'info', got %q", out["level"])
	}
	if out["message"] != "hello world" {
		t.Errorf("message: expected 'hello world', got %q", out["message"])
	}
}
