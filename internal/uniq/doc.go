// Package uniq provides per-field uniqueness enforcement for log records.
//
// Unlike the dedupe package which hashes entire records, uniq tracks a single
// named field and suppresses records whose field value has already been seen
// within a configurable sliding time window.
//
// # Basic Usage
//
//	 u, err := uniq.New("request_id", 5*time.Minute)
//	 if err != nil { ... }
//
//	 if u.Allow(record) {
//	     // record has a unique request_id — forward it
//	 }
//
// Records where the key field is absent or not a string are always forwarded.
// Expired entries are pruned lazily on each Allow call.
package uniq
