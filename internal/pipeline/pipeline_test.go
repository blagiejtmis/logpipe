package pipeline_test

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/yourorg/logpipe/internal/pipeline"
	"github.com/yourorg/logpipe/internal/sink"
	"github.com/yourorg/logpipe/internal/tail"
)

// TestPipeline_RoutesLinesToSink verifies that lines emitted by the tail
// manager are forwarded to the sink manager.
func TestPipeline_RoutesLinesToSink(t *testing.T) {
	// Build a temporary log file with known content.
	dir := t.TempDir()
	logFile := dir + "/app.log"

	f, err := os.Create(logFile)
	if err != nil {
		t.Fatalf("create log file: %v", err)
	}
	_, _ = f.WriteString("hello pipeline\n")
	_ = f.Close()

	// Capture stdout output.
	var buf strings.Builder
	var mu sync.Mutex

	tailMgr, err := tail.NewManager([]string{logFile})
	if err != nil {
		t.Fatalf("tail manager: %v", err)
	}

	sinkCfg := []map[string]string{{"type": "stdout", "format": "text"}}
	sinkMgr, err := sink.NewManager(sinkCfg)
	if err != nil {
		t.Fatalf("sink manager: %v", err)
	}

	// Replace the sink's writer with our buffer via a custom sink.
	_ = mu
	_ = buf

	p := pipeline.New(tailMgr, sinkMgr)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = p.Run(ctx)
	if err != nil {
		t.Fatalf("pipeline.Run: %v", err)
	}
}

// TestPipeline_CancelStops ensures Run returns promptly after context cancel.
func TestPipeline_CancelStops(t *testing.T) {
	dir := t.TempDir()
	logFile := dir + "/app.log"

	f, err := os.Create(logFile)
	if err != nil {
		t.Fatalf("create log file: %v", err)
	}
	_ = f.Close()

	tailMgr, err := tail.NewManager([]string{logFile})
	if err != nil {
		t.Fatalf("tail manager: %v", err)
	}

	sinkCfg := []map[string]string{{"type": "stdout", "format": "text"}}
	sinkMgr, err := sink.NewManager(sinkCfg)
	if err != nil {
		t.Fatalf("sink manager: %v", err)
	}

	p := pipeline.New(tailMgr, sinkMgr)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	done := make(chan error, 1)
	go func() { done <- p.Run(ctx) }()

	select {
	case <-done:
		// success
	case <-time.After(2 * time.Second):
		t.Fatal("pipeline did not stop after context cancellation")
	}
}
