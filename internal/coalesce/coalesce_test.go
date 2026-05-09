package coalesce

import (
	"testing"
)

func base() map[string]any {
	return map[string]any{
		"msg":     "hello",
		"message": "",
		"text":    "fallback",
	}
}

func TestNew_NoRules_ReturnsError(t *testing.T) {
	_, err := New(nil)
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestNew_EmptySources_ReturnsError(t *testing.T) {
	_, err := New([]Rule{{Sources: nil, Dest: "out"}})
	if err == nil {
		t.Fatal("expected error for empty sources")
	}
}

func TestNew_EmptyDest_ReturnsError(t *testing.T) {
	_, err := New([]Rule{{Sources: []string{"a"}, Dest: ""}})
	if err == nil {
		t.Fatal("expected error for empty dest")
	}
}

func TestApply_FirstNonEmpty(t *testing.T) {
	c, err := New([]Rule{
		{Sources: []string{"message", "msg", "text"}, Dest: "message"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rec := base()
	c.Apply(rec)
	if rec["message"] != "hello" {
		t.Errorf("expected 'hello', got %v", rec["message"])
	}
}

func TestApply_SkipsMissingFields(t *testing.T) {
	c, _ := New([]Rule{
		{Sources: []string{"absent", "text"}, Dest: "out"},
	})
	rec := base()
	c.Apply(rec)
	if rec["out"] != "fallback" {
		t.Errorf("expected 'fallback', got %v", rec["out"])
	}
}

func TestApply_RemovesSourcesByDefault(t *testing.T) {
	c, _ := New([]Rule{
		{Sources: []string{"msg", "text"}, Dest: "canonical"},
	})
	rec := base()
	c.Apply(rec)
	if _, ok := rec["msg"]; ok {
		t.Error("expected msg to be removed")
	}
}

func TestApply_KeepSources(t *testing.T) {
	c, _ := New([]Rule{
		{Sources: []string{"msg", "text"}, Dest: "canonical", KeepSources: true},
	})
	rec := base()
	c.Apply(rec)
	if _, ok := rec["msg"]; !ok {
		t.Error("expected msg to be retained")
	}
}

func TestApply_NothingChosen_RecordUnchanged(t *testing.T) {
	c, _ := New([]Rule{
		{Sources: []string{"absent1", "absent2"}, Dest: "out"},
	})
	rec := map[string]any{"level": "info"}
	c.Apply(rec)
	if _, ok := rec["out"]; ok {
		t.Error("expected out to not be set")
	}
}
