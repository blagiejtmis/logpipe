package sink

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

var fixedTime = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

func TestStdoutSink_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	s := newStdoutSinkWriter(&buf, "json")

	entry := Entry{Source: "app.log", Line: "hello world", Timestamp: fixedTime}
	if err := s.Write(entry); err != nil {
		t.Fatalf("Write() error: %v", err)
	}

	var got Entry
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &got); err != nil {
		t.Fatalf("failed to parse JSON output: %v", err)
	}
	if got.Source != entry.Source {
		t.Errorf("source: got %q, want %q", got.Source, entry.Source)
	}
	if got.Line != entry.Line {
		t.Errorf("line: got %q, want %q", got.Line, entry.Line)
	}
}

func TestStdoutSink_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	s := newStdoutSinkWriter(&buf, "text")

	entry := Entry{Source: "svc.log", Line: "starting up", Timestamp: fixedTime}
	if err := s.Write(entry); err != nil {
		t.Fatalf("Write() error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "svc.log") {
		t.Errorf("output missing source: %q", out)
	}
	if !strings.Contains(out, "starting up") {
		t.Errorf("output missing line: %q", out)
	}
	if !strings.Contains(out, "2024-01-15T12:00:00Z") {
		t.Errorf("output missing timestamp: %q", out)
	}
}

func TestNewStdoutSink_InvalidFormat(t *testing.T) {
	_, err := NewStdoutSink("xml")
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
	if !strings.Contains(err.Error(), "xml") {
		t.Errorf("error should mention bad format: %v", err)
	}
}

func TestStdoutSink_Close(t *testing.T) {
	var buf bytes.Buffer
	s := newStdoutSinkWriter(&buf, "json")
	if err := s.Close(); err != nil {
		t.Errorf("Close() returned unexpected error: %v", err)
	}
}
