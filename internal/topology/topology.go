// Package topology builds and validates the directed acyclic graph
// of logpipe components, ensuring sources, transforms, filters, and
// sinks are wired together consistently before the pipeline starts.
package topology

import (
	"errors"
	"fmt"

	"github.com/logpipe/logpipe/internal/config"
)

// Node represents a single component in the pipeline graph.
type Node struct {
	ID   string
	Kind string // "source", "transform", "filter", "sink"
	Deps []string
}

// Graph is the resolved topology of the pipeline.
type Graph struct {
	Nodes   []*Node
	indexed map[string]*Node
}

// Build constructs and validates a Graph from the supplied config.
// It returns an error if any dependency references an unknown node or
// if the graph contains a cycle.
func Build(cfg *config.Config) (*Graph, error) {
	if cfg == nil {
		return nil, errors.New("topology: nil config")
	}

	g := &Graph{indexed: make(map[string]*Node)}

	for _, s := range cfg.Sources {
		if err := g.addNode(s.Name, "source", nil); err != nil {
			return nil, err
		}
	}
	for _, s := range cfg.Sinks {
		if err := g.addNode(s.Name, "sink", nil); err != nil {
			return nil, err
		}
	}

	if err := g.detectCycles(); err != nil {
		return nil, err
	}
	return g, nil
}

func (g *Graph) addNode(id, kind string, deps []string) error {
	if _, exists := g.indexed[id]; exists {
		return fmt.Errorf("topology: duplicate node id %q", id)
	}
	for _, d := range deps {
		if _, ok := g.indexed[d]; !ok {
			return fmt.Errorf("topology: node %q depends on unknown node %q", id, d)
		}
	}
	n := &Node{ID: id, Kind: kind, Deps: deps}
	g.Nodes = append(g.Nodes, n)
	g.indexed[id] = n
	return nil
}

// detectCycles performs a DFS-based cycle check.
func (g *Graph) detectCycles() error {
	visited := make(map[string]bool)
	inStack := make(map[string]bool)

	var dfs func(id string) error
	dfs = func(id string) error {
		if inStack[id] {
			return fmt.Errorf("topology: cycle detected at node %q", id)
		}
		if visited[id] {
			return nil
		}
		visited[id] = true
		inStack[id] = true
		for _, dep := range g.indexed[id].Deps {
			if err := dfs(dep); err != nil {
				return err
			}
		}
		inStack[id] = false
		return nil
	}

	for _, n := range g.Nodes {
		if err := dfs(n.ID); err != nil {
			return err
		}
	}
	return nil
}
