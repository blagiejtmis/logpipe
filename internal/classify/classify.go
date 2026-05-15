// Package classify assigns a category label to log records based on
// field-matching rules. Each rule specifies a field, a regex pattern,
// and a label to write into a destination field when the pattern matches.
package classify

import (
	"errors"
	"fmt"
	"regexp"
)

// Rule describes a single classification rule.
type Rule struct {
	// Field is the source field to inspect.
	Field string
	// Pattern is the regular expression to match against the field value.
	Pattern string
	// Label is the value written to DestField when Pattern matches.
	Label string
	// DestField is the field where the label is stored (defaults to "category").
	DestField string
}

// Classifier applies ordered rules to a record and writes the first matching
// label into the destination field.
type Classifier struct {
	rules []compiledRule
}

type compiledRule struct {
	field     string
	re        *regexp.Regexp
	label     string
	destField string
}

// New builds a Classifier from the provided rules. Rules are evaluated in
// order; the first match wins. Returns an error if any rule is invalid.
func New(rules []Rule) (*Classifier, error) {
	if len(rules) == 0 {
		return nil, errors.New("classify: at least one rule is required")
	}
	compiled := make([]compiledRule, 0, len(rules))
	for i, r := range rules {
		if r.Field == "" {
			return nil, fmt.Errorf("classify: rule %d: field must not be empty", i)
		}
		if r.Label == "" {
			return nil, fmt.Errorf("classify: rule %d: label must not be empty", i)
		}
		re, err := regexp.Compile(r.Pattern)
		if err != nil {
			return nil, fmt.Errorf("classify: rule %d: invalid pattern: %w", i, err)
		}
		dest := r.DestField
		if dest == "" {
			dest = "category"
		}
		compiled = append(compiled, compiledRule{
			field:     r.Field,
			re:        re,
			label:     r.Label,
			destField: dest,
		})
	}
	return &Classifier{rules: compiled}, nil
}

// Apply evaluates the rules against rec. The first matching rule writes its
// label into the destination field. If no rule matches the record is unchanged.
func (c *Classifier) Apply(rec map[string]any) map[string]any {
	for _, r := range c.rules {
		val, ok := rec[r.field]
		if !ok {
			continue
		}
		s, ok := val.(string)
		if !ok {
			continue
		}
		if r.re.MatchString(s) {
			rec[r.destField] = r.label
			return rec
		}
	}
	return rec
}
