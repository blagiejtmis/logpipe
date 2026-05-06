// Package healthcheck provides an HTTP health endpoint for logpipe.
package healthcheck

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// Status represents the health status of a component.
type Status string

const (
	StatusOK      Status = "ok"
	StatusDegraded Status = "degraded"
)

// ComponentHealth holds the health state of a named component.
type ComponentHealth struct {
	Status  Status `json:"status"`
	Message string `json:"message,omitempty"`
}

// Report is the full health report returned by the handler.
type Report struct {
	Status     Status                     `json:"status"`
	Timestamp  time.Time                  `json:"timestamp"`
	Components map[string]ComponentHealth `json:"components"`
}

// Checker aggregates component health checks.
type Checker struct {
	mu         sync.RWMutex
	components map[string]ComponentHealth
}

// New creates a new Checker with no registered components.
func New() *Checker {
	return &Checker{
		components: make(map[string]ComponentHealth),
	}
}

// Set registers or updates the health of a named component.
func (c *Checker) Set(name string, h ComponentHealth) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.components[name] = h
}

// Report builds the current health report. Overall status is degraded
// if any component is degraded.
func (c *Checker) Report() Report {
	c.mu.RLock()
	defer c.mu.RUnlock()

	overall := StatusOK
	comps := make(map[string]ComponentHealth, len(c.components))
	for k, v := range c.components {
		comps[k] = v
		if v.Status == StatusDegraded {
			overall = StatusDegraded
		}
	}
	return Report{
		Status:     overall,
		Timestamp:  time.Now().UTC(),
		Components: comps,
	}
}

// Handler returns an http.Handler that serves the health report as JSON.
func (c *Checker) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		report := c.Report()
		code := http.StatusOK
		if report.Status == StatusDegraded {
			code = http.StatusServiceUnavailable
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		_ = json.NewEncoder(w).Encode(report)
	})
}
