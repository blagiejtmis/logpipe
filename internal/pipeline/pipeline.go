// Package pipeline wires tail sources to sink destinations.
package pipeline

import (
	"context"
	"log"
	"strings"

	"github.com/yourorg/logpipe/internal/filter"
	"github.com/yourorg/logpipe/internal/sink"
	"github.com/yourorg/logpipe/internal/tail"
)

// Pipeline reads lines from a tail manager and writes them to a sink manager.
type Pipeline struct {
	tails  *tail.Manager
	sinks  *sink.Manager
	filter *filter.Filter
}

// New creates a Pipeline. If rules is nil or empty no filtering is applied.
func New(tails *tail.Manager, sinks *sink.Manager, rules []filter.Rule) (*Pipeline, error) {
	f, err := filter.New(rules)
	if err != nil {
		return nil, err
	}
	return &Pipeline{tails: tails, sinks: sinks, filter: f}, nil
}

// Run starts the pipeline and blocks until ctx is cancelled.
func (p *Pipeline) Run(ctx context.Context) {
	lines := p.tails.Lines()
	for {
		select {
		case <-ctx.Done():
			return
		case line, ok := <-lines:
			if !ok {
				return
			}
			if !p.filter.NoRules() {
				fields := extractFields(line)
				if !p.filter.Match(fields) {
					continue
				}
			}
			if err := p.sinks.Write(line); err != nil {
				log.Printf("pipeline: write error: %v", err)
			}
		}
	}
}

// extractFields builds a minimal field map from a raw log line.
// It always sets "_raw"; for lines that look like "key=value" pairs it also
// populates individual fields.
func extractFields(line string) map[string]string {
	fields := map[string]string{"_raw": line}
	for _, token := range strings.Fields(line) {
		if kv := strings.SplitN(token, "=", 2); len(kv) == 2 {
			fields[kv[0]] = strings.Trim(kv[1], `"`)
		}
	}
	return fields
}
