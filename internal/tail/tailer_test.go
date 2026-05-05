package tail_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/yourorg/logpipe/internal/tail"
)

func writeLine(t *testing.T, f *os.File, line string) {
	t.Helper()
	if _, err := f.WriteString(line + "\n"); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func TestTailer_EmitsLines(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "tail-*.log")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	defer f.Close()

	output := make(chan tail.Line, 10)
	tr := tail.New(f.Name(), output)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() { _ = tr.Run(ctx) }()

	// Give the tailer time to seek to end.
	time.Sleep(50 * time.Millisecond)

	writeLine(t, f, "hello world")
	writeLine(t, f, "second line")

	for _, want := range []string{"hello world\n", "second line\n"} {
		select {
		case got := <-output:
			if got.Text != want {
				t.Errorf("got %q, want %q", got.Text, want)
			}
			if got.Source != f.Name() {
				t.Errorf("source %q, want %q", got.Source, f.Name())
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("timeout waiting for line %q", want)
		}
	}
}

func TestTailer_FileNotFound(t *testing.T) {
	output := make(chan tail.Line, 1)
	tr := tail.New("/nonexistent/path/file.log", output)
	err := tr.Run(context.Background())
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
