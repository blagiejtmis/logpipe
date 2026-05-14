package timeseries

import (
	"fmt"
	"sync"
	"time"
)

// Manager holds a Series per source (or a global default).
type Manager struct {
	mu      sync.RWMutex
	global  *Series
	sources map[string]*Series
}

// Config describes time-series configuration.
type Config struct {
	Field          string            `yaml:"field"`
	Window         time.Duration     `yaml:"window"`
	Buckets        int               `yaml:"buckets"`
	SourceOverride map[string]string `yaml:"source_field_override"`
}

// NewManager builds a Manager from cfg. A nil cfg returns a no-op manager.
func NewManager(cfg *Config) (*Manager, error) {
	if cfg == nil {
		return &Manager{sources: map[string]*Series{}}, nil
	}
	window := cfg.Window
	if window <= 0 {
		window = time.Minute
	}
	buckets := cfg.Buckets
	if buckets < 1 {
		buckets = 6
	}
	var global *Series
	if cfg.Field != "" {
		var err error
		global, err = New(cfg.Field, window, buckets)
		if err != nil {
			return nil, fmt.Errorf("timeseries: global: %w", err)
		}
	}
	sources := make(map[string]*Series)
	for src, field := range cfg.SourceOverride {
		s, err := New(field, window, buckets)
		if err != nil {
			return nil, fmt.Errorf("timeseries: source %q: %w", src, err)
		}
		sources[src] = s
	}
	return &Manager{global: global, sources: sources}, nil
}

// Record routes rec to the appropriate Series based on source.
func (m *Manager) Record(source string, rec map[string]any) {
	if s := m.seriesFor(source); s != nil {
		s.Record(rec)
	}
}

// Counts returns the count snapshot for the given source.
func (m *Manager) Counts(source string) map[string]int64 {
	if s := m.seriesFor(source); s != nil {
		return s.Counts()
	}
	return map[string]int64{}
}

func (m *Manager) seriesFor(source string) *Series {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if s, ok := m.sources[source]; ok {
		return s
	}
	return m.global
}
