// Package jitter provides randomised delay injection for log records.
//
// # Overview
//
// When log sources emit records in large bursts, downstream sinks can become
// overwhelmed. The jitter package lets operators configure a random sleep
// window so that records are spread across time before being forwarded.
//
// # Usage
//
//	// Create a jitter that waits between 0 and 20 ms.
//	j, err := jitter.New(0, 20*time.Millisecond)
//	if err != nil {
//		log.Fatal(err)
//	}
//	j.Wait() // blocks for a random duration in [0, 20ms]
//
// # Manager
//
// NewManager builds a Manager from a Config, which can specify a default
// delay range and per-source overrides. Sources with no explicit rule use
// the default; a nil Config produces a zero-delay (no-op) manager.
package jitter
