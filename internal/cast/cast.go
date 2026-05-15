// Package cast provides field type-casting for log records.
// Rules specify a source field, a target type, and an optional destination
// field. When dest is omitted the source field is overwritten in-place.
package cast

import (
	"fmt"
	"strconv"

	"logpipe/internal/clone"
)

// Rule describes a single cast operation.
type Rule struct {
	Field  string // source field name
	Type   string // "string", "int", "float", "bool"
	Dest   string // destination field; defaults to Field when empty
}

// Caster applies a set of type-cast rules to a log record.
type Caster struct {
	rules []Rule
}

// New validates rules and returns a ready Caster.
func New(rules []Rule) (*Caster, error) {
	if len(rules) == 0 {
		return nil, fmt.Errorf("cast: at least one rule is required")
	}
	for i, r := range rules {
		if r.Field == "" {
			return nil, fmt.Errorf("cast: rule[%d]: field must not be empty", i)
		}
		switch r.Type {
		case "string", "int", "float", "bool":
		default:
			return nil, fmt.Errorf("cast: rule[%d]: unknown type %q", i, r.Type)
		}
	}
	return &Caster{rules: rules}, nil
}

// Apply returns a shallow-cloned record with the cast rules applied.
// Fields that are missing or cannot be converted are left unchanged.
func (c *Caster) Apply(rec map[string]any) map[string]any {
	out := clone.MustDeep(rec)
	for _, r := range c.rules {
		v, ok := out[r.Field]
		if !ok {
			continue
		}
		casted, err := castValue(v, r.Type)
		if err != nil {
			continue
		}
		dest := r.Dest
		if dest == "" {
			dest = r.Field
		}
		out[dest] = casted
	}
	return out
}

func castValue(v any, typ string) (any, error) {
	s := fmt.Sprintf("%v", v)
	switch typ {
	case "string":
		return s, nil
	case "int":
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			// try float first then truncate
			f, err2 := strconv.ParseFloat(s, 64)
			if err2 != nil {
				return nil, err
			}
			return int64(f), nil
		}
		return n, nil
	case "float":
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}
		return f, nil
	case "bool":
		b, err := strconv.ParseBool(s)
		if err != nil {
			return nil, err
		}
		return b, nil
	}
	return nil, fmt.Errorf("unknown type")
}
