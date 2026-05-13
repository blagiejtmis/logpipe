// Package sequence provides monotonically increasing sequence number stamping
// for log records flowing through the logpipe pipeline.
//
// # Overview
//
// A [Sequencer] atomically increments a counter for every record it processes
// and writes the resulting value into a configurable field (default: "seq").
// This makes it straightforward to detect dropped records or restore ordering
// after fan-out.
//
// # Usage
//
//	s, err := sequence.New("seq", 1)
//	if err != nil { ... }
//	stamped := s.Apply(record)
//
// # Manager
//
// [NewManager] wires sequencers from a [config.SequenceConfig].  A global
// Default rule covers all sources; entries in Sources override the default
// for a named source.  When no configuration is provided, records pass
// through unchanged.
package sequence
