// Package alert provides threshold-based alerting for log records.
//
// An Alerter evaluates each log record against a set of rules. Each rule
// specifies a field name, a regex pattern to match the field value, a hit
// threshold, and a callback that fires once the threshold is reached.
//
// # Basic usage
//
//	a, err := alert.New(
//		alert.WithRule("level", "error", 5, func(field, value string, count int) {
//			log.Printf("alert: %s=%s fired %d times", field, value, count)
//		}),
//	)
//	if err != nil { ... }
//	a.Evaluate(record)
//
// # Manager
//
// NewManager wires global and per-source Alerters from a ManagerConfig,
// optionally attaching a metrics.Registry so that each alert trigger is
// reflected as an incrementing counter.
//
// AlerterFor(source) returns the most specific Alerter for a given source,
// falling back to the global one when no source-specific rule exists.
package alert
