// Package sequence assigns a monotonically increasing sequence number to
// each log record as it passes through the pipeline. This is useful for
// detecting dropped or out-of-order records downstream.
package sequence

import (
	"fmt"
	"sync/atomic"

	"logpipe/internal/config"
)

// Sequencer stamps records with an incrementing counter.
type Sequencer struct {
	counter uint64
	field   string
	start   uint64
}

// New creates a Sequencer that writes the sequence number into field.
// start is the first value emitted (typically 1).
func New(field string, start uint64) (*Sequencer, error) {
	if field == "" {
		return nil, fmt.Errorf("sequence: field name must not be empty")
	}
	return &Sequencer{
		counter: start - 1,
		field:   field,
		start:   start,
	}, nil
}

// Apply stamps record with the next sequence number and returns it.
func (s *Sequencer) Apply(record map[string]any) map[string]any {
	n := atomic.AddUint64(&s.counter, 1)
	out := make(map[string]any, len(record)+1)
	for k, v := range record {
		out[k] = v
	}
	out[s.field] = n
	return out
}

// Reset resets the counter back to the configured start value.
func (s *Sequencer) Reset() {
	atomic.StoreUint64(&s.counter, s.start-1)
}

// Config holds per-source or global sequencing configuration.
type Config struct {
	Field string `yaml:"field"`
	Start uint64 `yaml:"start"`
}

// applyDefaults fills in zero-value fields with sensible defaults.
func applyDefaults(c *Config) {
	if c.Field == "" {
		c.Field = "seq"
	}
	if c.Start == 0 {
		c.Start = 1
	}
}

// Manager holds sequencers keyed by source name.
type Manager struct {
	global  *Sequencer
	sources map[string]*Sequencer
}

// NewManager builds a Manager from the pipeline config.
func NewManager(cfg *config.SequenceConfig) (*Manager, error) {
	if cfg == nil {
		return &Manager{sources: map[string]*Sequencer{}}, nil
	}

	m := &Manager{sources: make(map[string]*Sequencer, len(cfg.Sources))}

	if cfg.Default != nil {
		applyDefaults(cfg.Default)
		s, err := New(cfg.Default.Field, cfg.Default.Start)
		if err != nil {
			return nil, fmt.Errorf("sequence: default: %w", err)
		}
		m.global = s
	}

	for src, sc := range cfg.Sources {
		applyDefaults(sc)
		s, err := New(sc.Field, sc.Start)
		if err != nil {
			return nil, fmt.Errorf("sequence: source %q: %w", src, err)
		}
		m.sources[src] = s
	}
	return m, nil
}

// Apply stamps the record for the given source.
// If no source-specific sequencer exists the global one is used.
// If neither is configured the record is returned unchanged.
func (m *Manager) Apply(source string, record map[string]any) map[string]any {
	if s, ok := m.sources[source]; ok {
		return s.Apply(record)
	}
	if m.global != nil {
		return m.global.Apply(record)
	}
	return record
}
