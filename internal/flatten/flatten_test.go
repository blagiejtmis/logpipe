package flatten

import (
	"testing"
)

func base() map[string]any {
	return map[string]any{
		"message": "hello",
		"meta": map[string]any{
			"host": "srv1",
			"region": "us-east",
		},
	}
}

func TestNew_NoRules_ReturnsError(t *testing.T) {
	_, err := New(nil)
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestNew_BlankField_ReturnsError(t *testing.T) {
	_, err := New([]Rule{{Field: "  ", Separator: "."}})
	if err == nil {
		t.Fatal("expected error for blank field")
	}
}

func TestNew_EmptySeparator_DefaultsToDot(t *testing.T) {
	f, err := New([]Rule{{Field: "meta"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.rules[0].Separator != "." {
		t.Errorf("expected separator '.', got %q", f.rules[0].Separator)
	}
}

func TestApply_FlattensNestedMap(t *testing.T) {
	f, _ := New([]Rule{{Field: "meta", Separator: "."}})
	rec := base()
	out := f.Apply(rec)

	if out["meta.host"] != "srv1" {
		t.Errorf("expected meta.host=srv1, got %v", out["meta.host"])
	}
	if out["meta.region"] != "us-east" {
		t.Errorf("expected meta.region=us-east, got %v", out["meta.region"])
	}
	if _, exists := out["meta"]; !exists {
		t.Error("source field 'meta' should be preserved when DropSource=false")
	}
}

func TestApply_DropSource(t *testing.T) {
	f, _ := New([]Rule{{Field: "meta", DropSource: true}})
	rec := base()
	out := f.Apply(rec)

	if _, exists := out["meta"]; exists {
		t.Error("source field 'meta' should be removed when DropSource=true")
	}
	if out["meta.host"] != "srv1" {
		t.Errorf("expected meta.host=srv1, got %v", out["meta.host"])
	}
}

func TestApply_CustomPrefix(t *testing.T) {
	f, _ := New([]Rule{{Field: "meta", Prefix: "tags", Separator: "_"}})
	rec := base()
	out := f.Apply(rec)

	if out["tags_host"] != "srv1" {
		t.Errorf("expected tags_host=srv1, got %v", out["tags_host"])
	}
}

func TestApply_DeepNesting(t *testing.T) {
	f, _ := New([]Rule{{Field: "deep", DropSource: true}})
	rec := map[string]any{
		"deep": map[string]any{
			"a": map[string]any{
				"b": "value",
			},
		},
	}
	out := f.Apply(rec)
	if out["a.b"] != "value" {
		t.Errorf("expected a.b=value, got %v", out["a.b"])
	}
}

func TestApply_NonMapField_Skipped(t *testing.T) {
	f, _ := New([]Rule{{Field: "message"}})
	rec := base()
	out := f.Apply(rec)
	if out["message"] != "hello" {
		t.Error("non-map field should be left unchanged")
	}
}

func TestApply_MissingField_NoOp(t *testing.T) {
	f, _ := New([]Rule{{Field: "absent"}})
	rec := base()
	out := f.Apply(rec)
	if len(out) != len(base()) {
		t.Error("missing field should not alter record")
	}
}
