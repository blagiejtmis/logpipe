package ratelimit

import (
	"fmt"

	"github.com/logpipe/logpipe/internal/config"
)

// Manager holds a rate limiter per source, keyed by source path.
type Manager struct {
	limiters map[string]*Limiter
	default_ *Limiter
}

// NewManager creates a Manager from the provided config.
// If a global rate limit is configured it is used as the default for any
// source that does not have its own entry. Sources with no limit configured
// and no global default are allowed through unconditionally.
func NewManager(cfg *config.Config) (*Manager, error) {
	m := &Manager{
		limiters: make(map[string]*Limiter),
	}

	if cfg.RateLimit.Default.Rate > 0 {
		l, err := New(cfg.RateLimit.Default.Rate, cfg.RateLimit.Default.WindowSecs)
		if err != nil {
			return nil, fmt.Errorf("ratelimit: default: %w", err)
		}
		m.default_ = l
	}

	for _, sr := range cfg.RateLimit.Sources {
		l, err := New(sr.Rate, sr.WindowSecs)
		if err != nil {
			return nil, fmt.Errorf("ratelimit: source %q: %w", sr.Source, err)
		}
		m.limiters[sr.Source] = l
	}

	return m, nil
}

// Allow reports whether the line from the given source should be forwarded.
// It uses the source-specific limiter when available, falling back to the
// default limiter. If neither exists the line is always allowed.
func (m *Manager) Allow(source string) bool {
	if l, ok := m.limiters[source]; ok {
		return l.Allow(source)
	}
	if m.default_ != nil {
		return m.default_.Allow(source)
	}
	return true
}
