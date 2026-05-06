// Package healthcheck provides a lightweight health-check subsystem for
// logpipe.
//
// # Overview
//
// A [Checker] maintains a map of named component health states. Any part of
// the application (tailers, sinks, pipeline) can call Set to mark itself as
// [StatusOK] or [StatusDegraded]. The overall status exposed via [Report] is
// degraded if at least one component is degraded.
//
// # HTTP endpoint
//
// [Checker.Handler] returns an http.Handler that serialises the current
// [Report] as JSON. The HTTP status code is 200 when all components are
// healthy and 503 when any component is degraded.
//
// [Server] wraps the handler in a standalone HTTP server with graceful
// shutdown driven by a context, making it easy to integrate with the
// orchestrator.
//
// # Usage
//
//	checker := healthcheck.New()
//	checker.Set("tailer", healthcheck.ComponentHealth{Status: healthcheck.StatusOK})
//	http.Handle("/healthz", checker.Handler())
package healthcheck
