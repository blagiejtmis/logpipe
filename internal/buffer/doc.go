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
package buffer
