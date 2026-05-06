package healthcheck_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/logpipe/internal/healthcheck"
)

func TestChecker_DefaultReport_IsOK(t *testing.T) {
	c := healthcheck.New()
	r := c.Report()
	if r.Status != healthcheck.StatusOK {
		t.Fatalf("expected ok, got %s", r.Status)
	}
	if len(r.Components) != 0 {
		t.Fatalf("expected no components, got %d", len(r.Components))
	}
}

func TestChecker_DegradedComponent_DegradedOverall(t *testing.T) {
	c := healthcheck.New()
	c.Set("tailer", healthcheck.ComponentHealth{Status: healthcheck.StatusOK})
	c.Set("sink", healthcheck.ComponentHealth{Status: healthcheck.StatusDegraded, Message: "write error"})

	r := c.Report()
	if r.Status != healthcheck.StatusDegraded {
		t.Fatalf("expected degraded, got %s", r.Status)
	}
}

func TestChecker_AllOK_ReportsOK(t *testing.T) {
	c := healthcheck.New()
	c.Set("tailer", healthcheck.ComponentHealth{Status: healthcheck.StatusOK})
	c.Set("sink", healthcheck.ComponentHealth{Status: healthcheck.StatusOK})

	if got := c.Report().Status; got != healthcheck.StatusOK {
		t.Fatalf("expected ok, got %s", got)
	}
}

func TestHandler_OKResponse(t *testing.T) {
	c := healthcheck.New()
	c.Set("pipeline", healthcheck.ComponentHealth{Status: healthcheck.StatusOK})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	c.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var report healthcheck.Report
	if err := json.NewDecoder(rec.Body).Decode(&report); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if report.Status != healthcheck.StatusOK {
		t.Fatalf("expected ok in body")
	}
}

func TestHandler_DegradedReturns503(t *testing.T) {
	c := healthcheck.New()
	c.Set("sink", healthcheck.ComponentHealth{Status: healthcheck.StatusDegraded, Message: "disk full"})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	c.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", rec.Code)
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	c := healthcheck.New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/healthz", nil)
	c.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestHandler_ContentType(t *testing.T) {
	c := healthcheck.New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	c.Handler().ServeHTTP(rec, req)

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Fatalf("expected application/json, got %s", ct)
	}
}
