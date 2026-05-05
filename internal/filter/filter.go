// Package filter provides log line filtering based on configurable rules.
package filter

import (
	"regexp"
	"strings"
)

// Rule defines a single filter rule applied to log lines.
type Rule struct {
	// Field is the log field to match against (e.g. "level", "message").
	// Use "_raw" to match against the entire raw line.
	Field string `yaml:"field"`
	// Pattern is a regular expression pattern to match.
	Pattern string `yaml:"pattern"`
	// Invert inverts the match (i.e. exclude lines that match).
	Invert bool `yaml:"invert"`
}

// Filter evaluates log lines against a set of rules.
// All rules must match for a line to pass (AND semantics).
type Filter struct {
	rules   []Rule
	compiled []*regexp.Regexp
}

// New creates a Filter from the provided rules.
// Returns an error if any pattern fails to compile.
func New(rules []Rule) (*Filter, error) {
	compiled := make([]*regexp.Regexp, len(rules))
	for i, r := range rules {
		re, err := regexp.Compile(r.Pattern)
		if err != nil {
			return nil, err
		}
		compiled[i] = re
	}
	return &Filter{rules: rules, compiled: compiled}, nil
}

// Match reports whether the given fields (or raw line) satisfy all rules.
// fields is a map of field name to value; use key "_raw" for the full line.
func (f *Filter) Match(fields map[string]string) bool {
	for i, rule := range f.rules {
		val, ok := fields[rule.Field]
		if !ok {
			// Field not present; treat as empty string.
			val = ""
		}
		matched := f.compiled[i].MatchString(strings.TrimSpace(val))
		if rule.Invert {
			matched = !matched
		}
		if !matched {
			return false
		}
	}
	return true
}

// NoRules reports whether the filter has no rules (i.e. passes everything).
func (f *Filter) NoRules() bool {
	return len(f.rules) == 0
}
