// Package coalesce provides field coalescing: merging multiple source fields
// into a single destination field using the first non-empty value found.
package coalesce

import (
	"errors"
	"fmt"
)

// Rule defines a single coalesce operation.
type Rule struct {
	// Sources is the ordered list of field names to check.
	Sources []string
	// Dest is the field name to write the result into.
	Dest string
	// KeepSources, when true, retains the source fields after coalescing.
	KeepSources bool
}

// Coalescer applies a set of coalesce rules to log records.
type Coalescer struct {
	rules []Rule
}

// New creates a Coalescer from the given rules.
// Returns an error if any rule is misconfigured.
func New(rules []Rule) (*Coalescer, error) {
	for i, r := range rules {
		if len(r.Sources) == 0 {
			return nil, fmt.Errorf("coalesce rule %d: sources must not be empty", i)
		}
		if r.Dest == "" {
			return nil, fmt.Errorf("coalesce rule %d: dest must not be empty", i)
		}
	}
	if len(rules) == 0 {
		return nil, errors.New("coalesce: at least one rule is required")
	}
	return &Coalescer{rules: rules}, nil
}

// Apply runs all coalesce rules against record, mutating it in place.
func (c *Coalescer) Apply(record map[string]any) {
	for _, r := range c.rules {
		var chosen any
		var chosenKey string
		for _, src := range r.Sources {
			v, ok := record[src]
			if !ok {
				continue
			}
			if s, ok := v.(string); ok && s == "" {
				continue
			}
			chosen = v
			chosenKey = src
			break
		}
		if chosen != nil {
			record[r.Dest] = chosen
			if !r.KeepSources {
				for _, src := range r.Sources {
					if src != r.Dest && src != chosenKey || src == chosenKey {
						delete(record, src)
					}
				}
			}
		}
	}
}
