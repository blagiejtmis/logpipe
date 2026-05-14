// Package tag provides a Tagger that adds or overwrites a fixed set of
// string tags on every log record it processes. Tags are expressed as
// key/value pairs and are useful for annotating records with environment,
// region, or deployment metadata before they reach a sink.
package tag

import (
	"errors"
	"fmt"
)

// Rule describes a single tag to apply.
type Rule struct {
	// Field is the record key to set.
	Field string
	// Value is the literal string value to assign.
	Value string
	// Overwrite controls whether an existing field is replaced.
	Overwrite bool
}

// Tagger applies a fixed list of tag rules to a record.
type Tagger struct {
	rules []Rule
}

// New creates a Tagger from the provided rules.
// It returns an error if any rule has a blank Field.
func New(rules []Rule) (*Tagger, error) {
	if len(rules) == 0 {
		return nil, errors.New("tag: at least one rule is required")
	}
	for i, r := range rules {
		if r.Field == "" {
			return nil, fmt.Errorf("tag: rule[%d]: field must not be empty", i)
		}
	}
	return &Tagger{rules: rules}, nil
}

// Apply writes tags onto rec according to the configured rules.
// It returns the (possibly modified) record.
func (t *Tagger) Apply(rec map[string]any) map[string]any {
	for _, r := range t.rules {
		if _, exists := rec[r.Field]; exists && !r.Overwrite {
			continue
		}
		rec[r.Field] = r.Value
	}
	return rec
}
