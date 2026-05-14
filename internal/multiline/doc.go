// Package multiline implements multi-line log record assembly for logpipe.
//
// Many log sources emit records that span several raw lines — Java stack traces,
// Python tracebacks, or any structured block with a header and indented
// continuation lines.  This package collects those raw lines and joins them
// into a single logical record before the record enters the rest of the
// processing pipeline.
//
// # Usage
//
// Create an [Assembler] with a [Rule] that describes either:
//   - StartPattern — a regular expression that matches the *first* line of a
//     new record.  Every line that matches starts a fresh record; preceding
//     buffered lines are flushed.
//   - ContinuePattern — a regular expression that matches *continuation* lines.
//     A line that does not match triggers a flush of the current buffer.
//
// Only one of the two patterns may be set at a time.
//
// Optional limits:
//   - MaxLines — flush when the buffer reaches this many lines.
//   - Timeout  — flush when no new line has arrived within this duration.
//   - Field    — the record field that receives the joined text (default: "message").
//
// # Manager
//
// [NewManager] reads a [Config] block and provides per-source [Assembler]
// instances via [Manager.Assembler].  Sources without an explicit rule inherit
// the default rule; if no default is configured the method returns nil
// (pass-through, no assembly).
package multiline
