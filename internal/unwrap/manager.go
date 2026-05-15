package unwrap

import "fmt"

// Config holds global and per-source unwrap rule sets.
type Config struct {
	Global []Rule
	// Sources maps source name to its rule set; overrides Global for that source.
	Sources map[string][]Rule
}

// Manager resolves the correct Unwrapper for a given source.
type Manager struct {
	global  *Unwrapper
	sources map[string]*Unwrapper
}

// NewManager builds a Manager from cfg. A nil cfg produces a no-op manager.
func NewManager(cfg *Config) (*Manager, error) {
	if cfg == nil {
		return &Manager{sources: map[string]*Unwrapper{}}, nil
	}

	var global *Unwrapper
	if len(cfg.Global) > 0 {
		var err error
		global, err = New(cfg.Global)
		if err != nil {
			return nil, fmt.Errorf("unwrap manager: global: %w", err)
		}
	}

	sources := make(map[string]*Unwrapper, len(cfg.Sources))
	for src, rules := range cfg.Sources {
		u, err := New(rules)
		if err != nil {
			return nil, fmt.Errorf("unwrap manager: source %q: %w", src, err)
		}
		sources[src] = u
	}

	return &Manager{global: global, sources: sources}, nil
}

// For returns the Unwrapper for the given source, or nil if none applies.
func (m *Manager) For(source string) *Unwrapper {
	if u, ok := m.sources[source]; ok {
		return u
	}
	return m.global
}
