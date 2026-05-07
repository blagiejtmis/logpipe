// Package checkpoint provides durable offset tracking for logpipe tail sources.
//
// When logpipe restarts it needs to resume reading each source file from the
// position it last processed, rather than re-emitting old log lines or missing
// lines that arrived during downtime.
//
// # Usage
//
//	store, err := checkpoint.New("/var/lib/logpipe/checkpoint.json")
//	if err != nil { ... }
//
//	// Read the last saved offset before opening a tailer.
//	offset := store.Get("/var/log/app.log")
//
//	// After each batch of lines is processed, persist the new offset.
//	if err := store.Set("/var/log/app.log", newOffset); err != nil { ... }
//
// The store is safe for concurrent use and flushes atomically via a
// write-then-rename strategy to avoid partial writes.
package checkpoint
