package normalize

import "fmt"

// RuleConfig is the serialisable form of a single normalization rule.
type RuleConfig struct {
	Field string `yaml:"field"`
	Op    string `yaml:"op"`
}

// Config holds global and per-source normalization configuration.
type Config struct {
	Global []RuleConfig            `yaml:"global"`
	Sources map[string][]RuleConfig `yaml:"sources"`
}

// Manager holds per-source and global Normalizers.
type Manager struct {
	global  *Normalizer
	sources map[string]*Normalizer
}

// NewManager builds a Manager from cfg.
// Returns an error if any rule set is invalid.
func NewManager(cfg *Config) (*Manager, error) {
	m := &Manager{sources: make(map[string]*Normalizer)}
	if cfg == nil {
		return m, nil
	}
	if len(cfg.Global) > 0 {
		n, err := buildNormalizer(cfg.Global)
		if err != nil {
			return nil, fmt.Errorf("normalize manager: global: %w", err)
		}
		m.global = n
	}
	for src, rules := range cfg.Sources {
		n, err := buildNormalizer(rules)
		if err != nil {
			return nil, fmt.Errorf("normalize manager: source %q: %w", src, err)
		}
		m.sources[src] = n
	}
	return m, nil
}

// Apply normalizes rec for the given source.
// Source-specific rules are applied after global rules.
func (m *Manager) Apply(source string, rec map[string]any) map[string]any {
	if m.global != nil {
		rec = m.global.Apply(rec)
	}
	if n, ok := m.sources[source]; ok {
		rec = n.Apply(rec)
	}
	return rec
}

func buildNormalizer(cfgs []RuleConfig) (*Normalizer, error) {
	rules := make([]Rule, len(cfgs))
	for i, c := range cfgs {
		rules[i] = Rule{Field: c.Field, Op: Op(c.Op)}
	}
	return New(rules)
}
