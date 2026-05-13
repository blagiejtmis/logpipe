// Package lookup provides table-driven field enrichment for log records.
//
// A Lookup applier holds one or more LookupRules. Each rule specifies:
//
//   - KeyField  – the record field whose value is used as the lookup key.
//   - Table     – a map from key strings to a map of fields to merge into the
//     record when the key matches.
//   - DestField – used when the table row contains a single value; otherwise
//     all fields in the matched row are merged.
//   - OnMiss    – "keep" (default) leaves the record unchanged when the key is
//     absent from the table; "drop" signals that the record should be
//     discarded by returning ErrDrop.
//
// The Manager selects the most specific set of rules for each source:
// source-specific rules take precedence over global rules. When neither
// applies the record is passed through unmodified.
//
// Typical usage:
//
//	m, err := lookup.NewManager(cfg.Lookup)
//	if err != nil { … }
//	enriched, err := m.Apply(record["source"].(string), record)
package lookup
