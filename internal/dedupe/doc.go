// Package dedupe provides log record deduplication for logpipe.
//
// It suppresses repeated log records that share the same value for a
// configured field within a sliding time window.
//
// # Deduplicator
//
// A Deduplicator tracks seen values for a single field over a fixed window
// duration. Calling Allow returns false when the same value has already been
// seen within the current window.
//
// # Manager
//
// NewManager constructs a Manager from configuration. The manager holds one
// Deduplicator per source (or a global fallback) and exposes a single Allow
// method that selects the appropriate deduplicator by source name.
//
// Example configuration:
//
//	dedup:
//	  global:
//	    field: message
//	    window: 10s
//	  sources:
//	    app.log:
//	      field: request_id
//	      window: 30s
package dedupe
