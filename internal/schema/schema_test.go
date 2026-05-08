package schema

import (
	"testing"
)

func TestNew_UnknownType_ReturnsError(t *testing.T) {
	_, err := New(Config{
		Fields: []FieldDef{{Name: "level", Type: "integer", Required: true}},
	})
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func TestNew_EmptyFieldName_ReturnsError(t *testing.T) {
	_, err := New(Config{
		Fields: []FieldDef{{Name: "  ", Type: FieldTypeString}},
	})
	if err == nil {
		t.Fatal("expected error for empty field name")
	}
}

func TestValidate_RequiredFieldMissing(t *testing.T) {
	v, err := New(Config{
		Fields: []FieldDef{{Name: "level", Type: FieldTypeString, Required: true}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := v.Validate(map[string]any{"msg": "hello"}); err == nil {
		t.Fatal("expected missing required field error")
	}
}

func TestValidate_OptionalFieldMissing_Passes(t *testing.T) {
	v, _ := New(Config{
		Fields: []FieldDef{{Name: "trace_id", Type: FieldTypeString, Required: false}},
	})
	if err := v.Validate(map[string]any{"msg": "hello"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_WrongType_ReturnsError(t *testing.T) {
	v, _ := New(Config{
		Fields: []FieldDef{{Name: "latency", Type: FieldTypeNumber, Required: true}},
	})
	err := v.Validate(map[string]any{"latency": "fast"})
	if err == nil {
		t.Fatal("expected type mismatch error")
	}
}

func TestValidate_CorrectTypes_Passes(t *testing.T) {
	v, _ := New(Config{
		Fields: []FieldDef{
			{Name: "msg", Type: FieldTypeString, Required: true},
			{Name: "latency", Type: FieldTypeNumber, Required: true},
			{Name: "ok", Type: FieldTypeBool, Required: true},
		},
	})
	err := v.Validate(map[string]any{
		"msg":     "hello",
		"latency": float64(42),
		"ok":      true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_BoolWrongType_ReturnsError(t *testing.T) {
	v, _ := New(Config{
		Fields: []FieldDef{{Name: "ok", Type: FieldTypeBool, Required: true}},
	})
	if err := v.Validate(map[string]any{"ok": "yes"}); err == nil {
		t.Fatal("expected type mismatch error for bool field")
	}
}

func TestValidate_NoFields_AlwaysPasses(t *testing.T) {
	v, _ := New(Config{})
	if err := v.Validate(map[string]any{"anything": 123}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
