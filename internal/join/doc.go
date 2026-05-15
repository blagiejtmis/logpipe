// Package join provides a streaming record joiner for logpipe.
//
// A Joiner correlates log records arriving from two distinct sources
// ("left" and "right") by matching on a shared key field. When both
// sides arrive within the configured time window the records are merged
// and emitted as a single output record.
//
// # Usage
//
//	rule := join.Rule{
//	    LeftSource:  "app",
//	    RightSource: "db",
//	    KeyField:    "request_id",
//	    Window:      5 * time.Second,
//	    OutputField: "db_span", // optional; omit to merge fields inline
//	}
//	
//	j, err := join.New(rule)
//	if err != nil { ... }
//	
//	// Feed records from each source:
//	if merged, ok := j.Add(sourceName, rec); ok {
//	    // merged record is ready for downstream processing
//	}
//
// Records that do not find a counterpart within the window are silently
// discarded. Expired entries are evicted lazily on each call to Add.
package join
