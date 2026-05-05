package transform

import (
	"testing"
)

func baseRecord() map[string]string {
	return map[string]string{
		"level":   "info",
		"message": "hello world",
		"host":    "srv-01",
	}
}

func TestTransform_AddField(t *testing.T) {
	tr, err := New([]Rule{{Op: "add", Field: "env", Value: "production"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := tr.Apply(baseRecord())
	if out["env"] != "production" {
		t.Errorf("expected env=production, got %q", out["env"])
	}
}

func TestTransform_RemoveField(t *testing.T) {
	tr, _ := New([]Rule{{Op: "remove", Field: "host"}})
	out := tr.Apply(baseRecord())
	if _, ok := out["host"]; ok {
		t.Error("expected host to be removed")
	}
}

func TestTransform_RenameField(t *testing.T) {
	tr, _ := New([]Rule{{Op: "rename", Field: "message", NewName: "msg"}})
	out := tr.Apply(baseRecord())
	if out["msg"] != "hello world" {
		t.Errorf("expected msg=hello world, got %q", out["msg"])
	}
	if _, ok := out["message"]; ok {
		t.Error("original field should be removed after rename")
	}
}

func TestTransform_ReplaceWithPattern(t *testing.T) {
	tr, err := New([]Rule{{Op: "replace", Field: "level", Pattern: `^info$`, Value: "INFO"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := tr.Apply(baseRecord())
	if out["level"] != "INFO" {
		t.Errorf("expected level=INFO, got %q", out["level"])
	}
}

func TestTransform_ReplaceWithoutPattern(t *testing.T) {
	tr, _ := New([]Rule{{Op: "replace", Field: "host", Value: "srv-99"}})
	out := tr.Apply(baseRecord())
	if out["host"] != "srv-99" {
		t.Errorf("expected host=srv-99, got %q", out["host"])
	}
}

func TestTransform_OriginalUnmodified(t *testing.T) {
	tr, _ := New([]Rule{{Op: "add", Field: "extra", Value: "x"}})
	original := baseRecord()
	tr.Apply(original)
	if _, ok := original["extra"]; ok {
		t.Error("Apply should not mutate the original record")
	}
}

func TestNew_UnknownOp(t *testing.T) {
	_, err := New([]Rule{{Op: "uppercase", Field: "level"}})
	if err == nil {
		t.Error("expected error for unknown op")
	}
}

func TestNew_MissingField(t *testing.T) {
	_, err := New([]Rule{{Op: "add", Value: "x"}})
	if err == nil {
		t.Error("expected error when field is empty")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New([]Rule{{Op: "replace", Field: "level", Pattern: `[invalid`}})
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}
