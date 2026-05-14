// Package fingerprint computes stable hashes over log record fields
// to uniquely identify records for deduplication, correlation, or tracing.
package fingerprint

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
)

// Record is the minimal interface a log record must satisfy.
type Record = map[string]any

// Fingerprinter computes a hex fingerprint over a configurable set of fields.
type Fingerprinter struct {
	fields []string // ordered list of fields to include
	sep    string   // separator between field values
}

// Config holds the configuration for a Fingerprinter.
type Config struct {
	// Fields lists the record keys whose values are hashed.
	// If empty, all keys are used (sorted for stability).
	Fields []string `yaml:"fields"`
	// Separator is placed between concatenated values. Defaults to "|".
	Separator string `yaml:"separator"`
}

// New creates a Fingerprinter from cfg.
func New(cfg Config) (*Fingerprinter, error) {
	sep := cfg.Separator
	if sep == "" {
		sep = "|"
	}
	for _, f := range cfg.Fields {
		if strings.TrimSpace(f) == "" {
			return nil, errors.New("fingerprint: field name must not be blank")
		}
	}
	return &Fingerprinter{fields: cfg.Fields, sep: sep}, nil
}

// Compute returns a SHA-256 hex digest derived from the configured fields
// of rec. If no fields are configured, all keys are used in sorted order.
func (fp *Fingerprinter) Compute(rec Record) string {
	keys := fp.fields
	if len(keys) == 0 {
		keys = make([]string, 0, len(rec))
		for k := range rec {
			keys = append(keys, k)
		}
		sort.Strings(keys)
	}

	var sb strings.Builder
	for i, k := range keys {
		if i > 0 {
			sb.WriteString(fp.sep)
		}
		sb.WriteString(fmt.Sprintf("%v", rec[k]))
	}

	sum := sha256.Sum256([]byte(sb.String()))
	return hex.EncodeToString(sum[:])
}

// Equal reports whether two records produce the same fingerprint
// under this Fingerprinter. This is a convenience wrapper around
// Compute that avoids the caller having to compare hex strings manually.
func (fp *Fingerprinter) Equal(a, b Record) bool {
	return fp.Compute(a) == fp.Compute(b)
}
