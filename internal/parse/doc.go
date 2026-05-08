// Package parse implements log line parsers for logpipe.
//
// A Parser converts a raw string (one line from a tailed source) into a
// structured Record (map[string]string) that the rest of the pipeline can
// filter, transform, redact, and route.
//
// # Supported Formats
//
//   - "json"   – expects a JSON object; numeric/boolean values are coerced to
//     their string representations.
//   - "logfmt" – expects space-separated key=value pairs, optionally
//     double-quoted values.  A bare token (no '=') is stored under the
//     key "msg" if no explicit msg key has been seen yet.
//
// # Time Injection
//
// Both parsers inject a "time" field (RFC-3339, UTC) when the incoming line
// does not already carry one, ensuring downstream components always have a
// timestamp to work with.
//
// # Usage
//
//	p, err := parse.New("json")
//	if err != nil { /* unsupported format */ }
//	rec, err := p.Parse(line)
package parse
