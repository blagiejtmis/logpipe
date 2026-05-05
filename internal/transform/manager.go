package transform

import (
	"fmt"

	"github.com/user/logpipe/internal/config"
)

// Manager holds a Transformer per named source, plus an optional global
// transformer that is applied to every record regardless of source.
type Manager struct {
	global  *Transformer
	perSource map[string]*Transformer
}

// NewManager builds a Manager from the application config.
// config.Transform entries with Source == "" are treated as global rules.
func NewManager(cfg *config.Config) (*Manager, error) {
	m := &Manager{
		perSource: make(map[string]*Transformer),
	}
	var globalRules []Rule
	sourceRules := make(map[string][]Rule)

	for _, tc := range cfg.Transforms {
		rules := make([]Rule, len(tc.Rules))
		for i, r := range tc.Rules {
			rules[i] = Rule{
				Op:      r.Op,
				Field:   r.Field,
				Value:   r.Value,
				NewName: r.NewName,
				Pattern: r.Pattern,
			}
		}
		if tc.Source == "" {
			globalRules = append(globalRules, rules...)
		} else {
			sourceRules[tc.Source] = append(sourceRules[tc.Source], rules...)
		}
	}

	if len(globalRules) > 0 {
		tr, err := New(globalRules)
		if err != nil {
			return nil, fmt.Errorf("transform manager global rules: %w", err)
		}
		m.global = tr
	}

	for src, rules := range sourceRules {
		tr, err := New(rules)
		if err != nil {
			return nil, fmt.Errorf("transform manager source %q: %w", src, err)
		}
		m.perSource[src] = tr
	}
	return m, nil
}

// Apply runs the global transformer (if any) followed by the source-specific
// transformer (if any) on the given record, returning the transformed copy.
func (m *Manager) Apply(source string, record map[string]string) map[string]string {
	out := record
	if m.global != nil {
		out = m.global.Apply(out)
	}
	if tr, ok := m.perSource[source]; ok {
		out = tr.Apply(out)
	}
	return out
}
