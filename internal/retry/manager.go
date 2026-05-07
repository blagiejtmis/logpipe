package retry

import (
	"fmt"
	"time"
)

// SinkConfig holds per-sink or default retry configuration from YAML.
type SinkConfig struct {
	Default  *Config            `yaml:"default"`
	Overrides map[string]Config `yaml:"overrides"`
}

// Manager resolves a Retryer for a given sink name.
type Manager struct {
	defaultRetryer *Retryer
	overrides      map[string]*Retryer
}

// defaultConfig is used when no retry config is provided.
var defaultConfig = Config{
	MaxAttempts:  3,
	InitialDelay: 100 * time.Millisecond,
	MaxDelay:     5 * time.Second,
	Multiplier:   2.0,
}

// NewManager builds a Manager from optional config.
// If cfg is nil, a sensible default policy is applied to all sinks.
func NewManager(cfg *SinkConfig) (*Manager, error) {
	base := defaultConfig
	if cfg != nil && cfg.Default != nil {
		base = *cfg.Default
	}
	def, err := New(base)
	if err != nil {
		return nil, fmt.Errorf("retry manager: default policy: %w", err)
	}

	overrides := make(map[string]*Retryer)
	if cfg != nil {
		for sink, oc := range cfg.Overrides {
			r, err := New(oc)
			if err != nil {
				return nil, fmt.Errorf("retry manager: sink %q: %w", sink, err)
			}
			overrides[sink] = r
		}
	}
	return &Manager{defaultRetryer: def, overrides: overrides}, nil
}

// For returns the Retryer for the given sink name.
func (m *Manager) For(sink string) *Retryer {
	if r, ok := m.overrides[sink]; ok {
		return r
	}
	return m.defaultRetryer
}
