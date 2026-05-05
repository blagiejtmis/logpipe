// Package filter provides configurable log line filtering for logpipe.
//
// A Filter is constructed from a slice of Rule values, each specifying:
//   - Field:   the log field name to inspect (use "_raw" for the full line).
//   - Pattern: a regular expression the field value must satisfy.
//   - Invert:  when true, the rule passes only when the pattern does NOT match.
//
// Rules are evaluated with AND semantics: every rule must pass for a log line
// to be forwarded downstream. A Filter with no rules passes all lines.
//
// Typical usage inside a pipeline stage:
//
//	fields := map[string]string{
//		"level":   entry.Level,
//		"message": entry.Message,
//		"_raw":    entry.Raw,
//	}
//	if f.Match(fields) {
//		// forward entry
//	}
package filter
