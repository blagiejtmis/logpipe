package lookup

import (
	"fmt"

	"github.com/yourorg/logpipe/internal/config"
)

// Manager holds per-source and global Lookup appliers.
type Manager struct {
	global  *Lookup
	sources map[string]*Lookup
}

// NewManager constructs a Manager from config. Returns an error if any rule
// set is invalid.
func NewManager(cfg *config.LookupConfig) (*Manager, error) {
	if cfg == nil {
		return &Manager{sources: make(map[string]*Lookup)}, nil
	}

	m := &Manager{sources: make(map[string]*Lookup)}

	if len(cfg.Global) > 0 {
		l, err := New(cfg.Global)
		if err != nil {
			return nil, fmt.Errorf("lookup: global rules: %w", err)
		}
		m.global = l
	}

	for src, rules := range cfg.Sources {
		l, err := New(rules)
		if err != nil {
			return nil, fmt.Errorf("lookup: source %q: %w", src, err)
		}
		m.sources[src] = l
	}

	return m, nil
}

// Apply enriches the record using the most specific lookup rules available.
// Source-specific rules take precedence over global rules. If no rules match
// the source, global rules are applied. Returns the (possibly modified) record.
func (m *Manager) Apply(source string, record map[string]any) (map[string]any, error) {
	if l, ok := m.sources[source]; ok {
		return l.Apply(record)
	}
	if m.global != nil {
		return m.global.Apply(record)
	}
	return record, nil
}
