package enrich

import (
	"fmt"

	"github.com/yourorg/logpipe/internal/config"
)

// Manager holds per-source and global Enrichers derived from configuration.
type Manager struct {
	global *Enricher
	sources map[string]*Enricher
}

// NewManager builds a Manager from the enrichment section of the config.
// Returns an error if any rule set is invalid.
func NewManager(cfg *config.EnrichConfig) (*Manager, error) {
	if cfg == nil {
		return &Manager{sources: make(map[string]*Enricher)}, nil
	}

	var global *Enricher
	if len(cfg.Global) > 0 {
		rules := toRules(cfg.Global)
		var err error
		global, err = New(rules)
		if err != nil {
			return nil, fmt.Errorf("enrich manager: global: %w", err)
		}
	}

	sources := make(map[string]*Enricher, len(cfg.Sources))
	for src, fields := range cfg.Sources {
		rules := toRules(fields)
		e, err := New(rules)
		if err != nil {
			return nil, fmt.Errorf("enrich manager: source %q: %w", src, err)
		}
		sources[src] = e
	}

	return &Manager{global: global, sources: sources}, nil
}

// Apply enriches record for the given source. Global rules are applied first,
// then source-specific rules (which may overwrite global values).
func (m *Manager) Apply(source string, record map[string]any) map[string]any {
	if m.global != nil {
		m.global.Apply(record)
	}
	if e, ok := m.sources[source]; ok {
		e.Apply(record)
	}
	return record
}

// toRules converts a map of field→value pairs into a slice of Rule.
func toRules(fields map[string]string) []Rule {
	rules := make([]Rule, 0, len(fields))
	for k, v := range fields {
		rules = append(rules, Rule{Field: k, Value: v})
	}
	return rules
}
