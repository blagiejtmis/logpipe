package multiline

import (
	"fmt"
	"time"
)

// Config holds the multiline configuration block from the top-level config.
type Config struct {
	Default  *RuleConfig            `yaml:"default"`
	Sources  map[string]*RuleConfig `yaml:"sources"`
}

// RuleConfig is the YAML-serialisable form of a Rule.
type RuleConfig struct {
	StartPattern    string `yaml:"start_pattern"`
	ContinuePattern string `yaml:"continue_pattern"`
	MaxLines        int    `yaml:"max_lines"`
	TimeoutMs       int    `yaml:"timeout_ms"`
	Field           string `yaml:"field"`
}

// Manager hands out Assemblers keyed by source name.
type Manager struct {
	defaultRule *Rule
	sources     map[string]*Rule
}

// NewManager builds a Manager from cfg.
// If cfg is nil every source receives a nil Assembler (pass-through).
func NewManager(cfg *Config) (*Manager, error) {
	m := &Manager{sources: make(map[string]*Rule)}
	if cfg == nil {
		return m, nil
	}
	if cfg.Default != nil {
		r, err := toRule(cfg.Default)
		if err != nil {
			return nil, fmt.Errorf("multiline default: %w", err)
		}
		m.defaultRule = r
	}
	for src, rc := range cfg.Sources {
		r, err := toRule(rc)
		if err != nil {
			return nil, fmt.Errorf("multiline source %q: %w", src, err)
		}
		m.sources[src] = r
	}
	return m, nil
}

// Assembler returns a new Assembler for the given source, or nil if no rule
// applies to that source.
func (m *Manager) Assembler(source string) (*Assembler, error) {
	rule := m.defaultRule
	if r, ok := m.sources[source]; ok {
		rule = r
	}
	if rule == nil {
		return nil, nil //nolint:nilnil
	}
	return New(*rule)
}

func toRule(rc *RuleConfig) (*Rule, error) {
	if rc == nil {
		return nil, nil //nolint:nilnil
	}
	r := &Rule{
		StartPattern:    rc.StartPattern,
		ContinuePattern: rc.ContinuePattern,
		MaxLines:        rc.MaxLines,
		Field:           rc.Field,
	}
	if rc.TimeoutMs > 0 {
		r.Timeout = time.Duration(rc.TimeoutMs) * time.Millisecond
	}
	return r, nil
}
