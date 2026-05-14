package aggregate

import (
	"fmt"
	"time"
)

// Config holds aggregation configuration loaded from the top-level config.
type Config struct {
	Default  []RuleConfig            `yaml:"default"`
	Sources  map[string][]RuleConfig `yaml:"sources"`
}

// RuleConfig is the YAML-serialisable form of a Rule.
type RuleConfig struct {
	Field  string `yaml:"field"`
	Op     string `yaml:"op"`
	Window string `yaml:"window"`
}

// Manager holds per-source (or default) Aggregators.
type Manager struct {
	defaultAgg *Aggregator
	sources    map[string]*Aggregator
}

// NewManager builds a Manager from cfg. A nil cfg produces a no-op Manager.
func NewManager(cfg *Config) (*Manager, error) {
	if cfg == nil {
		return &Manager{}, nil
	}
	defAgg, err := buildAggregator(cfg.Default)
	if err != nil {
		return nil, fmt.Errorf("aggregate: default rules: %w", err)
	}
	srcMap := make(map[string]*Aggregator, len(cfg.Sources))
	for src, ruleCfgs := range cfg.Sources {
		agg, err := buildAggregator(ruleCfgs)
		if err != nil {
			return nil, fmt.Errorf("aggregate: source %q: %w", src, err)
		}
		srcMap[src] = agg
	}
	return &Manager{defaultAgg: defAgg, sources: srcMap}, nil
}

// Add routes a record to the appropriate Aggregator.
func (m *Manager) Add(source string, record map[string]any) {
	if agg, ok := m.sources[source]; ok {
		agg.Add(source, record)
		return
	}
	if m.defaultAgg != nil {
		m.defaultAgg.Add(source, record)
	}
}

// Snapshot returns aggregated state for all aggregators keyed by scope
// ("default" or source name).
func (m *Manager) Snapshot() map[string]map[string]map[string]float64 {
	out := make(map[string]map[string]map[string]float64)
	if m.defaultAgg != nil {
		out["default"] = m.defaultAgg.Snapshot()
	}
	for src, agg := range m.sources {
		out[src] = agg.Snapshot()
	}
	return out
}

// Sources returns the list of source names that have explicit aggregation rules
// configured, in no particular order.
func (m *Manager) Sources() []string {
	if len(m.sources) == 0 {
		return nil
	}
	names := make([]string, 0, len(m.sources))
	for src := range m.sources {
		names = append(names, src)
	}
	return names
}

func buildAggregator(cfgs []RuleConfig) (*Aggregator, error) {
	if len(cfgs) == 0 {
		return nil, nil
	}
	rules := make([]Rule, 0, len(cfgs))
	for _, rc := range cfgs {
		d, err := time.ParseDuration(rc.Window)
		if err != nil {
			return nil, fmt.Errorf("invalid window %q: %w", rc.Window, err)
		}
		rules = append(rules, Rule{Field: rc.Field, Op: Op(rc.Op), Window: d})
	}
	return New(rules)
}
