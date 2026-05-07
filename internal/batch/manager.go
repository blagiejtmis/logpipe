package batch

import (
	"fmt"
	"time"
)

// Config holds batching configuration, typically sourced from the top-level
// config file.
type Config struct {
	// MaxSize is the maximum number of records per batch. Defaults to 100.
	MaxSize int `yaml:"max_size"`
	// FlushInterval is the maximum time between flushes. Defaults to 5s.
	FlushInterval time.Duration `yaml:"flush_interval"`
}

// Manager creates per-source Batchers sharing a common default configuration.
type Manager struct {
	defaultMaxSize  int
	defaultInterval time.Duration
}

// NewManager validates cfg and returns a Manager. If cfg is nil, built-in
// defaults are used (maxSize=100, interval=5s).
func NewManager(cfg *Config) (*Manager, error) {
	m := &Manager{
		defaultMaxSize:  100,
		defaultInterval: 5 * time.Second,
	}
	if cfg == nil {
		return m, nil
	}
	if cfg.MaxSize != 0 {
		if cfg.MaxSize < 1 {
			return nil, fmt.Errorf("batch: max_size must be >= 1")
		}
		m.defaultMaxSize = cfg.MaxSize
	}
	if cfg.FlushInterval != 0 {
		if cfg.FlushInterval <= 0 {
			return nil, fmt.Errorf("batch: flush_interval must be > 0")
		}
		m.defaultInterval = cfg.FlushInterval
	}
	return m, nil
}

// NewBatcher creates a Batcher using the manager's default settings.
func (m *Manager) NewBatcher(flush FlushFunc) (*Batcher, error) {
	return New(m.defaultMaxSize, m.defaultInterval, flush)
}
