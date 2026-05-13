// Package topology provides pipeline graph construction and validation
// for logpipe.
//
// # Overview
//
// Before any data flows, logpipe builds a directed acyclic graph (DAG)
// whose nodes represent sources, transforms, filters, and sinks.  The
// topology package is responsible for:
//
//   - Parsing the config into a set of named Node values.
//   - Verifying that every dependency reference resolves to a known node.
//   - Detecting cycles that would cause the pipeline to deadlock.
//
// # Usage
//
//	g, err := topology.Build(cfg)
//	if err != nil {
//	    log.Fatalf("invalid topology: %v", err)
//	}
//
// The resulting Graph is consumed by the orchestrator to determine
// start-up and shutdown ordering.
package topology
