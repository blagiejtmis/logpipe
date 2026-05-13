package topology

import (
	"testing"

	"github.com/logpipe/logpipe/internal/config"
)

func makeConfig(sources, sinks []string) *config.Config {
	cfg := &config.Config{}
	for _, s := range sources {
		cfg.Sources = append(cfg.Sources, config.Source{Name: s, Path: "/tmp/" + s})
	}
	for _, s := range sinks {
		cfg.Sinks = append(cfg.Sinks, config.Sink{Name: s, Type: "stdout"})
	}
	return cfg
}

func TestBuild_ValidConfig(t *testing.T) {
	cfg := makeConfig([]string{"app", "sys"}, []string{"console"})
	g, err := Build(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(g.Nodes) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(g.Nodes))
	}
}

func TestBuild_NilConfig_ReturnsError(t *testing.T) {
	_, err := Build(nil)
	if err == nil {
		t.Fatal("expected error for nil config")
	}
}

func TestBuild_DuplicateSourceName_ReturnsError(t *testing.T) {
	cfg := makeConfig([]string{"app", "app"}, []string{"out"})
	_, err := Build(cfg)
	if err == nil {
		t.Fatal("expected error for duplicate node id")
	}
}

func TestBuild_DuplicateSinkName_ReturnsError(t *testing.T) {
	cfg := makeConfig([]string{"src"}, []string{"sink", "sink"})
	_, err := Build(cfg)
	if err == nil {
		t.Fatal("expected error for duplicate sink id")
	}
}

func TestBuild_EmptySourcesAndSinks(t *testing.T) {
	cfg := makeConfig(nil, nil)
	g, err := Build(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(g.Nodes) != 0 {
		t.Fatalf("expected 0 nodes, got %d", len(g.Nodes))
	}
}

func TestBuild_NodeKindsAssigned(t *testing.T) {
	cfg := makeConfig([]string{"src"}, []string{"dst"})
	g, err := Build(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	kinds := map[string]string{}
	for _, n := range g.Nodes {
		kinds[n.ID] = n.Kind
	}
	if kinds["src"] != "source" {
		t.Errorf("expected kind=source for src, got %q", kinds["src"])
	}
	if kinds["dst"] != "sink" {
		t.Errorf("expected kind=sink for dst, got %q", kinds["dst"])
	}
}
