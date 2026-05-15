package jitter

import (
	"fmt"
	"time"
)

// Config holds jitter configuration loaded from the top-level config file.
type Config struct {
	// Default applies to every source that has no specific override.
	Default *Rule `yaml:"default"`
	// Sources maps source name patterns to per-source rules.
	Sources map[string]*Rule `yaml:"sources"`
}

// Rule is a single jitter rule with a min and max delay expressed as strings.
type Rule struct {
	Min string `yaml:"min"`
	Max string `yaml:"max"`
}

// Manager resolves a Jitter instance for a given source.
type Manager struct {
	defaultJitter *Jitter
	sources       map[string]*Jitter
}

// NewManager builds a Manager from cfg. If cfg is nil every source gets a
// no-op (zero delay) Jitter so callers can always call Wait() unconditionally.
func NewManager(cfg *Config) (*Manager, error) {
	noop, _ := New(0, 0)
	m := &Manager{
		defaultJitter: noop,
		sources:       make(map[string]*Jitter),
	}
	if cfg == nil {
		return m, nil
	}
	if cfg.Default != nil {
		j, err := parseRule(cfg.Default)
		if err != nil {
			return nil, fmt.Errorf("jitter: default rule: %w", err)
		}
		m.defaultJitter = j
	}
	for src, r := range cfg.Sources {
		j, err := parseRule(r)
		if err != nil {
			return nil, fmt.Errorf("jitter: source %q: %w", src, err)
		}
		m.sources[src] = j
	}
	return m, nil
}

// For returns the Jitter for the given source, falling back to the default.
func (m *Manager) For(source string) *Jitter {
	if j, ok := m.sources[source]; ok {
		return j
	}
	return m.defaultJitter
}

func parseRule(r *Rule) (*Jitter, error) {
	var minD, maxD time.Duration
	var err error
	if r.Min != "" {
		minD, err = time.ParseDuration(r.Min)
		if err != nil {
			return nil, fmt.Errorf("invalid min %q: %w", r.Min, err)
		}
	}
	if r.Max != "" {
		maxD, err = time.ParseDuration(r.Max)
		if err != nil {
			return nil, fmt.Errorf("invalid max %q: %w", r.Max, err)
		}
	}
	return New(minD, maxD)
}
