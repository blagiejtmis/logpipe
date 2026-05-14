package clone

import (
	"testing"
)

func base() Record {
	return Record{
		"level":   "info",
		"message": "hello",
		"count":   int64(42),
		"meta": map[string]any{
			"host": "srv-1",
			"tags": []any{"a", "b"},
		},
	}
}

func TestDeep_ReturnsCopy(t *testing.T) {
	orig := base()
	copy, err := Deep(orig)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if copy["level"] != orig["level"] {
		t.Errorf("expected level %q, got %q", orig["level"], copy["level"])
	}
}

func TestDeep_MutationDoesNotAffectOriginal(t *testing.T) {
	orig := base()
	copy, _ := Deep(orig)
	copy["level"] = "error"
	if orig["level"] != "info" {
		t.Errorf("original was mutated")
	}
}

func TestDeep_NestedMapIsIndependent(t *testing.T) {
	orig := base()
	copy, _ := Deep(orig)
	copy["meta"].(map[string]any)["host"] = "srv-99"
	origMeta := orig["meta"].(map[string]any)
	if origMeta["host"] != "srv-1" {
		t.Errorf("nested map in original was mutated")
	}
}

func TestDeep_NestedSliceIsIndependent(t *testing.T) {
	orig := base()
	copy, _ := Deep(orig)
	slice := copy["meta"].(map[string]any)["tags"].([]any)
	slice[0] = "z"
	origTags := orig["meta"].(map[string]any)["tags"].([]any)
	if origTags[0] != "a" {
		t.Errorf("nested slice in original was mutated")
	}
}

func TestDeep_NilValueAllowed(t *testing.T) {
	r := Record{"x": nil}
	copy, err := Deep(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if copy["x"] != nil {
		t.Errorf("expected nil, got %v", copy["x"])
	}
}

func TestDeep_UnsupportedType_ReturnsError(t *testing.T) {
	type custom struct{ V int }
	r := Record{"obj": custom{V: 1}}
	_, err := Deep(r)
	if err == nil {
		t.Fatal("expected error for unsupported type, got nil")
	}
}

func TestMustDeep_Panics_OnUnsupportedType(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic, got none")
		}
	}()
	type custom struct{}
	MustDeep(Record{"x": custom{}})
}
