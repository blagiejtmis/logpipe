// Package orchestrator wires together the tail, filter, transform,
// router, sink, and metrics subsystems into a single runnable unit.
package orchestrator

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/example/logpipe/internal/config"
	"github.com/example/logpipe/internal/filter"
	"github.com/example/logpipe/internal/metrics"
	"github.com/example/logpipe/internal/pipeline"
	"github.com/example/logpipe/internal/router"
	"github.com/example/logpipe/internal/sink"
	"github.com/example/logpipe/internal/tail"
	"github.com/example/logpipe/internal/transform"
)

// Orchestrator owns the lifecycle of all subsystems.
type Orchestrator struct {
	cfg      *config.Config
	logger   *slog.Logger
	reg      *metrics.Registry
	pipeline *pipeline.Pipeline
}

// New builds an Orchestrator from cfg. It constructs every subsystem
// and returns an error if any configuration is invalid.
func New(cfg *config.Config, logger *slog.Logger, reg *metrics.Registry) (*Orchestrator, error) {
	if logger == nil {
		logger = slog.Default()
	}
	if reg == nil {
		reg = metrics.NewRegistry()
	}

	// Sinks
	sinkMgr, err := sink.NewManager(cfg)
	if err != nil {
		return nil, fmt.Errorf("orchestrator: sinks: %w", err)
	}

	// Router
	rtr, err := router.NewFromConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("orchestrator: router: %w", err)
	}

	// Filter
	flt, err := filter.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("orchestrator: filter: %w", err)
	}

	// Transform
	txMgr, err := transform.NewManager(cfg)
	if err != nil {
		return nil, fmt.Errorf("orchestrator: transform: %w", err)
	}

	// Tail manager
	tailMgr, err := tail.NewManager(cfg)
	if err != nil {
		return nil, fmt.Errorf("orchestrator: tail: %w", err)
	}

	// Pipeline
	p := pipeline.New(tailMgr, flt, txMgr, rtr, sinkMgr, reg, logger)

	return &Orchestrator{
		cfg:      cfg,
		logger:   logger,
		reg:      reg,
		pipeline: p,
	}, nil
}

// Run starts all subsystems and blocks until ctx is cancelled.
func (o *Orchestrator) Run(ctx context.Context) error {
	o.logger.Info("orchestrator: starting")
	if err := o.pipeline.Run(ctx); err != nil {
		return fmt.Errorf("orchestrator: pipeline: %w", err)
	}
	o.logger.Info("orchestrator: stopped")
	return nil
}

// Registry returns the metrics registry used by this orchestrator.
func (o *Orchestrator) Registry() *metrics.Registry { return o.reg }
