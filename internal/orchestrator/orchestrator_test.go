package orchestrator_test

import (
	"context"
	"testing"
	"time"

	"github.com/example/logpipe/internal/orchestrator"
	"github.com/example/logpipe/internal/config"
)

// stubSink satisfies the sink.Writer interface for testing.
type stubSink struct {
	writes []map[string]string
	closed bool
}

func (s *stubSink) Write(record map[string]string) error {
	s.writes = append(s.writes, record)
	return nil
}

func (s *stubSink) Close() error {
	s.closed = true
	return nil
}

func TestOrchestrator_StartsAndStops(t *testing.T) {
	cfg := &config.Config{
		Sources: []config.Source{
			{Path: "/dev/null"},
		},
		Sinks: []config.Sink{
			{Type: "stdout", Format: "text"},
		},
	}

	o, err := orchestrator.New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- o.Run(ctx)
	}()

	select {
	case err := <-done:
		if err != nil && err != context.DeadlineExceeded {
			t.Errorf("Run returned unexpected error: %v", err)
		}
	case <-time.After(500 * time.Millisecond):
		t.Error("orchestrator did not stop within timeout")
	}
}

func TestOrchestrator_InvalidConfig_ReturnsError(t *testing.T) {
	cfg := &config.Config{
		Sources: []config.Source{},
		Sinks:   []config.Sink{},
	}

	_, err := orchestrator.New(cfg)
	if err == nil {
		t.Fatal("expected error for empty config, got nil")
	}
}
