package classify

import (
	"testing"
)

func base() map[string]any {
	return map[string]any{
		"message": "user login failed",
		"level":   "error",
	}
}

func TestNew_NoRules_ReturnsError(t *testing.T) {
	_, err := New(nil)
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestNew_EmptyField_ReturnsError(t *testing.T) {
	_, err := New([]Rule{{Field: "", Pattern: ".*", Label: "any"}})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_EmptyLabel_ReturnsError(t *testing.T) {
	_, err := New([]Rule{{Field: "message", Pattern: ".*", Label: ""}})
	if err == nil {
		t.Fatal("expected error for empty label")
	}
}

func TestNew_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := New([]Rule{{Field: "message", Pattern: "[", Label: "bad"}})
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestNew_ValidRules_NoError(t *testing.T) {
	_, err := New([]Rule{{Field: "message", Pattern: "login", Label: "auth"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApply_MatchesFirstRule(t *testing.T) {
	c, _ := New([]Rule{
		{Field: "message", Pattern: "login", Label: "auth"},
		{Field: "message", Pattern: "error", Label: "fault"},
	})
	rec := base()
	out := c.Apply(rec)
	if got := out["category"]; got != "auth" {
		t.Fatalf("expected auth, got %v", got)
	}
}

func TestApply_FallsToSecondRule(t *testing.T) {
	c, _ := New([]Rule{
		{Field: "message", Pattern: "^payment", Label: "billing"},
		{Field: "level", Pattern: "error", Label: "fault"},
	})
	out := c.Apply(base())
	if got := out["category"]; got != "fault" {
		t.Fatalf("expected fault, got %v", got)
	}
}

func TestApply_NoMatch_RecordUnchanged(t *testing.T) {
	c, _ := New([]Rule{
		{Field: "message", Pattern: "^payment", Label: "billing"},
	})
	rec := base()
	out := c.Apply(rec)
	if _, exists := out["category"]; exists {
		t.Fatal("expected no category field when no rule matches")
	}
}

func TestApply_CustomDestField(t *testing.T) {
	c, _ := New([]Rule{
		{Field: "message", Pattern: "login", Label: "auth", DestField: "kind"},
	})
	out := c.Apply(base())
	if got := out["kind"]; got != "auth" {
		t.Fatalf("expected kind=auth, got %v", got)
	}
	if _, exists := out["category"]; exists {
		t.Fatal("expected no default category field")
	}
}

func TestApply_MissingField_Skipped(t *testing.T) {
	c, _ := New([]Rule{
		{Field: "nonexistent", Pattern: ".*", Label: "x"},
		{Field: "level", Pattern: "error", Label: "fault"},
	})
	out := c.Apply(base())
	if got := out["category"]; got != "fault" {
		t.Fatalf("expected fault, got %v", got)
	}
}
