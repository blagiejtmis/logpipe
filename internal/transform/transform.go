// Package transform provides log record transformation capabilities,
// allowing fields to be added, removed, or renamed before routing to sinks.
package transform

import (
	"fmt"
	"regexp"
)

// Rule describes a single transformation to apply to a log record.
type Rule struct {
	// Op is one of: "add", "remove", "rename", "replace".
	Op string `yaml:"op"`
	// Field is the target field name.
	Field string `yaml:"field"`
	// Value is used by "add" and "replace" ops.
	Value string `yaml:"value,omitempty"`
	// NewName is used by the "rename" op.
	NewName string `yaml:"new_name,omitempty"`
	// Pattern is a regex used by "replace" to match the existing value.
	Pattern string `yaml:"pattern,omitempty"`
}

// Transformer applies a sequence of Rules to log records.
type Transformer struct {
	rules   []Rule
	patterns map[int]*regexp.Regexp
}

// New validates and compiles the provided rules, returning a Transformer.
func New(rules []Rule) (*Transformer, error) {
	patterns := make(map[int]*regexp.Regexp)
	for i, r := range rules {
		switch r.Op {
		case "add", "remove", "rename", "replace":
			// valid
		default:
			return nil, fmt.Errorf("transform: unknown op %q at index %d", r.Op, i)
		}
		if r.Field == "" {
			return nil, fmt.Errorf("transform: rule at index %d missing field", i)
		}
		if r.Op == "replace" && r.Pattern != "" {
			re, err := regexp.Compile(r.Pattern)
			if err != nil {
				return nil, fmt.Errorf("transform: rule %d bad pattern: %w", i, err)
			}
			patterns[i] = re
		}
	}
	return &Transformer{rules: rules, patterns: patterns}, nil
}

// Apply returns a new map with all rules applied to record.
func (t *Transformer) Apply(record map[string]string) map[string]string {
	out := make(map[string]string, len(record))
	for k, v := range record {
		out[k] = v
	}
	for i, r := range t.rules {
		switch r.Op {
		case "add":
			out[r.Field] = r.Value
		case "remove":
			delete(out, r.Field)
		case "rename":
			if val, ok := out[r.Field]; ok {
				out[r.NewName] = val
				delete(out, r.Field)
			}
		case "replace":
			if val, ok := out[r.Field]; ok {
				if re, hasRe := t.patterns[i]; hasRe {
					out[r.Field] = re.ReplaceAllString(val, r.Value)
				} else {
					out[r.Field] = r.Value
				}
			}
		}
	}
	return out
}
