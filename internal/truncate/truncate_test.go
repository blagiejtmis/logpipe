package truncate

import (
	"testing"
)

func base() map[string]any {
	return map[string]any{
		"message": "hello world",
		"detail":  "some longer detail text",
		"count":   42,
	}
}

func TestNew_NoRules_ReturnsError(t *testing.T) {
	_, err := New(nil)
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestNew_EmptyField_ReturnsError(t *testing.T) {
	_, err := New([]Rule{{Field: "", MaxLen: 10}})
	if err == nil {
		t.Fatal("expected error for blank field")
	}
}

func TestNew_ZeroMaxLen_ReturnsError(t *testing.T) {
	_, err := New([]Rule{{Field: "message", MaxLen: 0}})
	if err == nil {
		t.Fatal("expected error for zero max_len")
	}
}

func TestNew_ValidRules_NoError(t *testing.T) {
	_, err := New([]Rule{{Field: "message", MaxLen: 5}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApply_NoTruncation_WhenUnderLimit(t *testing.T) {
	tr, _ := New([]Rule{{Field: "message", MaxLen: 100}})
	rec := base()
	tr.Apply(rec)
	if rec["message"] != "hello world" {
		t.Errorf("expected unchanged message, got %q", rec["message"])
	}
}

func TestApply_TruncatesAtMaxLen(t *testing.T) {
	tr, _ := New([]Rule{{Field: "message", MaxLen: 5}})
	rec := base()
	tr.Apply(rec)
	if rec["message"] != "hello" {
		t.Errorf("expected %q, got %q", "hello", rec["message"])
	}
}

func TestApply_AppendsSuffix(t *testing.T) {
	tr, _ := New([]Rule{{Field: "message", MaxLen: 5, Suffix: "..."}})
	rec := base()
	tr.Apply(rec)
	if rec["message"] != "hello..." {
		t.Errorf("expected %q, got %q", "hello...", rec["message"])
	}
}

func TestApply_NonStringField_Skipped(t *testing.T) {
	tr, _ := New([]Rule{{Field: "count", MaxLen: 1}})
	rec := base()
	tr.Apply(rec)
	if rec["count"] != 42 {
		t.Errorf("expected count unchanged, got %v", rec["count"])
	}
}

func TestApply_MissingField_Skipped(t *testing.T) {
	tr, _ := New([]Rule{{Field: "nonexistent", MaxLen: 5}})
	rec := base()
	tr.Apply(rec) // must not panic
}

func TestApply_MultipleRules(t *testing.T) {
	tr, _ := New([]Rule{
		{Field: "message", MaxLen: 5},
		{Field: "detail", MaxLen: 4, Suffix: "~"},
	})
	rec := base()
	tr.Apply(rec)
	if rec["message"] != "hello" {
		t.Errorf("message: expected %q, got %q", "hello", rec["message"])
	}
	if rec["detail"] != "some~" {
		t.Errorf("detail: expected %q, got %q", "some~", rec["detail"])
	}
}
