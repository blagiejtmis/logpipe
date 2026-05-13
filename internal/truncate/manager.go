package truncate

import "fmt"

// Config holds truncation configuration, mirroring the top-level config shape
// used by other managers in logpipe.
type Config struct {
	// Global rules apply to every source unless overridden.
	Global []Rule
	// PerSource maps a source name to source-specific rules that replace
	// the global rules for that source.
	PerSource map[string][]Rule
}

// Manager resolves a Truncator for a given source, applying per-source
// overrides before falling back to global rules.
type Manager struct {
	global    *Truncator
	perSource map[string]*Truncator
}

// NewManager constructs a Manager from cfg. Returns an error if any rule set
// is invalid. A nil cfg produces a no-op manager.
func NewManager(cfg *Config) (*Manager, error) {
	if cfg == nil {
		return &Manager{}, nil
	}

	var global *Truncator
	if len(cfg.Global) > 0 {
		var err error
		global, err = New(cfg.Global)
		if err != nil {
			return nil, fmt.Errorf("truncate manager: global: %w", err)
		}
	}

	ps := make(map[string]*Truncator, len(cfg.PerSource))
	for src, rules := range cfg.PerSource {
		tr, err := New(rules)
		if err != nil {
			return nil, fmt.Errorf("truncate manager: source %q: %w", src, err)
		}
		ps[src] = tr
	}

	return &Manager{global: global, perSource: ps}, nil
}

// For returns the Truncator for source. Returns nil when no rules apply,
// meaning the caller should pass the record through unchanged.
func (m *Manager) For(source string) *Truncator {
	if t, ok := m.perSource[source]; ok {
		return t
	}
	return m.global
}
