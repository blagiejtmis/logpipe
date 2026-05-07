// Package buffer implements a bounded, thread-safe ring buffer for
// log records produced by tail sources.
//
// # Overview
//
// When a downstream sink is slow or temporarily unavailable, records
// accumulate in the buffer rather than being dropped immediately. Two
// overflow strategies are supported:
//
//   - Evict-oldest (dropNew=false): the oldest unprocessed record is
//     silently discarded to make room for the new one. This preserves
//     the most recent data.
//
//   - Drop-new (dropNew=true): the incoming record is rejected and
//     ErrFull is returned to the caller. The existing queue is
//     preserved intact.
//
// # Errors
//
// ErrFull is returned by Push when the buffer is at capacity and the
// drop-new strategy is active. ErrInvalidCapacity is returned by New
// when the requested capacity is less than 1.
//
// # Usage
//
//	buf, err := buffer.New(1024, false)
//	if err != nil { ... }
//
//	buf.Push(buffer.Record{Source: "app.log", Fields: fields})
//
//	if r, ok := buf.Pop(); ok {
//	    // forward r to sink
//	}
//
// Dropped returns the cumulative number of records lost to overflow,
// which can be exposed via the metrics registry.
//
// # Concurrency
//
// All methods on Buffer are safe for concurrent use by multiple
// goroutines without additional synchronisation.
package buffer
