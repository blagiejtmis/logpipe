package sink

import (
	"fmt"

	"github.com/yourusername/logpipe/internal/config"
)

// Sink is the interface that all log sinks must implement.
type Sink interface {
	Write(line string) error
	Close() error
}

// Manager holds and manages multiple sink instances.
type Manager struct {
	sinks []Sink
}

// NewManager creates a new Manager from the provided sink configs.
// It initialises each sink based on its type field.
func NewManager(cfgs []config.SinkConfig) (*Manager, error) {
	sinks := make([]Sink, 0, len(cfgs))

	for _, cfg := range cfgs {
		var s Sink
		var err error

		switch cfg.Type {
		case "stdout":
			s, err = NewStdoutSink(cfg)
		case "file":
			s, err = NewFileSink(cfg)
		default:
			return nil, fmt.Errorf("sink: unknown type %q", cfg.Type)
		}

		if err != nil {
			return nil, fmt.Errorf("sink: failed to create %q sink: %w", cfg.Type, err)
		}

		sinks = append(sinks, s)
	}

	return &Manager{sinks: sinks}, nil
}

// Write fans out a log line to every registered sink.
// All sinks are attempted; the first error encountered is returned.
func (m *Manager) Write(line string) error {
	for _, s := range m.sinks {
		if err := s.Write(line); err != nil {
			return err
		}
	}
	return nil
}

// Close shuts down every registered sink.
func (m *Manager) Close() error {
	for _, s := range m.sinks {
		if err := s.Close(); err != nil {
			return err
		}
	}
	return nil
}
