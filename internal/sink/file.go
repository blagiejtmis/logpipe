package sink

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// FileSink writes log entries to a file on disk.
type FileSink struct {
	mu     sync.Mutex
	f      *os.File
	format string
}

// NewFileSink opens (or creates) the file at path and returns a FileSink.
// format must be "json" or "text".
func NewFileSink(path, format string) (*FileSink, error) {
	if format != "json" && format != "text" {
		return nil, fmt.Errorf("sink/file: unsupported format %q (want json or text)", format)
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("sink/file: open %s: %w", path, err)
	}
	return &FileSink{f: f, format: format}, nil
}

// Write appends a log entry to the file.
func (s *FileSink) Write(entry map[string]any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var line string
	switch s.format {
	case "json":
		b, err := json.Marshal(entry)
		if err != nil {
			return fmt.Errorf("sink/file: marshal: %w", err)
		}
		line = string(b) + "\n"
	case "text":
		ts, _ := entry["time"].(string)
		if ts == "" {
			ts = time.Now().UTC().Format(time.RFC3339)
		}
		lvl, _ := entry["level"].(string)
		msg, _ := entry["message"].(string)
		line = fmt.Sprintf("%s\t%s\t%s\n", ts, lvl, msg)
	}

	_, err := fmt.Fprint(s.f, line)
	return err
}

// Close flushes and closes the underlying file.
func (s *FileSink) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.f.Close()
}
