package metrics_test

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/metrics"
)

func newBufLogger(buf *bytes.Buffer) *slog.Logger {
	return slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
}

func TestReporter_EmitsOnCancel(t *testing.T) {
	reg := metrics.NewRegistry()
	reg.Counter("lines_in").Add(42)

	var buf bytes.Buffer
	logger := newBufLogger(&buf)

	reporter := metrics.NewReporter(reg, 10*time.Second, logger)
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		reporter.Run(ctx)
		close(done)
	}()

	cancel()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("reporter did not stop after context cancel")
	}

	output := buf.String()
	if !strings.Contains(output, "metrics snapshot") {
		t.Fatalf("expected 'metrics snapshot' in output, got: %s", output)
	}
	if !strings.Contains(output, "lines_in") {
		t.Fatalf("expected 'lines_in' in output, got: %s", output)
	}
}

func TestReporter_EmitsOnTick(t *testing.T) {
	reg := metrics.NewRegistry()
	reg.Counter("events").Inc()

	var buf bytes.Buffer
	logger := newBufLogger(&buf)

	reporter := metrics.NewReporter(reg, 50*time.Millisecond, logger)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	reporter.Run(ctx)

	output := buf.String()
	if !strings.Contains(output, "events") {
		t.Fatalf("expected 'events' counter in output, got: %s", output)
	}
}

func TestReporter_NilLogger_UsesDefault(t *testing.T) {
	reg := metrics.NewRegistry()
	// should not panic
	r := metrics.NewReporter(reg, time.Second, nil)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	r.Run(ctx) // immediate return, empty registry — no output expected
}
