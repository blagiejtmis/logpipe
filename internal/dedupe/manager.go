package dedupe

import (
	"fmt"
	"time"

	"github.com/your-org/logpipe/internal/config"
)

// Manager holds per-source (and optional global) Deduplicators.
type Manager struct {
	global  *Deduplicator
	sources map[string]*Deduplicator
}

// NewManager builds a Manager from the dedupe section of the config.
func NewManager(cfg config.DedupeConfig) (*Manager, error) {
	m := &Manager{sources: make(map[string]*Deduplicator)}

	if cfg.Global != nil {
		d, err := New(cfg.Global.Fields, time.Duration(cfg.Global.WindowSecs)*time.Second)
		if err != nil {
			return nil, fmt.Errorf("dedupe global: %w", err)
		}
		m.global = d
	}

	for src, rule := range cfg.Sources {
		d, err := New(rule.Fields, time.Duration(rule.WindowSecs)*time.Second)
		if err != nil {
			return nil, fmt.Errorf("dedupe source %q: %w", src, err)
		}
		m.sources[src] = d
	}
	return m, nil
}

// Allow returns true when the record should pass through for the given source.
// Source-specific rules take precedence over the global rule. If no rule
// applies, the record is always allowed.
func (m *Manager) Allow(source string, rec Record) bool {
	if d, ok := m.sources[source]; ok {
		return d.Allow(rec)
	}
	if m.global != nil {
		return m.global.Allow(rec)
	}
	return true
}

// PurgeAll evicts expired fingerprints from every deduplicator.
func (m *Manager) PurgeAll() {
	if m.global != nil {
		m.global.Purge()
	}
	for _, d := range m.sources {
		d.Purge()
	}
}
