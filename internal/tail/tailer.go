package tail

import (
	"bufio"
	"context"
	"io"
	"os"
	"time"
)

// Line represents a single tailed log line with its source.
type Line struct {
	Source string
	Text   string
	Time   time.Time
}

// Tailer tails a file and emits lines on a channel.
type Tailer struct {
	path   string
	output chan<- Line
}

// New creates a new Tailer for the given file path.
func New(path string, output chan<- Line) *Tailer {
	return &Tailer{path: path, output: output}
}

// Run starts tailing the file until ctx is cancelled.
func (t *Tailer) Run(ctx context.Context) error {
	f, err := os.Open(t.path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Seek to end so we only tail new lines.
	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		return err
	}

	reader := bufio.NewReader(f)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				time.Sleep(200 * time.Millisecond)
				continue
			}
			return err
		}

		if len(line) > 0 {
			t.output <- Line{
				Source: t.path,
				Text:   line,
				Time:   time.Now(),
			}
		}
	}
}
