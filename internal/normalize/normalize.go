// Package normalize provides field-level normalization for log records,
// supporting operations such as lowercase, uppercase, and trimming whitespace.
package normalize

import (
	"fmt"
	"strings"
)

// Op represents a normalization operation.
type Op string

const (
	OpLowercase Op = "lowercase"
	OpUppercase Op = "uppercase"
	OpTrimSpace Op = "trim_space"
)

// Rule defines a normalization rule for a single field.
type Rule struct {
	Field string
	Op    Op
}

// Normalizer applies a set of normalization rules to log records.
type Normalizer struct {
	rules []Rule
}

// New creates a Normalizer from the given rules.
// Returns an error if any rule has an empty field or unknown operation.
func New(rules []Rule) (*Normalizer, error) {
	for _, r := range rules {
		if r.Field == "" {
			return nil, fmt.Errorf("normalize: rule has empty field name")
		}
		switch r.Op {
		case OpLowercase, OpUppercase, OpTrimSpace:
		default:
			return nil, fmt.Errorf("normalize: unknown operation %q for field %q", r.Op, r.Field)
		}
	}
	return &Normalizer{rules: rules}, nil
}

// Apply runs all normalization rules against rec, modifying string fields in place.
// Non-string fields and missing fields are silently skipped.
func (n *Normalizer) Apply(rec map[string]any) map[string]any {
	for _, r := range n.rules {
		v, ok := rec[r.Field]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok {
			continue
		}
		switch r.Op {
		case OpLowercase:
			rec[r.Field] = strings.ToLower(s)
		case OpUppercase:
			rec[r.Field] = strings.ToUpper(s)
		case OpTrimSpace:
			rec[r.Field] = strings.TrimSpace(s)
		}
	}
	return rec
}
