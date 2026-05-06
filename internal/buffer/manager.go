package buffer

import (
	"fmt"

	"github.com/logpipe/logpipe/internal/config"
)

// Manager holds per-source ring buffers and a global fallback.
type Manager struct {
	defaultBuf *RingBuffer
	sourceBufs map[string]*RingBuffer
}

// NewManager constructs a Manager from config.
// Each source may override the global buffer settings.
func NewManager(cfg config.BufferConfig) (*Manager, error) {
	var defaultBuf *RingBuffer
	if cfg.Capacity > 0 {
		buf, err := New(cfg.Capacity, cfg.Policy)
		if err != nil {
			return nil, fmt.Errorf("buffer: default: %w", err)
		}
		defaultBuf = buf
	}

	sourceBufs := make(map[string]*RingBuffer, len(cfg.Sources))
	for src, sc := range cfg.Sources {
		buf, err := New(sc.Capacity, sc.Policy)
		if err != nil {
			return nil, fmt.Errorf("buffer: source %q: %w", src, err)
		}
		sourceBufs[src] = buf
	}

	return &Manager{
		defaultBuf: defaultBuf,
		sourceBufs: sourceBufs,
	}, nil
}

// BufferFor returns the RingBuffer for the given source, falling back to the
// default buffer. Returns nil when no buffer is configured.
func (m *Manager) BufferFor(source string) *RingBuffer {
	if buf, ok := m.sourceBufs[source]; ok {
		return buf
	}
	return m.defaultBuf
}
