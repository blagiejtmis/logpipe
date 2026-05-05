package sink

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileSink_JSONFormat(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.log")

	s, err := NewFileSink(path, "json")
	if err != nil {
		t.Fatalf("NewFileSink: %v", err)
	}
	defer s.Close()

	entry := map[string]any{"level": "info", "message": "hello", "time": "2024-01-01T00:00:00Z"}
	if err := s.Write(entry); err != nil {
		t.Fatalf("Write: %v", err)
	}
	_ = s.Close()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal output: %v", err)
	}
	if got["message"] != "hello" {
		t.Errorf("expected message=hello, got %v", got["message"])
	}
}

func TestFileSink_TextFormat(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.log")

	s, err := NewFileSink(path, "text")
	if err != nil {
		t.Fatalf("NewFileSink: %v", err)
	}

	entry := map[string]any{"level": "warn", "message": "disk full", "time": "2024-06-01T12:00:00Z"}
	if err := s.Write(entry); err != nil {
		t.Fatalf("Write: %v", err)
	}
	_ = s.Close()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	line := string(data)
	if !strings.Contains(line, "warn") || !strings.Contains(line, "disk full") {
		t.Errorf("unexpected text output: %q", line)
	}
}

func TestNewFileSink_InvalidFormat(t *testing.T) {
	dir := t.TempDir()
	_, err := NewFileSink(filepath.Join(dir, "x.log"), "xml")
	if err == nil {
		t.Fatal("expected error for invalid format, got nil")
	}
}

func TestNewFileSink_BadPath(t *testing.T) {
	_, err := NewFileSink("/nonexistent/dir/out.log", "json")
	if err == nil {
		t.Fatal("expected error for bad path, got nil")
	}
}
