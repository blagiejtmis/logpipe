package debounce

import (
	"fmt"
	"time"
)

// Rule configures debounce behaviour for a source.
type Rule struct {
	Field  string
	Window time.Duration
}

// Config holds global and per-source debounce rules.
type Config struct {
	Default *Rule
	Sources map[string]*Rule
}

// Manager resolves the correct Debouncer for each source, building instances
// lazily from the provided Config.
type Manager struct {
	defaultDebouncer *Debouncer
	sources          map[string]*Debouncer
}

// NewManager constructs a Manager from cfg. A nil cfg means no debouncing.
func NewManager(cfg *Config) (*Manager, error) {
	if cfg == nil {
		return &Manager{}, nil
	}

	m := &Manager{
		sources: make(map[string]*Debouncer),
	}

	if cfg.Default != nil {
		d, err := New(cfg.Default.Field, cfg.Default.Window)
		if err != nil {
			return nil, fmt.Errorf("debounce manager: default rule: %w", err)
		}
		m.defaultDebouncer = d
	}

	for src, rule := range cfg.Sources {
		if rule == nil {
			continue
		}
		d, err := New(rule.Field, rule.Window)
		if err != nil {
			return nil, fmt.Errorf("debounce manager: source %q: %w", src, err)
		}
		m.sources[src] = d
	}

	return m, nil
}

// Allow returns true when the record should be forwarded for the given source.
// If no debouncer is configured for the source, all records are allowed.
func (m *Manager) Allow(source string, record map[string]any) bool {
	if d, ok := m.sources[source]; ok {
		return d.Allow(record)
	}
	if m.defaultDebouncer != nil {
		return m.defaultDebouncer.Allow(record)
	}
	return true
}
