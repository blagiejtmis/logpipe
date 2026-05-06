// Package metrics provides a lightweight in-process metrics registry
// for tracking counters across logpipe components.
//
// # Registry
//
// NewRegistry returns a thread-safe registry. Counters are created on
// first access via Registry.Counter and can be snapshotted at any time.
//
// # Reporter
//
// NewReporter periodically logs a snapshot of all counters at a
// configurable interval. It respects context cancellation and emits a
// final snapshot on shutdown.
//
// # HTTP
//
// Handler exposes the current snapshot as a JSON endpoint suitable for
// scraping by monitoring systems.
package metrics
