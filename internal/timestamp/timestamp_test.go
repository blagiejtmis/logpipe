package timestamp_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/timestamp"
)

func base() map[string]any {
	return map[string]any{
		"ts":  "2024-03-15T10:30:00Z",
		"msg": "hello",
	}
}

func TestNew_NoRules_ReturnsError(t *testing.T) {
	_, err := timestamp.New(nil)
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestNew_EmptyField_ReturnsError(t *testing.T) {
	_, err := timestamp.New([]timestamp.Rule{{Field: "", Layouts: []string{time.RFC3339}}})
	if err == nil {
		t.Fatal("expected error for blank field")
	}
}

func TestNew_NoLayouts_ReturnsError(t *testing.T) {
	_, err := timestamp.New([]timestamp.Rule{{Field: "ts", Layouts: nil}})
	if err == nil {
		t.Fatal("expected error for empty layouts")
	}
}

func TestNew_ValidRule_NoError(t *testing.T) {
	_, err := timestamp.New([]timestamp.Rule{{Field: "ts", Layouts: []string{time.RFC3339}}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApply_ParsesAndNormalizes(t *testing.T) {
	n, _ := timestamp.New([]timestamp.Rule{
		{Field: "ts", Layouts: []string{time.RFC3339}},
	})
	rec := base()
	n.Apply(rec)
	want := "2024-03-15T10:30:00Z"
	if rec["ts"] != want {
		t.Fatalf("got %q, want %q", rec["ts"], want)
	}
}

func TestApply_WritesToOutputField(t *testing.T) {
	n, _ := timestamp.New([]timestamp.Rule{
		{Field: "ts", Layouts: []string{time.RFC3339}, OutputField: "@timestamp"},
	})
	rec := base()
	n.Apply(rec)
	if _, ok := rec["@timestamp"]; !ok {
		t.Fatal("expected @timestamp field to be set")
	}
	if _, ok := rec["ts"]; !ok {
		t.Fatal("original ts field should be preserved")
	}
}

func TestApply_UnparsableValue_LeftUnchanged(t *testing.T) {
	n, _ := timestamp.New([]timestamp.Rule{
		{Field: "ts", Layouts: []string{time.RFC3339}},
	})
	rec := map[string]any{"ts": "not-a-timestamp"}
	n.Apply(rec)
	if rec["ts"] != "not-a-timestamp" {
		t.Fatalf("expected original value, got %v", rec["ts"])
	}
}

func TestApply_MissingField_NoOp(t *testing.T) {
	n, _ := timestamp.New([]timestamp.Rule{
		{Field: "ts", Layouts: []string{time.RFC3339}},
	})
	rec := map[string]any{"msg": "hi"}
	n.Apply(rec)
	if _, ok := rec["ts"]; ok {
		t.Fatal("ts field should not have been created")
	}
}

func TestApply_TriesLayoutsInOrder(t *testing.T) {
	n, _ := timestamp.New([]timestamp.Rule{
		{
			Field:   "ts",
			Layouts: []string{"2006-01-02", time.RFC3339},
		},
	})
	rec := map[string]any{"ts": "2024-03-15"}
	n.Apply(rec)
	if rec["ts"] != "2024-03-15T00:00:00Z" {
		t.Fatalf("unexpected value: %v", rec["ts"])
	}
}
