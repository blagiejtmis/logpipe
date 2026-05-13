package topology

import (
	"fmt"

	"github.com/logpipe/logpipe/internal/config"
)

// Manager holds the validated pipeline graph and provides
// convenience accessors used by the orchestrator.
type Manager struct {
	graph *Graph
}

// NewManager builds the topology graph from cfg and returns a Manager.
// An error is returned if the graph is invalid.
func NewManager(cfg *config.Config) (*Manager, error) {
	g, err := Build(cfg)
	if err != nil {
		return nil, fmt.Errorf("topology manager: %w", err)
	}
	return &Manager{graph: g}, nil
}

// SourceIDs returns the IDs of all source nodes.
func (m *Manager) SourceIDs() []string {
	return m.nodeIDsByKind("source")
}

// SinkIDs returns the IDs of all sink nodes.
func (m *Manager) SinkIDs() []string {
	return m.nodeIDsByKind("sink")
}

// NodeCount returns the total number of nodes in the graph.
func (m *Manager) NodeCount() int {
	return len(m.graph.Nodes)
}

func (m *Manager) nodeIDsByKind(kind string) []string {
	var ids []string
	for _, n := range m.graph.Nodes {
		if n.Kind == kind {
			ids = append(ids, n.ID)
		}
	}
	return ids
}
