package parse

import (
	"fmt"

	"github.com/logpipe/logpipe/internal/config"
)

// Manager holds per-source and default parsers.
type Manager struct {
	defaultParser *Parser
	sourceParsers map[string]*Parser
}

// NewManager constructs a Manager from config.
// A default parser is built from cfg.DefaultFormat (falls back to "json").
// Per-source overrides are built from cfg.Sources.
func NewManager(cfg *config.ParseConfig) (*Manager, error) {
	defaultFormat := "json"
	if cfg != nil && cfg.DefaultFormat != "" {
		defaultFormat = cfg.DefaultFormat
	}

	dp, err := New(defaultFormat)
	if err != nil {
		return nil, fmt.Errorf("parse manager: default parser: %w", err)
	}

	m := &Manager{
		defaultParser: dp,
		sourceParsers: make(map[string]*Parser),
	}

	if cfg == nil {
		return m, nil
	}

	for src, format := range cfg.Sources {
		p, err := New(format)
		if err != nil {
			return nil, fmt.Errorf("parse manager: source %q: %w", src, err)
		}
		m.sourceParsers[src] = p
	}

	return m, nil
}

// ParserFor returns the parser configured for the given source, falling back
// to the default parser when no source-specific override exists.
func (m *Manager) ParserFor(source string) *Parser {
	if p, ok := m.sourceParsers[source]; ok {
		return p
	}
	return m.defaultParser
}
