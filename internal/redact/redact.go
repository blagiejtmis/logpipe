// Package redact provides field-level redaction for log records.
// It supports masking sensitive fields by replacing their values
// with a configurable placeholder or a partial mask.
package redact

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule describes a single redaction rule.
type Rule struct {
	// Field is the record key to redact.
	Field string
	// Pattern, if non-empty, only redacts values matching this regex.
	Pattern string
	// Placeholder is the replacement string. Defaults to "***".
	Placeholder string

	re *regexp.Regexp
}

// Redactor applies a set of redaction rules to log records.
type Redactor struct {
	rules []Rule
}

// New compiles and returns a Redactor from the provided rules.
// Returns an error if any rule contains an invalid regex pattern.
func New(rules []Rule) (*Redactor, error) {
	compiled := make([]Rule, len(rules))
	for i, r := range rules {
		if r.Placeholder == "" {
			r.Placeholder = "***"
		}
		if r.Pattern != "" {
			re, err := regexp.Compile(r.Pattern)
			if err != nil {
				return nil, fmt.Errorf("redact: rule %d invalid pattern %q: %w", i, r.Pattern, err)
			}
			r.re = re
		}
		compiled[i] = r
	}
	return &Redactor{rules: compiled}, nil
}

// Apply returns a copy of record with sensitive fields redacted according
// to the configured rules. The original map is never mutated.
func (r *Redactor) Apply(record map[string]string) map[string]string {
	if len(r.rules) == 0 {
		return record
	}
	out := make(map[string]string, len(record))
	for k, v := range record {
		out[k] = v
	}
	for _, rule := range r.rules {
		v, ok := out[rule.Field]
		if !ok {
			continue
		}
		if rule.re != nil {
			if rule.re.MatchString(v) {
				out[rule.Field] = rule.re.ReplaceAllString(v, strings.Repeat("*", len(rule.Placeholder)))
			}
		} else {
			out[rule.Field] = rule.Placeholder
		}
	}
	return out
}
