package metrics

import (
	"context"
	"log/slog"
	"time"
)

// Reporter periodically logs a snapshot of a Registry at the given interval.
type Reporter struct {
	reg      *Registry
	interval time.Duration
	logger   *slog.Logger
}

// NewReporter creates a Reporter that will emit snapshots every interval.
// If logger is nil, slog.Default() is used.
func NewReporter(reg *Registry, interval time.Duration, logger *slog.Logger) *Reporter {
	if logger == nil {
		logger = slog.Default()
	}
	return &Reporter{reg: reg, interval: interval, logger: logger}
}

// Run blocks, logging snapshots until ctx is cancelled.
func (r *Reporter) Run(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			r.emit() // final snapshot on shutdown
			return
		case <-ticker.C:
			r.emit()
		}
	}
}

func (r *Reporter) emit() {
	snap := r.reg.Snapshot()
	if len(snap) == 0 {
		return
	}
	attrs := make([]any, 0, len(snap)*2)
	for k, v := range snap {
		attrs = append(attrs, k, v)
	}
	r.logger.Info("metrics snapshot", attrs...)
}
