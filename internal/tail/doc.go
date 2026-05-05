// Package tail provides file-tailing primitives for logpipe.
//
// A Tailer watches a single file for new lines and emits them on a channel.
// A Manager supervises multiple Tailers, fanning all output into a single
// read-only channel that downstream sinks can consume.
//
// Basic usage:
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	mgr := tail.NewManager(100)
//	mgr.Add(ctx, "/var/log/app.log")
//	mgr.Add(ctx, "/var/log/access.log")
//
//	for line := range mgr.Output() {
//		fmt.Printf("[%s] %s", line.Source, line.Text)
//	}
package tail
