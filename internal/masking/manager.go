package masking

import (
	"fmt"
)

// RuleConfig is the serialisable form of a masking rule loaded from config.
type RuleConfig struct {
	Field       string `yaml:"field"`
	Pattern     string `yaml:"pattern"`
	Placeholder string `yaml:"placeholder"`
}

// Config holds the masking configuration block.
type Config struct {
	Global []RuleConfig            `yaml:"global"`
	Sources map[string][]RuleConfig `yaml:"sources"`
}

// Manager resolves the correct Masker for a given source.
type Manager struct {
	global  *Masker
	sources map[string]*Masker
}

// NewManager builds a Manager from the supplied Config.
// A nil Config produces a no-op manager.
func NewManager(cfg *Config) (*Manager, error) {
	mgr := &Manager{sources: make(map[string]*Masker)}
	if cfg == nil {
		return mgr, nil
	}
	if len(cfg.Global) > 0 {
		rules, err := toRules(cfg.Global)
		if err != nil {
			return nil, fmt.Errorf("masking: global: %w", err)
		}
		m, err := New(rules)
		if err != nil {
			return nil, fmt.Errorf("masking: global: %w", err)
		}
		mgr.global = m
	}
	for src, cfgRules := range cfg.Sources {
		rules, err := toRules(cfgRules)
		if err != nil {
			return nil, fmt.Errorf("masking: source %q: %w", src, err)
		}
		m, err := New(rules)
		if err != nil {
			return nil, fmt.Errorf("masking: source %q: %w", src, err)
		}
		mgr.sources[src] = m
	}
	return mgr, nil
}

// Apply masks the record using the source-specific masker if present,
// falling back to the global masker. If neither is configured the
// record is returned unmodified.
func (mgr *Manager) Apply(source string, record map[string]any) map[string]any {
	if m, ok := mgr.sources[source]; ok {
		return m.Apply(record)
	}
	if mgr.global != nil {
		return mgr.global.Apply(record)
	}
	return record
}

func toRules(cfgRules []RuleConfig) ([]Rule, error) {
	rules := make([]Rule, len(cfgRules))
	for i, r := range cfgRules {
		rules[i] = Rule{
			Field:       r.Field,
			Pattern:     r.Pattern,
			Placeholder: r.Placeholder,
		}
	}
	return rules, nil
}
