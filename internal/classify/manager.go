package classify

import (
	"fmt"
)

// Config holds classify configuration loaded from the main config file.
type Config struct {
	// Default rules applied when no source-specific rules match.
	Default []Rule `yaml:"default"`
	// Sources maps source names to their own rule sets.
	Sources map[string][]Rule `yaml:"sources"`
}

// Manager resolves the correct Classifier for a given source.
type Manager struct {
	defaultClassifier *Classifier
	sources           map[string]*Classifier
}

// NewManager builds a Manager from cfg. A nil cfg produces a no-op manager
// that never classifies records.
func NewManager(cfg *Config) (*Manager, error) {
	m := &Manager{sources: make(map[string]*Classifier)}
	if cfg == nil {
		return m, nil
	}
	if len(cfg.Default) > 0 {
		c, err := New(cfg.Default)
		if err != nil {
			return nil, fmt.Errorf("classify manager: default rules: %w", err)
		}
		m.defaultClassifier = c
	}
	for src, rules := range cfg.Sources {
		c, err := New(rules)
		if err != nil {
			return nil, fmt.Errorf("classify manager: source %q: %w", src, err)
		}
		m.sources[src] = c
	}
	return m, nil
}

// Apply classifies rec using the rules for source. If no source-specific rules
// exist the default classifier is used. If neither is configured rec is
// returned unchanged.
func (m *Manager) Apply(source string, rec map[string]any) map[string]any {
	if c, ok := m.sources[source]; ok {
		return c.Apply(rec)
	}
	if m.defaultClassifier != nil {
		return m.defaultClassifier.Apply(rec)
	}
	return rec
}
