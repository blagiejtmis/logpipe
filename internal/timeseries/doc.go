// Package timeseries implements a sliding-window frequency counter for
// log-record field values.
//
// # Overview
//
// A [Series] divides a configurable time window into N equal-sized buckets.
// Each call to [Series.Record] increments the count for the value of the
// tracked field in the current bucket. Expired buckets are lazily cleared
// on the next Record or Counts call, keeping memory usage bounded.
//
// # Usage
//
//	s, err := timeseries.New("level", time.Minute, 6)
//	if err != nil { ... }
//	s.Record(rec)          // called per log record
//	counts := s.Counts()   // map[string]int64 snapshot
//
// # Manager
//
// [NewManager] builds a source-aware Manager from a [Config], applying
// per-source field overrides when present and falling back to the global
// field for all other sources.
package timeseries
