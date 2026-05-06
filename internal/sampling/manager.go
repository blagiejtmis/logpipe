package sampling

import (
	"fmt"

	"github.com/logpipe/logpipe/internal/config"
)

// Manager holds per-source and global samplers.
type Manager struct {
	global  *Sampler
	sources map[string]*Sampler
}

// NewManager constructs a Manager from config.
// A global default rate may be set; source-specific rates override it.
func NewManager(cfg config.SamplingConfig) (*Manager, error) {
	m := &Manager{
		sources: make(map[string]*Sampler),
	}

	if cfg.DefaultRate != 0 {
		s, err := New(cfg.DefaultRate)
		if err != nil {
			return nil, fmt.Errorf("sampling: invalid default rate: %w", err)
		}
		m.global = s
	}

	for source, rate := range cfg.Sources {
		s, err := New(rate)
		if err != nil {
			return nil, fmt.Errorf("sampling: invalid rate for source %q: %w", source, err)
		}
		m.sources[source] = s
	}

	return m, nil
}

// Allow returns true if the record from the given source should be forwarded.
// Source-specific samplers take precedence over the global sampler.
// If no sampler is configured, all records are allowed.
func (m *Manager) Allow(source string) bool {
	if s, ok := m.sources[source]; ok {
		return s.Allow(source)
	}
	if m.global != nil {
		return m.global.Allow(source)
	}
	return true
}
