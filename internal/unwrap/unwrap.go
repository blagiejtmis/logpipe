// Package unwrap provides a processor that extracts a nested object field
// and promotes its keys to the top-level record, optionally under a prefix.
package unwrap

import (
	"fmt"
	"strings"
)

// Rule describes a single unwrap operation.
type Rule struct {
	// Field is the dot-separated path to the nested map field to unwrap.
	Field string
	// Prefix is prepended to each promoted key. If empty, keys are promoted as-is.
	Prefix string
	// Remove controls whether the original nested field is deleted after unwrapping.
	Remove bool
}

// Unwrapper promotes nested map fields to the top level of a log record.
type Unwrapper struct {
	rules []Rule
}

// New validates rules and returns an Unwrapper.
func New(rules []Rule) (*Unwrapper, error) {
	if len(rules) == 0 {
		return nil, fmt.Errorf("unwrap: at least one rule is required")
	}
	for i, r := range rules {
		if strings.TrimSpace(r.Field) == "" {
			return nil, fmt.Errorf("unwrap: rule[%d]: field must not be blank", i)
		}
	}
	return &Unwrapper{rules: rules}, nil
}

// Apply runs all rules against record and returns the modified copy.
func (u *Unwrapper) Apply(record map[string]any) map[string]any {
	for _, r := range u.rules {
		val, ok := record[r.Field]
		if !ok {
			continue
		}
		nested, ok := val.(map[string]any)
		if !ok {
			continue
		}
		for k, v := range nested {
			key := k
			if r.Prefix != "" {
				key = r.Prefix + k
			}
			record[key] = v
		}
		if r.Remove {
			delete(record, r.Field)
		}
	}
	return record
}
