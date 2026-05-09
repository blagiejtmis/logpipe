// Package splitter provides field-based log record splitting.
// A Splitter fans out a single record into multiple records by
// expanding an array-valued field into one record per element.
package splitter

import (
	"errors"
	"fmt"
)

// Record is a map of string keys to arbitrary values, matching the
// convention used throughout logpipe.
type Record = map[string]any

// Splitter expands an array-valued field in a record into multiple
// records, one per element. All other fields are copied verbatim.
type Splitter struct {
	field  string
	prefix string
}

// Config holds configuration for a Splitter.
type Config struct {
	// Field is the name of the array field to split on (required).
	Field string
	// Prefix, when non-empty, is prepended to the index to form a new
	// field name stored in each child record, e.g. "item_0".
	Prefix string
}

// New returns a Splitter configured by cfg.
func New(cfg Config) (*Splitter, error) {
	if cfg.Field == "" {
		return nil, errors.New("splitter: field must not be empty")
	}
	return &Splitter{field: cfg.Field, prefix: cfg.Prefix}, nil
}

// Split expands rec into one record per element of the array stored
// at the configured field. If the field is absent or is not a slice,
// Split returns the original record unchanged in a single-element
// slice. Each child record receives an "_index" field (or
// "<prefix>_index" when a prefix is configured) indicating its
// position in the original array.
func (s *Splitter) Split(rec Record) []Record {
	v, ok := rec[s.field]
	if !ok {
		return []Record{rec}
	}

	elems, ok := toSlice(v)
	if !ok {
		return []Record{rec}
	}

	indexKey := "_index"
	if s.prefix != "" {
		indexKey = s.prefix + "_index"
	}

	out := make([]Record, 0, len(elems))
	for i, elem := range elems {
		child := make(Record, len(rec)+1)
		for k, val := range rec {
			if k != s.field {
				child[k] = val
			}
		}
		child[s.field] = elem
		child[indexKey] = fmt.Sprintf("%d", i)
		out = append(out, child)
	}
	return out
}

// toSlice attempts to convert v to []any.
func toSlice(v any) ([]any, bool) {
	switch t := v.(type) {
	case []any:
		return t, true
	case []string:
		result := make([]any, len(t))
		for i, s := range t {
			result[i] = s
		}
		return result, true
	}
	return nil, false
}
