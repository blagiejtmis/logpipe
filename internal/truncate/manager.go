package truncate

import (
	"fmt"

	"github.com/logpipe/logpipe/internal/config"
)

// Manager holds per-source and global Truncator instances.
type Manager struct {
	global *Truncator
	sources map[string]*Truncator
}

// NewManager builds a Manager from config. A nil config returns a no-op manager.
func NewManager(cfg *config.TruncateConfig) (*Manager, error) {
	if cfg == nil {
		return &Manager{sources: make(map[string]*Truncator)}, nil
	}

	m := &Manager{sources: make(map[string]*Truncator)}

	if len(cfg.Default) > 0 {
		t, err := New(cfg.Default)
		if err != nil {
			return nil, fmt.Errorf("truncate: default rules: %w", err)
		}
		m.global = t
	}

	for src, rules := range cfg.Sources {
		t, err := New(rules)
		if err != nil {
			return nil, fmt.Errorf("truncate: source %q: %w", src, err)
		}
		m.sources[src] = t
	}

	return m, nil
}

// For returns the Truncator for the given source, falling back to the global
// truncator. Returns nil if no rules apply.
func (m *Manager) For(source string) *Truncator {
	if t, ok := m.sources[source]; ok {
		return t
	}
	return m.global
}
