package alert

import (
	"sync/atomic"
	"testing"
	"time"

	"logpipe/internal/metrics"
)

func newReg(t *testing.T) *metrics.Registry {
	t.Helper()
	return metrics.NewRegistry()
}

func TestNew_NilRegistry_ReturnsError(t *testing.T) {
	_, err := New(nil, []Rule{{CounterName: "x", Threshold: 1}}, func(string, int64) {})
	if err == nil {
		t.Fatal("expected error for nil registry")
	}
}

func TestNew_NilCallback_ReturnsError(t *testing.T) {
	_, err := New(newReg(t), []Rule{{CounterName: "x", Threshold: 1}}, nil)
	if err == nil {
		t.Fatal("expected error for nil callback")
	}
}

func TestNew_EmptyRules_ReturnsError(t *testing.T) {
	_, err := New(newReg(t), nil, func(string, int64) {})
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestNew_ZeroThreshold_ReturnsError(t *testing.T) {
	_, err := New(newReg(t), []Rule{{CounterName: "x", Threshold: 0}}, func(string, int64) {})
	if err == nil {
		t.Fatal("expected error for zero threshold")
	}
}

func TestEvaluate_BelowThreshold_NoAlert(t *testing.T) {
	reg := newReg(t)
	c := reg.Counter("errors")
	c.Inc()

	fired := false
	a, _ := New(reg, []Rule{{CounterName: "errors", Threshold: 5}}, func(string, int64) { fired = true })
	a.Evaluate()

	if fired {
		t.Fatal("alert should not fire below threshold")
	}
}

func TestEvaluate_AtThreshold_FiresAlert(t *testing.T) {
	reg := newReg(t)
	c := reg.Counter("errors")
	c.Add(5)

	var fired int32
	a, _ := New(reg, []Rule{{CounterName: "errors", Threshold: 5, Cooldown: 0}}, func(string, int64) {
		atomic.AddInt32(&fired, 1)
	})
	a.Evaluate()

	if atomic.LoadInt32(&fired) != 1 {
		t.Fatalf("expected 1 alert, got %d", fired)
	}
}

func TestEvaluate_CooldownSuppressesRepeat(t *testing.T) {
	reg := newReg(t)
	c := reg.Counter("drops")
	c.Add(10)

	var fired int32
	a, _ := New(reg, []Rule{{CounterName: "drops", Threshold: 1, Cooldown: 10 * time.Second}},
		func(string, int64) { atomic.AddInt32(&fired, 1) })

	a.Evaluate()
	a.Evaluate()

	if atomic.LoadInt32(&fired) != 1 {
		t.Fatalf("cooldown should suppress second alert; got %d fires", fired)
	}
}

func TestEvaluate_UnknownCounter_NoAlert(t *testing.T) {
	reg := newReg(t)

	fired := false
	a, _ := New(reg, []Rule{{CounterName: "missing", Threshold: 1}}, func(string, int64) { fired = true })
	a.Evaluate()

	if fired {
		t.Fatal("should not alert on unknown counter")
	}
}
