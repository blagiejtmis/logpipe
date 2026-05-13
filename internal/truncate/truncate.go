// Package truncate provides field-level truncation for log records,
// capping string values at a configurable maximum byte length.
package truncate

import (
	"errors"
	"fmt"
)

// Rule describes how a single field should be truncated.
type Rule struct {
	// Field is the record key to truncate.
	Field string
	// MaxLen is the maximum number of bytes allowed. Values longer than
	// this are sliced and an optional suffix is appended.
	MaxLen int
	// Suffix is appended when a value is truncated (e.g. "...").
	// Defaults to empty string.
	Suffix string
}

// Truncator applies a set of truncation rules to log records.
type Truncator struct {
	rules []Rule
}

// New constructs a Truncator from the provided rules.
// Returns an error if any rule is invalid.
func New(rules []Rule) (*Truncator, error) {
	if len(rules) == 0 {
		return nil, errors.New("truncate: at least one rule is required")
	}
	for i, r := range rules {
		if r.Field == "" {
			return nil, fmt.Errorf("truncate: rule[%d]: field must not be empty", i)
		}
		if r.MaxLen <= 0 {
			return nil, fmt.Errorf("truncate: rule[%d]: max_len must be > 0", i)
		}
	}
	return &Truncator{rules: rules}, nil
}

// Apply truncates fields in rec according to the configured rules.
// rec is modified in place and returned.
func (t *Truncator) Apply(rec map[string]any) map[string]any {
	for _, r := range t.rules {
		v, ok := rec[r.Field]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok {
			continue
		}
		if len(s) > r.MaxLen {
			s = s[:r.MaxLen] + r.Suffix
			rec[r.Field] = s
		}
	}
	return rec
}
