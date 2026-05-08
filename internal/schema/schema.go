// Package schema provides log record field validation against
// a declared schema, ensuring required fields are present and
// that field values match expected types.
package schema

import (
	"fmt"
	"strings"
)

// FieldType represents the expected type of a schema field.
type FieldType string

const (
	FieldTypeString  FieldType = "string"
	FieldTypeNumber  FieldType = "number"
	FieldTypeBool    FieldType = "bool"
)

// FieldDef describes a single field in the schema.
type FieldDef struct {
	Name     string
	Type     FieldType
	Required bool
}

// Config holds the schema configuration.
type Config struct {
	Fields []FieldDef
}

// Validator validates log records against a schema.
type Validator struct {
	fields []FieldDef
}

// New creates a new Validator from the given Config.
// Returns an error if any FieldDef has an unsupported type.
func New(cfg Config) (*Validator, error) {
	for _, f := range cfg.Fields {
		switch f.Type {
		case FieldTypeString, FieldTypeNumber, FieldTypeBool:
			// valid
		default:
			return nil, fmt.Errorf("schema: unknown field type %q for field %q", f.Type, f.Name)
		}
		if strings.TrimSpace(f.Name) == "" {
			return nil, fmt.Errorf("schema: field name must not be empty")
		}
	}
	return &Validator{fields: cfg.Fields}, nil
}

// Validate checks the record against the schema.
// Returns a non-nil error describing the first violation found.
func (v *Validator) Validate(record map[string]any) error {
	for _, f := range v.fields {
		val, ok := record[f.Name]
		if !ok {
			if f.Required {
				return fmt.Errorf("schema: required field %q is missing", f.Name)
			}
			continue
		}
		if err := checkType(f.Name, f.Type, val); err != nil {
			return err
		}
	}
	return nil
}

func checkType(name string, want FieldType, val any) error {
	switch want {
	case FieldTypeString:
		if _, ok := val.(string); !ok {
			return fmt.Errorf("schema: field %q must be a string, got %T", name, val)
		}
	case FieldTypeNumber:
		switch val.(type) {
		case int, int64, float64, float32:
			// ok
		default:
			return fmt.Errorf("schema: field %q must be a number, got %T", name, val)
		}
	case FieldTypeBool:
		if _, ok := val.(bool); !ok {
			return fmt.Errorf("schema: field %q must be a bool, got %T", name, val)
		}
	}
	return nil
}
