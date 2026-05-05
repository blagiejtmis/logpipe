// Package sink provides log routing destinations for logpipe.
package sink

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Entry represents a structured log entry routed to a sink.
type Entry struct {
	Source    string    `json:"source"`
	Line      string    `json:"line"`
	Timestamp time.Time `json:"timestamp"`
}

// Sink is the interface that all log sinks must implement.
type Sink interface {
	Write(entry Entry) error
	Close() error
}

// StdoutSink writes structured log entries to an io.Writer (defaults to os.Stdout).
type StdoutSink struct {
	out    io.Writer
	format string // "json" or "text"
}

// NewStdoutSink creates a new StdoutSink. format must be "json" or "text".
func NewStdoutSink(format string) (*StdoutSink, error) {
	if format != "json" && format != "text" {
		return nil, fmt.Errorf("sink/stdout: unsupported format %q, must be \"json\" or \"text\"", format)
	}
	return &StdoutSink{
		out:    os.Stdout,
		format: format,
	}, nil
}

// newStdoutSinkWriter is used in tests to inject a custom writer.
func newStdoutSinkWriter(w io.Writer, format string) *StdoutSink {
	return &StdoutSink{out: w, format: format}
}

// Write encodes and writes a log entry to the sink's output.
func (s *StdoutSink) Write(entry Entry) error {
	switch s.format {
	case "json":
		data, err := json.Marshal(entry)
		if err != nil {
			return fmt.Errorf("sink/stdout: marshal error: %w", err)
		}
		_, err = fmt.Fprintln(s.out, string(data))
		return err
	case "text":
		_, err := fmt.Fprintf(s.out, "[%s] %s: %s\n",
			entry.Timestamp.Format(time.RFC3339), entry.Source, entry.Line)
		return err
	}
	return nil
}

// Close is a no-op for StdoutSink but satisfies the Sink interface.
func (s *StdoutSink) Close() error { return nil }
