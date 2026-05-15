package cast

import (
	"testing"
)

func base() map[string]any {
	return map[string]any{
		"count":  "42",
		"ratio":  "3.14",
		"active": "true",
		"label":  123,
	}
}

func TestNew_NoRules_ReturnsError(t *testing.T) {
	_, err := New(nil)
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestNew_EmptyField_ReturnsError(t *testing.T) {
	_, err := New([]Rule{{Field: "", Type: "int"}})
	if err == nil {
		t.Fatal("expected error for blank field")
	}
}

func TestNew_UnknownType_ReturnsError(t *testing.T) {
	_, err := New([]Rule{{Field: "x", Type: "uuid"}})
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func TestNew_ValidRules_NoError(t *testing.T) {
	_, err := New([]Rule{{Field: "count", Type: "int"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApply_StringToInt(t *testing.T) {
	c, _ := New([]Rule{{Field: "count", Type: "int"}})
	out := c.Apply(base())
	if out["count"] != int64(42) {
		t.Fatalf("expected int64(42), got %v (%T)", out["count"], out["count"])
	}
}

func TestApply_StringToFloat(t *testing.T) {
	c, _ := New([]Rule{{Field: "ratio", Type: "float"}})
	out := c.Apply(base())
	if out["ratio"] != 3.14 {
		t.Fatalf("expected 3.14, got %v", out["ratio"])
	}
}

func TestApply_StringToBool(t *testing.T) {
	c, _ := New([]Rule{{Field: "active", Type: "bool"}})
	out := c.Apply(base())
	if out["active"] != true {
		t.Fatalf("expected true, got %v", out["active"])
	}
}

func TestApply_IntToString(t *testing.T) {
	c, _ := New([]Rule{{Field: "label", Type: "string"}})
	out := c.Apply(base())
	if out["label"] != "123" {
		t.Fatalf("expected \"123\", got %v", out["label"])
	}
}

func TestApply_DestField_WritesToNewKey(t *testing.T) {
	c, _ := New([]Rule{{Field: "count", Type: "int", Dest: "count_int"}})
	out := c.Apply(base())
	if out["count_int"] != int64(42) {
		t.Fatalf("expected int64(42) at count_int, got %v", out["count_int"])
	}
	// original key preserved
	if out["count"] != "42" {
		t.Fatalf("original field should be unchanged")
	}
}

func TestApply_MissingField_Skipped(t *testing.T) {
	c, _ := New([]Rule{{Field: "missing", Type: "int"}})
	out := c.Apply(base())
	if _, ok := out["missing"]; ok {
		t.Fatal("missing field should not be created")
	}
}

func TestApply_InvalidConversion_Skipped(t *testing.T) {
	c, _ := New([]Rule{{Field: "label", Type: "bool"}})
	out := c.Apply(base())
	// "123" cannot be parsed as bool; field should remain as-is
	if out["label"] != 123 {
		t.Fatalf("expected original value 123, got %v", out["label"])
	}
}

func TestApply_OriginalUnmodified(t *testing.T) {
	c, _ := New([]Rule{{Field: "count", Type: "int"}})
	orig := base()
	c.Apply(orig)
	if orig["count"] != "42" {
		t.Fatal("Apply must not mutate the original record")
	}
}
