// Package flatten provides a processor that flattens nested map fields
// in a log record into dot-notation top-level keys.
package flatten

import (
	"fmt"
	"strings"
)

// Rule describes a single flatten operation.
type Rule struct {
	// Field is the top-level field whose value (a nested map) will be flattened.
	Field string
	// Prefix, if non-empty, is prepended to every expanded key.
	Prefix string
	// Separator is placed between path segments. Defaults to ".".
	Separator string
	// DropSource removes the original nested field after expansion.
	DropSource bool
}

// Flattener expands nested map fields into dot-notation keys.
type Flattener struct {
	rules []Rule
}

// New creates a Flattener from the provided rules.
// Returns an error if any rule has a blank Field or an empty Separator.
func New(rules []Rule) (*Flattener, error) {
	if len(rules) == 0 {
		return nil, fmt.Errorf("flatten: at least one rule is required")
	}
	for i, r := range rules {
		if strings.TrimSpace(r.Field) == "" {
			return nil, fmt.Errorf("flatten: rule[%d]: field must not be blank", i)
		}
		if r.Separator == "" {
			rules[i].Separator = "."
		}
	}
	return &Flattener{rules: rules}, nil
}

// Apply expands nested map fields in rec according to the configured rules.
// rec is modified in place and returned.
func (f *Flattener) Apply(rec map[string]any) map[string]any {
	for _, r := range f.rules {
		v, ok := rec[r.Field]
		if !ok {
			continue
		}
		nested, ok := v.(map[string]any)
		if !ok {
			continue
		}
		expand(rec, nested, r.Prefix, r.Separator)
		if r.DropSource {
			delete(rec, r.Field)
		}
	}
	return rec
}

// expand recursively walks m and writes leaf values into dst using sep-joined paths.
func expand(dst map[string]any, m map[string]any, prefix, sep string) {
	for k, v := range m {
		key := k
		if prefix != "" {
			key = prefix + sep + k
		}
		if nested, ok := v.(map[string]any); ok {
			expand(dst, nested, key, sep)
		} else {
			dst[key] = v
		}
	}
}
