// Package timestamp provides field-level timestamp parsing and normalization
// for log records. It can parse timestamps from string fields using configurable
// layouts and write the result as a normalized RFC3339 value.
package timestamp

import (
	"fmt"
	"time"
)

// Rule describes how to parse and normalize a single timestamp field.
type Rule struct {
	// Field is the record key whose value should be parsed.
	Field string
	// Layouts is an ordered list of time.Parse layouts to try.
	Layouts []string
	// OutputField is where the parsed time is written (defaults to Field).
	OutputField string
	// OutputLayout is the format written to OutputField (defaults to time.RFC3339).
	OutputLayout string
}

// Normalizer parses and normalizes timestamp fields in log records.
type Normalizer struct {
	rules []Rule
}

// New creates a Normalizer from the given rules.
// Returns an error if any rule is misconfigured.
func New(rules []Rule) (*Normalizer, error) {
	if len(rules) == 0 {
		return nil, fmt.Errorf("timestamp: at least one rule is required")
	}
	for i, r := range rules {
		if r.Field == "" {
			return nil, fmt.Errorf("timestamp: rule[%d]: field must not be empty", i)
		}
		if len(r.Layouts) == 0 {
			return nil, fmt.Errorf("timestamp: rule[%d]: at least one layout is required", i)
		}
	}
	return &Normalizer{rules: rules}, nil
}

// Apply parses timestamp fields in rec according to the configured rules.
// Fields that cannot be parsed by any layout are left unchanged.
// rec is modified in place and also returned.
func (n *Normalizer) Apply(rec map[string]any) map[string]any {
	for _, r := range n.rules {
		v, ok := rec[r.Field]
		if !ok {
			continue
		}
		raw, ok := v.(string)
		if !ok {
			continue
		}
		var parsed time.Time
		var found bool
		for _, layout := range r.Layouts {
			t, err := time.Parse(layout, raw)
			if err == nil {
				parsed = t
				found = true
				break
			}
		}
		if !found {
			continue
		}
		outLayout := r.OutputLayout
		if outLayout == "" {
			outLayout = time.RFC3339
		}
		outField := r.OutputField
		if outField == "" {
			outField = r.Field
		}
		rec[outField] = parsed.UTC().Format(outLayout)
	}
	return rec
}
