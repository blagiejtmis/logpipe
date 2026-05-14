// Package clone provides deep-copy utilities for log records.
//
// A cloned record is a fully independent copy; mutations to the clone
// do not affect the original and vice-versa.
package clone

import "fmt"

// Record is the canonical log record type used across logpipe.
type Record = map[string]any

// Deep returns a deep copy of the given record.
// Nested maps and slices are recursively cloned.
// Returns an error if an unsupported value type is encountered.
func Deep(r Record) (Record, error) {
	out := make(Record, len(r))
	for k, v := range r {
		cv, err := cloneValue(v)
		if err != nil {
			return nil, fmt.Errorf("clone: field %q: %w", k, err)
		}
		out[k] = cv
	}
	return out, nil
}

// MustDeep is like Deep but panics on error.
func MustDeep(r Record) Record {
	out, err := Deep(r)
	if err != nil {
		panic(err)
	}
	return out
}

func cloneValue(v any) (any, error) {
	switch t := v.(type) {
	case nil:
		return nil, nil
	case bool, int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64, string:
		return t, nil
	case map[string]any:
		return Deep(t)
	case []any:
		cs := make([]any, len(t))
		for i, elem := range t {
			cv, err := cloneValue(elem)
			if err != nil {
				return nil, err
			}
			cs[i] = cv
		}
		return cs, nil
	default:
		return nil, fmt.Errorf("unsupported type %T", v)
	}
}
