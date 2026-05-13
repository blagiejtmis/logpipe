package topology

import (
	"testing"

	"github.com/logpipe/logpipe/internal/config"
)

func TestNewManager_ValidConfig(t *testing.T) {
	cfg := makeConfig([]string{"a", "b"}, []string{"out"})
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.NodeCount() != 3 {
		t.Fatalf("expected 3 nodes, got %d", m.NodeCount())
	}
}

func TestNewManager_InvalidConfig_ReturnsError(t *testing.T) {
	cfg := makeConfig([]string{"dup", "dup"}, []string{"out"})
	_, err := NewManager(cfg)
	if err == nil {
		t.Fatal("expected error for duplicate node")
	}
}

func TestManager_SourceIDs(t *testing.T) {
	cfg := makeConfig([]string{"s1", "s2"}, []string{"sink"})
	m, _ := NewManager(cfg)
	ids := m.SourceIDs()
	if len(ids) != 2 {
		t.Fatalf("expected 2 source IDs, got %d", len(ids))
	}
}

func TestManager_SinkIDs(t *testing.T) {
	cfg := makeConfig([]string{"src"}, []string{"k1", "k2"})
	m, _ := NewManager(cfg)
	ids := m.SinkIDs()
	if len(ids) != 2 {
		t.Fatalf("expected 2 sink IDs, got %d", len(ids))
	}
}

func TestNewManager_NilConfig_ReturnsError(t *testing.T) {
	_, err := NewManager(nil)
	if err == nil {
		t.Fatal("expected error for nil config")
	}
}

func TestManager_EmptyGraph(t *testing.T) {
	cfg := &config.Config{}
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.NodeCount() != 0 {
		t.Fatalf("expected 0 nodes, got %d", m.NodeCount())
	}
	if len(m.SourceIDs()) != 0 {
		t.Error("expected no source IDs")
	}
	if len(m.SinkIDs()) != 0 {
		t.Error("expected no sink IDs")
	}
}
