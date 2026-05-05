package sink

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/logpipe/internal/config"
)

func TestNewManager_StdoutSink(t *testing.T) {
	cfgs := []config.SinkConfig{
		{Type: "stdout", Format: "text"},
	}

	m, err := NewManager(cfgs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer m.Close()

	if len(m.sinks) != 1 {
		t.Fatalf("expected 1 sink, got %d", len(m.sinks))
	}
}

func TestNewManager_FileSink(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.log")

	cfgs := []config.SinkConfig{
		{Type: "file", Format: "json", Path: path},
	}

	m, err := NewManager(cfgs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer m.Close()

	if len(m.sinks) != 1 {
		t.Fatalf("expected 1 sink, got %d", len(m.sinks))
	}
}

func TestNewManager_UnknownType(t *testing.T) {
	cfgs := []config.SinkConfig{
		{Type: "kafka", Format: "json"},
	}

	_, err := NewManager(cfgs)
	if err == nil {
		t.Fatal("expected error for unknown sink type, got nil")
	}
}

func TestManager_Write_FansOut(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "fan.log")

	cfgs := []config.SinkConfig{
		{Type: "stdout", Format: "text"},
		{Type: "file", Format: "text", Path: path},
	}

	m, err := NewManager(cfgs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer m.Close()

	if err := m.Write("hello fan-out"); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read file sink output: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty file sink output")
	}
}

func TestManager_Close(t *testing.T) {
	cfgs := []config.SinkConfig{
		{Type: "stdout", Format: "json"},
	}

	m, err := NewManager(cfgs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := m.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}
}
