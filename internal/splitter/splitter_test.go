package splitter

import (
	"testing"
)

func base() Record {
	return Record{
		"source": "app",
		"level":  "info",
		"tags":   []any{"a", "b", "c"},
	}
}

func TestNew_EmptyField_ReturnsError(t *testing.T) {
	_, err := New(Config{Field: ""})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestNew_ValidConfig_NoError(t *testing.T) {
	s, err := New(Config{Field: "tags"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil splitter")
	}
}

func TestSplit_ArrayField_ExpandsRecords(t *testing.T) {
	s, _ := New(Config{Field: "tags"})
	out := s.Split(base())
	if len(out) != 3 {
		t.Fatalf("expected 3 records, got %d", len(out))
	}
	for i, want := range []string{"a", "b", "c"} {
		got, ok := out[i]["tags"]
		if !ok {
			t.Fatalf("record %d missing 'tags' field", i)
		}
		if got != want {
			t.Errorf("record %d: got tags=%q, want %q", i, got, want)
		}
	}
}

func TestSplit_IndexField_IsSet(t *testing.T) {
	s, _ := New(Config{Field: "tags"})
	out := s.Split(base())
	for i, rec := range out {
		want := []string{"0", "1", "2"}[i]
		if rec["_index"] != want {
			t.Errorf("record %d: _index=%v, want %q", i, rec["_index"], want)
		}
	}
}

func TestSplit_PrefixedIndexField(t *testing.T) {
	s, _ := New(Config{Field: "tags", Prefix: "item"})
	out := s.Split(base())
	if _, ok := out[0]["item_index"]; !ok {
		t.Error("expected 'item_index' field, not found")
	}
	if _, ok := out[0]["_index"]; ok {
		t.Error("unexpected '_index' field when prefix is set")
	}
}

func TestSplit_OtherFieldsCopied(t *testing.T) {
	s, _ := New(Config{Field: "tags"})
	out := s.Split(base())
	for i, rec := range out {
		if rec["source"] != "app" {
			t.Errorf("record %d: source=%v, want 'app'", i, rec["source"])
		}
		if rec["level"] != "info" {
			t.Errorf("record %d: level=%v, want 'info'", i, rec["level"])
		}
	}
}

func TestSplit_MissingField_PassesThrough(t *testing.T) {
	s, _ := New(Config{Field: "missing"})
	rec := Record{"msg": "hello"}
	out := s.Split(rec)
	if len(out) != 1 {
		t.Fatalf("expected 1 record, got %d", len(out))
	}
	if out[0]["msg"] != "hello" {
		t.Errorf("unexpected record: %v", out[0])
	}
}

func TestSplit_NonSliceField_PassesThrough(t *testing.T) {
	s, _ := New(Config{Field: "level"})
	out := s.Split(base())
	if len(out) != 1 {
		t.Fatalf("expected 1 record, got %d", len(out))
	}
}

func TestSplit_StringSlice_Expanded(t *testing.T) {
	s, _ := New(Config{Field: "tags"})
	rec := Record{"tags": []string{"x", "y"}}
	out := s.Split(rec)
	if len(out) != 2 {
		t.Fatalf("expected 2 records, got %d", len(out))
	}
}
