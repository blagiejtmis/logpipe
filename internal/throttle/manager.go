package throttle

import (
	"fmt"
	"sync"
)

// Config holds throttle configuration mirroring the top-level config structure.
type Config struct {
	DefaultRatePerSec int            `yaml:"default_rate_per_sec"`
	Sources           map[string]int `yaml:"sources"`
}

// Manager resolves a Throttler for each source, honouring source-specific
// overrides before falling back to the default rate.
type Manager struct {
	mu       sync.Mutex
	cfg      *Config
	cache    map[string]*Throttler
	nilThrot *noopThrottler
}

// noopThrottler is returned when no throttle is configured.
type noopThrottler struct{}

func (n *noopThrottler) Wait() {}

// Waiter is the interface satisfied by both Throttler and noopThrottler.
type Waiter interface {
	Wait()
}

// NewManager builds a Manager from cfg. If cfg is nil, all sources are
// passed through without throttling.
func NewManager(cfg *Config) (*Manager, error) {
	m := &Manager{
		cache:    make(map[string]*Throttler),
		cfg:      cfg,
		nilThrot: &noopThrottler{},
	}
	if cfg == nil {
		return m, nil
	}
	if cfg.DefaultRatePerSec < 0 {
		return nil, fmt.Errorf("throttle: default_rate_per_sec must be >= 0")
	}
	return m, nil
}

// For returns the Waiter for the given source, creating it on first use.
func (m *Manager) For(source string) (Waiter, error) {
	if m.cfg == nil {
		return m.nilThrot, nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if th, ok := m.cache[source]; ok {
		return th, nil
	}

	rate := m.cfg.DefaultRatePerSec
	if r, ok := m.cfg.Sources[source]; ok {
		rate = r
	}
	if rate == 0 {
		return m.nilThrot, nil
	}

	th, err := New(rate)
	if err != nil {
		return nil, fmt.Errorf("throttle: source %q: %w", source, err)
	}
	m.cache[source] = th
	return th, nil
}
