package ceiling

import (
	"fmt"
	"sync"
	"time"
)

// Rule configures a ceiling for a specific source or as a default.
type Rule struct {
	Max    int64
	Window time.Duration
}

// Config holds ceiling configuration keyed by source name, plus an optional default.
type Config struct {
	Default *Rule
	Sources map[string]*Rule
}

// Manager holds per-source Ceiling instances and a fallback.
type Manager struct {
	mu       sync.Mutex
	default_ *Ceiling
	sources  map[string]*Ceiling
}

// NewManager constructs a Manager from the provided Config.
// If cfg is nil, all records are allowed (no ceiling enforced).
func NewManager(cfg *Config) (*Manager, error) {
	m := &Manager{
		sources: make(map[string]*Ceiling),
	}
	if cfg == nil {
		return m, nil
	}
	if cfg.Default != nil {
		c, err := New(cfg.Default.Max, cfg.Default.Window)
		if err != nil {
			return nil, fmt.Errorf("ceiling manager: default: %w", err)
		}
		m.default_ = c
	}
	for src, rule := range cfg.Sources {
		if rule == nil {
			continue
		}
		c, err := New(rule.Max, rule.Window)
		if err != nil {
			return nil, fmt.Errorf("ceiling manager: source %q: %w", src, err)
		}
		m.sources[src] = c
	}
	return m, nil
}

// Allow returns true if the record from source should be accepted.
// It uses the source-specific ceiling when available, falling back to the default.
// If no ceiling is configured for the source, it always returns true.
func (m *Manager) Allow(source string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if c, ok := m.sources[source]; ok {
		return c.Allow()
	}
	if m.default_ != nil {
		return m.default_.Allow()
	}
	return true
}
