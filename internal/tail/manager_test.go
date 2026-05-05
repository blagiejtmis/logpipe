package tail_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/yourorg/logpipe/internal/tail"
)

func TestManager_MultipleSources(t *testing.T) {
	dir := t.TempDir()

	f1, err := os.CreateTemp(dir, "src1-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f1.Close()

	f2, err := os.CreateTemp(dir, "src2-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f2.Close()

	ctx, cancel := context.WithCancel(context.Background())

	mgr := tail.NewManager(20)
	mgr.Add(ctx, f1.Name())
	mgr.Add(ctx, f2.Name())

	time.Sleep(50 * time.Millisecond)

	writeLine(t, f1, "from-source-1")
	writeLine(t, f2, "from-source-2")

	received := make(map[string]bool)
	timeout := time.After(2 * time.Second)
	for len(received) < 2 {
		select {
		case line := <-mgr.Output():
			received[line.Source] = true
		case <-timeout:
			t.Fatalf("timeout; only received from: %v", received)
		}
	}

	cancel()
	errs := mgr.Wait()
	if len(errs) != 0 {
		t.Errorf("unexpected errors: %v", errs)
	}
}

func TestManager_InvalidSource(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mgr := tail.NewManager(4)
	mgr.Add(ctx, "/no/such/file.log")

	errs := mgr.Wait()
	if len(errs) == 0 {
		t.Fatal("expected error for invalid source")
	}
}
