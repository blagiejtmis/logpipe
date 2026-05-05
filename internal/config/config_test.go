package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/logpipe/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "logpipe.yaml")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTemp: %v", err)
	}
	return p
}

func TestLoad_Valid(t *testing.T) {
	raw := `
sources:
  - name: app-logs
    type: file
    path: /var/log/app.log
    format: json
sinks:
  - name: stdout-sink
    type: stdout
`
	cfg, err := config.Load(writeTemp(t, raw))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Sources) != 1 {
		t.Errorf("expected 1 source, got %d", len(cfg.Sources))
	}
	if cfg.Sources[0].Name != "app-logs" {
		t.Errorf("unexpected source name: %s", cfg.Sources[0].Name)
	}
	if len(cfg.Sinks) != 1 {
		t.Errorf("expected 1 sink, got %d", len(cfg.Sinks))
	}
}

func TestLoad_MissingSources(t *testing.T) {
	raw := `
sinks:
  - name: stdout-sink
    type: stdout
`
	_, err := config.Load(writeTemp(t, raw))
	if err == nil {
		t.Fatal("expected error for missing sources, got nil")
	}
}

func TestLoad_MissingSinks(t *testing.T) {
	raw := `
sources:
  - name: app-logs
    type: file
`
	_, err := config.Load(writeTemp(t, raw))
	if err == nil {
		t.Fatal("expected error for missing sinks, got nil")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := config.Load("/nonexistent/path/logpipe.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_UnknownField(t *testing.T) {
	raw := `
sources:
  - name: app-logs
    type: file
    unknown_field: oops
sinks:
  - name: out
    type: stdout
`
	_, err := config.Load(writeTemp(t, raw))
	if err == nil {
		t.Fatal("expected error for unknown field, got nil")
	}
}
