package alert

import (
	"fmt"

	"github.com/logpipe/logpipe/internal/metrics"
)

// Manager holds per-source and global Alerters.
type Manager struct {
	global  *Alerter
	sources map[string]*Alerter
}

// AlertRule mirrors the config shape used by NewManager.
type AlertRule struct {
	Field     string
	Pattern   string
	Threshold int
	Callback  func(field, value string, count int)
}

// ManagerConfig carries global and per-source alert rules.
type ManagerConfig struct {
	Global  []AlertRule
	Sources map[string][]AlertRule
}

// NewManager builds an alert Manager from cfg.
// reg is used to wire metrics counters into each Alerter.
func NewManager(cfg *ManagerConfig, reg *metrics.Registry) (*Manager, error) {
	if cfg == nil {
		return &Manager{sources: make(map[string]*Alerter)}, nil
	}

	m := &Manager{sources: make(map[string]*Alerter)}

	if len(cfg.Global) > 0 {
		a, err := buildAlerter(cfg.Global, reg)
		if err != nil {
			return nil, fmt.Errorf("alert manager: global rules: %w", err)
		}
		m.global = a
	}

	for src, rules := range cfg.Sources {
		a, err := buildAlerter(rules, reg)
		if err != nil {
			return nil, fmt.Errorf("alert manager: source %q: %w", src, err)
		}
		m.sources[src] = a
	}

	return m, nil
}

// AlerterFor returns the Alerter for the given source, falling back to the
// global Alerter when no source-specific one exists. Returns nil when neither
// is configured.
func (m *Manager) AlerterFor(source string) *Alerter {
	if a, ok := m.sources[source]; ok {
		return a
	}
	return m.global
}

func buildAlerter(rules []AlertRule, reg *metrics.Registry) (*Alerter, error) {
	var opts []Option
	for _, r := range rules {
		opts = append(opts, WithRule(r.Field, r.Pattern, r.Threshold, r.Callback))
	}
	if reg != nil {
		opts = append(opts, WithRegistry(reg))
	}
	return New(opts...)
}
