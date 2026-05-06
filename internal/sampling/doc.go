// Package sampling implements probabilistic sampling for log records.
//
// A Sampler accepts a rate in the range (0.0, 1.0] and uses a
// cryptographically-seeded pseudo-random source to decide whether each
// record should be forwarded to downstream sinks.
//
// # Sampler
//
// New creates a single-source sampler. Call Allow with a source key;
// the key is used only for logical grouping — each Sampler maintains its
// own independent state.
//
// # Manager
//
// NewManager reads a SamplingConfig and builds a set of samplers:
//
//   - A global (default) sampler applied to all sources not explicitly listed.
//   - Per-source samplers that override the global rate for specific sources.
//
// If no configuration is provided, Manager.Allow always returns true so that
// the pipeline operates without any sampling overhead.
package sampling
