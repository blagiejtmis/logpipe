// Package lookup provides field-value enrichment via static lookup tables.
// A Lookup matches a record's key field against a map of values and merges
// the corresponding attributes into the record.
package lookup

import (
	"errors"
	"fmt"

	"github.com/yourorg/logpipe/internal/pipeline"
)

// Rule describes a single lookup table configuration.
type Rule struct {
	// KeyField is the record field whose value is used to look up a row.
	KeyField string
	// Table maps key values to a set of fields that will be merged into
	// the record when the key matches.
	Table map[string]map[string]string
	// OnMiss controls behaviour when the key is absent from the table.
	// "ignore" (default) leaves the record unchanged; "drop" discards it.
	OnMiss string
}

// Lookup enriches records by merging fields from a static table.
type Lookup struct {
	rules []Rule
}

// New validates rules and returns a ready-to-use Lookup.
func New(rules []Rule) (*Lookup, error) {
	if len(rules) == 0 {
		return nil, errors.New("lookup: at least one rule is required")
	}
	for i, r := range rules {
		if r.KeyField == "" {
			return nil, fmt.Errorf("lookup: rule[%d]: key_field must not be empty", i)
		}
		if len(r.Table) == 0 {
			return nil, fmt.Errorf("lookup: rule[%d]: table must not be empty", i)
		}
		if r.OnMiss != "" && r.OnMiss != "ignore" && r.OnMiss != "drop" {
			return nil, fmt.Errorf("lookup: rule[%d]: on_miss must be \"ignore\" or \"drop\"", i)
		}
	}
	return &Lookup{rules: rules}, nil
}

// Apply runs all lookup rules against rec.
// It returns the (possibly modified) record and a bool indicating whether
// the record should be kept (false means drop).
func (l *Lookup) Apply(rec pipeline.Record) (pipeline.Record, bool) {
	for _, r := range l.rules {
		keyVal, ok := rec[r.KeyField]
		if !ok {
			if r.OnMiss == "drop" {
				return rec, false
			}
			continue
		}
		key := fmt.Sprintf("%v", keyVal)
		attrs, found := r.Table[key]
		if !found {
			if r.OnMiss == "drop" {
				return rec, false
			}
			continue
		}
		for k, v := range attrs {
			if _, exists := rec[k]; !exists {
				rec[k] = v
			}
		}
	}
	return rec, true
}
