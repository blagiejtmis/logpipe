package circuitbreaker

import (
	"fmt"
	"sync"
	"time"
)

// Config holds circuit-breaker settings for a single sink or a default.
type Config struct {
	MaxFailures int           `yaml:"max_failures"`
	Cooldown    time.Duration `yaml:"cooldown"`
}

// ManagerConfig is the top-level config section for circuit breakers.
type ManagerConfig struct {
	Default  *Config            `yaml:"default"`
	PerSink  map[string]*Config `yaml:"per_sink"`
}

const (
	defaultMaxFailures = 5
	defaultCooldown    = 30 * time.Second
)

// Manager holds a Breaker per sink, created lazily.
type Manager struct {
	mu      sync.Mutex
	cfg     *ManagerConfig
	breakers map[string]*Breaker
}

// NewManager creates a Manager from cfg. A nil cfg uses package defaults.
func NewManager(cfg *ManagerConfig) (*Manager, error) {
	if cfg == nil {
		cfg = &ManagerConfig{}
	}
	m := &Manager{
		cfg:      cfg,
		breakers: make(map[string]*Breaker),
	}
	// Eagerly validate per-sink configs.
	for sink, sc := range cfg.PerSink {
		if _, err := m.buildBreaker(sc); err != nil {
			return nil, fmt.Errorf("circuitbreaker: sink %q: %w", sink, err)
		}
	}
	return m, nil
}

// For returns the Breaker for the given sink ID, creating it if necessary.
func (m *Manager) For(sinkID string) (*Breaker, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if b, ok := m.breakers[sinkID]; ok {
		return b, nil
	}
	cfg := m.resolveConfig(sinkID)
	b, err := m.buildBreaker(cfg)
	if err != nil {
		return nil, err
	}
	m.breakers[sinkID] = b
	return b, nil
}

func (m *Manager) resolveConfig(sinkID string) *Config {
	if sc, ok := m.cfg.PerSink[sinkID]; ok && sc != nil {
		return sc
	}
	if m.cfg.Default != nil {
		return m.cfg.Default
	}
	return &Config{MaxFailures: defaultMaxFailures, Cooldown: defaultCooldown}
}

func (m *Manager) buildBreaker(cfg *Config) (*Breaker, error) {
	mf := cfg.MaxFailures
	if mf <= 0 {
		mf = defaultMaxFailures
	}
	cd := cfg.Cooldown
	if cd <= 0 {
		cd = defaultCooldown
	}
	return New(mf, cd)
}
