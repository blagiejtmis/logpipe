// Package sink provides log sink implementations for logpipe.
//
// A Sink accepts a single log line (plain string) and writes it to an
// underlying destination — currently stdout or a local file.  Each sink
// supports two output formats:
//
//   - "text" — the raw log line followed by a newline.
//   - "json" — the log line wrapped in a JSON object: {"log": "<line>"}.
//
// # Creating sinks
//
// Individual sinks are constructed via their typed constructors:
//
//	stdout, err := sink.NewStdoutSink(cfg)
//	file,   err := sink.NewFileSink(cfg)
//
// # Managing multiple sinks
//
// Use Manager to fan a single log line out to several sinks at once:
//
//	m, err := sink.NewManager(cfgs)
//	if err != nil { ... }
//	defer m.Close()
//	m.Write(line)
//
// # Error handling
//
// Write errors from individual sinks are collected and returned as a combined
// error so that a failure in one sink does not silently suppress writes to the
// remaining sinks.  Callers should check the returned error and decide whether
// to propagate or log it.
package sink
