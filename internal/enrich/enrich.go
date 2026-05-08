// Package enrich adds static or derived fields to log records.
package enrich

import "fmt"

// Rule describes a single enrichment operation.
type Rule struct {
	// Field is the target field name to set.
	Field string
	// Value is the static string value to assign.
	Value string
}

// Enricher applies a set of enrichment rules to a log record.
type Enricher struct {
	rules []Rule
}

// New creates an Enricher from the provided rules.
// Returns an error if any rule has an empty field name.
func New(rules []Rule) (*Enricher, error) {
	for i, r := range rules {
		if r.Field == "" {
			return nil, fmt.Errorf("enrich: rule %d has empty field name", i)
		}
	}
	return &Enricher{rules: rules}, nil
}

// Apply copies the record and appends enrichment fields.
// Existing fields are NOT overwritten.
func (e *Enricher) Apply(record map[string]any) map[string]any {
	out := make(map[string]any, len(record)+len(e.rules))
	for k, v := range record {
		out[k] = v
	}
	for _, r := range e.rules {
		if _, exists := out[r.Field]; !exists {
			out[r.Field] = r.Value
		}
	}
	return out
}
