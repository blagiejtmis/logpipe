package redact

import (
	"fmt"

	"github.com/logpipe/logpipe/internal/config"
)

// Manager holds per-source and global redaction rules.
type Manager struct {
	global  *Redactor
	sources map[string]*Redactor
}

// NewManager constructs a Manager from config.
// If cfg is nil, the manager passes all records through unchanged.
func NewManager(cfg *config.RedactConfig) (*Manager, error) {
	m := &Manager{
		sources: make(map[string]*Redactor),
	}
	if cfg == nil {
		return m, nil
	}

	if len(cfg.Global) > 0 {
		r, err := New(cfg.Global)
		if err != nil {
			return nil, fmt.Errorf("redact: global rules: %w", err)
		}
		m.global = r
	}

	for src, rules := range cfg.Sources {
		r, err := New(rules)
		if err != nil {
			return nil, fmt.Errorf("redact: source %q rules: %w", src, err)
		}
		m.sources[src] = r
	}

	return m, nil
}

// Apply redacts the record according to the rules for the given source,
// falling back to global rules when no source-specific rules exist.
func (m *Manager) Apply(source string, record map[string]any) map[string]any {
	if r, ok := m.sources[source]; ok {
		return r.Apply(record)
	}
	if m.global != nil {
		return m.global.Apply(record)
	}
	return record
}
