// Package router provides source-aware and field-aware routing of log records
// to named sink destinations.
//
// A Router is constructed from an ordered list of Rules. Each Rule carries:
//   - SourcePattern: a regular expression matched against the origin source name.
//   - FieldKey / FieldPattern: an optional field equality check on the parsed
//     log record.
//   - Sinks: the list of sink names to deliver matching records to.
//
// Rules are evaluated in declaration order; the first match wins. When no rule
// matches, the Router falls back to its configured default sinks.
//
// Example configuration (YAML):
//
//	routing:
//	  default_sinks: [stdout]
//	  rules:
//	    - source: "^app\\.log$"
//	      field:  level
//	      match:  "^error$"
//	      sinks:  [errors-file, stdout]
package router
