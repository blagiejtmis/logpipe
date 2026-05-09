// Package masking provides field-level value masking for log records,
// allowing sensitive fields to be partially or fully obscured before routing.
package masking

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule describes how a single field should be masked.
type Rule struct {
	// Field is the record key to mask.
	Field string
	// Pattern is an optional regex; only matching portions are masked.
	// If empty the entire value is replaced.
	Pattern string
	// Placeholder replaces the matched text. Defaults to "***".
	Placeholder string

	re *regexp.Regexp
}

// Masker applies a set of masking rules to log records.
type Masker struct {
	rules []Rule
}

// New validates and compiles the supplied rules, returning a Masker.
func New(rules []Rule) (*Masker, error) {
	compiled := make([]Rule, len(rules))
	for i, r := range rules {
		if strings.TrimSpace(r.Field) == "" {
			return nil, fmt.Errorf("masking: rule %d: field name must not be empty", i)
		}
		if r.Placeholder == "" {
			r.Placeholder = "***"
		}
		if r.Pattern != "" {
			re, err := regexp.Compile(r.Pattern)
			if err != nil {
				return nil, fmt.Errorf("masking: rule %d: invalid pattern: %w", i, err)
			}
			r.re = re
		}
		compiled[i] = r
	}
	return &Masker{rules: compiled}, nil
}

// Apply returns a shallow copy of record with masked field values.
func (m *Masker) Apply(record map[string]any) map[string]any {
	out := make(map[string]any, len(record))
	for k, v := range record {
		out[k] = v
	}
	for _, r := range m.rules {
		v, ok := out[r.Field]
		if !ok {
			continue
		}
		s := fmt.Sprintf("%v", v)
		if r.re != nil {
			out[r.Field] = r.re.ReplaceAllString(s, r.Placeholder)
		} else {
			out[r.Field] = r.Placeholder
		}
	}
	return out
}
