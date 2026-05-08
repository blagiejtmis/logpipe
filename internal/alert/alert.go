// Package alert provides threshold-based alerting over metric counters.
// When a counter exceeds a configured threshold within an evaluation window,
// an alert is fired via a user-supplied callback.
package alert

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"logpipe/internal/metrics"
)

// AlertFunc is called when a threshold is breached.
// name is the counter name, value is the current snapshot value.
type AlertFunc func(name string, value int64)

// Rule defines a single alerting rule.
type Rule struct {
	CounterName string
	Threshold   int64
	Cooldown    time.Duration // minimum time between repeated alerts
}

// Alerter watches a metrics registry and fires alerts when rules are breached.
type Alerter struct {
	reg      *metrics.Registry
	rules    []Rule
	onAlert  AlertFunc
	lastFire map[string]time.Time
	mu       sync.Mutex
}

// New creates an Alerter. reg and onAlert must not be nil; rules must be non-empty.
func New(reg *metrics.Registry, rules []Rule, onAlert AlertFunc) (*Alerter, error) {
	if reg == nil {
		return nil, errors.New("alert: registry must not be nil")
	}
	if onAlert == nil {
		return nil, errors.New("alert: onAlert callback must not be nil")
	}
	if len(rules) == 0 {
		return nil, errors.New("alert: at least one rule is required")
	}
	for i, r := range rules {
		if r.CounterName == "" {
			return nil, fmt.Errorf("alert: rule[%d] has empty CounterName", i)
		}
		if r.Threshold <= 0 {
			return nil, fmt.Errorf("alert: rule[%d] threshold must be > 0", i)
		}
	}
	return &Alerter{
		reg:      reg,
		rules:    rules,
		onAlert:  onAlert,
		lastFire: make(map[string]time.Time),
	}, nil
}

// Evaluate checks all rules against the current registry snapshot.
// It should be called periodically (e.g. from a ticker).
func (a *Alerter) Evaluate() {
	snap := a.reg.Snapshot()
	now := time.Now()

	a.mu.Lock()
	defer a.mu.Unlock()

	for _, rule := range a.rules {
		val, ok := snap[rule.CounterName]
		if !ok {
			continue
		}
		if val < rule.Threshold {
			continue
		}
		last, fired := a.lastFire[rule.CounterName]
		if fired && now.Sub(last) < rule.Cooldown {
			continue
		}
		a.lastFire[rule.CounterName] = now
		a.onAlert(rule.CounterName, val)
	}
}
