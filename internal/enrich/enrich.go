// Package enrich provides field enrichment for log records,
// allowing static or derived fields to be injected into every record.
package enrich

import (
	"fmt"
	"os"
	"strings"
)

// Rule describes a single enrichment operation.
type Rule struct {
	// Field is the key to add or overwrite in the record.
	Field string
	// Value is the static value to set. Supports the special token
	// "${hostname}" which is resolved at construction time.
	Value string
}

// Enricher applies a set of Rules to log records.
type Enricher struct {
	rules []Rule
}

// New constructs an Enricher from the provided rules.
// Returns an error if any rule has an empty field name.
func New(rules []Rule) (*Enricher, error) {
	resolved := make([]Rule, len(rules))
	for i, r := range rules {
		if strings.TrimSpace(r.Field) == "" {
			return nil, fmt.Errorf("enrich: rule %d has empty field name", i)
		}
		v := r.Value
		if v == "${hostname}" {
			h, err := os.Hostname()
			if err != nil {
				return nil, fmt.Errorf("enrich: resolving hostname: %w", err)
			}
			v = h
		}
		resolved[i] = Rule{Field: r.Field, Value: v}
	}
	return &Enricher{rules: resolved}, nil
}

// Apply adds or overwrites fields in record according to the enricher's rules.
// The record map is modified in place and also returned for convenience.
func (e *Enricher) Apply(record map[string]any) map[string]any {
	for _, r := range e.rules {
		record[r.Field] = r.Value
	}
	return record
}
