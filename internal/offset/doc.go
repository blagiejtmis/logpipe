// Package offset provides a lightweight persistent store for per-source byte
// offsets used by logpipe's tail subsystem.
//
// When logpipe restarts it needs to know where it left off in each source file
// so that it does not re-emit lines that were already processed.  The Store
// type serialises offsets as a JSON file on disk and exposes simple Get/Set/
// Delete operations.  All methods are safe for concurrent use.
//
// Typical usage:
//
//	store, err := offset.New("/var/lib/logpipe/offsets.json")
//	if err != nil { ... }
//
//	// Restore position before tailing.
//	pos := store.Get(sourceName)
//
//	// Persist position after each successful flush.
//	_ = store.Set(sourceName, newPos)
package offset
