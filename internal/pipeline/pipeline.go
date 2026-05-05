// Package pipeline wires tail sources to sink destinations,
// reading log entries from the manager and fanning them out.
package pipeline

import (
	"context"
	"log"

	"github.com/yourorg/logpipe/internal/sink"
	"github.com/yourorg/logpipe/internal/tail"
)

// Entry represents a single log line together with its origin.
type Entry struct {
	Source string
	Line   string
}

// Pipeline connects a tail.Manager to a sink.Manager and pumps entries
// between them until the context is cancelled.
type Pipeline struct {
	tails *tail.Manager
	sinks *sink.Manager
}

// New creates a Pipeline from the supplied managers.
func New(t *tail.Manager, s *sink.Manager) *Pipeline {
	return &Pipeline{tails: t, sinks: s}
}

// Run starts pumping log lines from all tail sources into all sinks.
// It blocks until ctx is cancelled, then drains any remaining lines and
// closes the sink manager.
func (p *Pipeline) Run(ctx context.Context) error {
	lines := p.tails.Lines()

	for {
		select {
		case entry, ok := <-lines:
			if !ok {
				// channel closed — all tailers finished
				return p.sinks.Close()
			}
			if err := p.sinks.Write(entry.Source, entry.Line); err != nil {
				log.Printf("pipeline: write error: %v", err)
			}
		case <-ctx.Done():
			// drain remaining buffered lines before stopping
			for {
				select {
				case entry, ok := <-lines:
					if !ok {
						return p.sinks.Close()
					}
					if err := p.sinks.Write(entry.Source, entry.Line); err != nil {
						log.Printf("pipeline: write error during drain: %v", err)
					}
				default:
					return p.sinks.Close()
				}
			}
		}
	}
}
