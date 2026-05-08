package schema

import (
	"fmt"

	"github.com/logpipe/logpipe/internal/config"
)

// Manager holds per-source and global schema validators.
type Manager struct {
	global  *Schema
	sources map[string]*Schema
}

// NewManager constructs a Manager from config.
// If cfg is nil, no validation is performed (all records pass).
func NewManager(cfg *config.SchemaConfig) (*Manager, error) {
	if cfg == nil {
		return &Manager{sources: make(map[string]*Schema)}, nil
	}

	m := &Manager{sources: make(map[string]*Schema)}

	if len(cfg.Global) > 0 {
		s, err := New(cfg.Global)
		if err != nil {
			return nil, fmt.Errorf("schema: global rules: %w", err)
		}
		m.global = s
	}

	for src, rules := range cfg.Sources {
		s, err := New(rules)
		if err != nil {
			return nil, fmt.Errorf("schema: source %q rules: %w", src, err)
		}
		m.sources[src] = s
	}

	return m, nil
}

// Validate checks record against the schema for the given source.
// Source-specific schema takes precedence over global. If no schema
// applies, Validate returns nil.
func (m *Manager) Validate(source string, record map[string]any) error {
	if s, ok := m.sources[source]; ok {
		return s.Validate(record)
	}
	if m.global != nil {
		return m.global.Validate(record)
	}
	return nil
}
