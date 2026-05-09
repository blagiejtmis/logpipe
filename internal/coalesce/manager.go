package coalesce

import (
	"fmt"
)

// RuleConfig is the config-layer representation of a coalesce rule.
type RuleConfig struct {
	Sources     []string `yaml:"sources"`
	Dest        string   `yaml:"dest"`
	KeepSources bool     `yaml:"keep_sources"`
}

// Config holds global and per-source coalesce configuration.
type Config struct {
	Global []RuleConfig            `yaml:"global"`
	Sources map[string][]RuleConfig `yaml:"sources"`
}

// Manager resolves a Coalescer per log source.
type Manager struct {
	global  *Coalescer
	sources map[string]*Coalescer
}

// NewManager builds a Manager from cfg. Returns an error if any rule is invalid.
func NewManager(cfg *Config) (*Manager, error) {
	m := &Manager{sources: make(map[string]*Coalescer)}
	if cfg == nil {
		return m, nil
	}
	if len(cfg.Global) > 0 {
		c, err := buildCoalescer(cfg.Global)
		if err != nil {
			return nil, fmt.Errorf("coalesce global: %w", err)
		}
		m.global = c
	}
	for src, rules := range cfg.Sources {
		c, err := buildCoalescer(rules)
		if err != nil {
			return nil, fmt.Errorf("coalesce source %q: %w", src, err)
		}
		m.sources[src] = c
	}
	return m, nil
}

// Get returns the Coalescer for the given source, falling back to the global
// one, or nil if none is configured.
func (m *Manager) Get(source string) *Coalescer {
	if c, ok := m.sources[source]; ok {
		return c
	}
	return m.global
}

func buildCoalescer(cfgRules []RuleConfig) (*Coalescer, error) {
	rules := make([]Rule, len(cfgRules))
	for i, r := range cfgRules {
		rules[i] = Rule{
			Sources:     r.Sources,
			Dest:        r.Dest,
			KeepSources: r.KeepSources,
		}
	}
	return New(rules)
}
