// Package output provides a unified writer that fans records through the
// processing pipeline (filter → transform → sink) for a single source.
package output

import (
	"context"
	"fmt"

	"github.com/yourorg/logpipe/internal/sink"
)

// Record is the minimal log record type shared across packages.
type Record = map[string]any

// Filterer decides whether a record should be forwarded.
type Filterer interface {
	Allow(source string, r Record) (bool, error)
}

// Transformer mutates a record, returning the (possibly new) record.
type Transformer interface {
	Apply(source string, r Record) Record
}

// Writer fans a record through filter → transform → sink.
type Writer struct {
	filter    Filterer
	transform Transformer
	sink      *sink.Manager
}

// New constructs a Writer. filter and transform may be nil (no-op).
func New(f Filterer, t Transformer, s *sink.Manager) (*Writer, error) {
	if s == nil {
		return nil, fmt.Errorf("output: sink manager must not be nil")
	}
	return &Writer{filter: f, transform: t, sink: s}, nil
}

// Write processes a single record from the named source.
// Returns (false, nil) when the record is dropped by the filter.
func (w *Writer) Write(ctx context.Context, source string, r Record) (bool, error) {
	if w.filter != nil {
		ok, err := w.filter.Allow(source, r)
		if err != nil {
			return false, fmt.Errorf("output filter [%s]: %w", source, err)
		}
		if !ok {
			return false, nil
		}
	}

	if w.transform != nil {
		r = w.transform.Apply(source, r)
	}

	if err := w.sink.Write(ctx, source, r); err != nil {
		return false, fmt.Errorf("output sink [%s]: %w", source, err)
	}
	return true, nil
}
