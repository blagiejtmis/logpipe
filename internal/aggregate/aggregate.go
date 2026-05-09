// Package aggregate provides time-window based field aggregation for log records.
package aggregate

import (
	"errors"
	"sync"
	"time"
)

// Op is an aggregation operation.
type Op string

const (
	OpCount Op = "count"
	OpSum   Op = "sum"
	OpMin   Op = "min"
	OpMax   Op = "max"
)

// Rule defines how a field should be aggregated.
type Rule struct {
	Field  string
	Op     Op
	Window time.Duration
}

// Aggregator accumulates values per source over a rolling window.
type Aggregator struct {
	rules  []Rule
	mu     sync.Mutex
	bucket map[string]map[string]*cell // source -> field -> cell
}

type cell struct {
	op      Op
	value   float64
	count   int64
	expires time.Time
}

// New creates an Aggregator from the given rules.
func New(rules []Rule) (*Aggregator, error) {
	for _, r := range rules {
		if r.Field == "" {
			return nil, errors.New("aggregate: field name must not be empty")
		}
		if r.Window <= 0 {
			return nil, errors.New("aggregate: window must be positive")
		}
		switch r.Op {
		case OpCount, OpSum, OpMin, OpMax:
		default:
			return nil, errors.New("aggregate: unknown op " + string(r.Op))
		}
	}
	return &Aggregator{
		rules:  rules,
		bucket: make(map[string]map[string]*cell),
	}, nil
}

// Add incorporates a record value into the aggregation state.
// record is a map of field name -> value. source identifies the log source.
func (a *Aggregator) Add(source string, record map[string]any) {
	a.mu.Lock()
	defer a.mu.Unlock()
	now := time.Now()
	if _, ok := a.bucket[source]; !ok {
		a.bucket[source] = make(map[string]*cell)
	}
	for _, r := range a.rules {
		v, _ := toFloat(record[r.Field])
		c, ok := a.bucket[source][r.Field]
		if !ok || now.After(c.expires) {
			c = &cell{op: r.Op, expires: now.Add(r.Window)}
			a.bucket[source][r.Field] = c
		}
		applyOp(c, v)
	}
}

// Snapshot returns the current aggregated values per source and field.
func (a *Aggregator) Snapshot() map[string]map[string]float64 {
	a.mu.Lock()
	defer a.mu.Unlock()
	out := make(map[string]map[string]float64, len(a.bucket))
	now := time.Now()
	for src, fields := range a.bucket {
		out[src] = make(map[string]float64, len(fields))
		for field, c := range fields {
			if now.After(c.expires) {
				out[src][field] = 0
				continue
			}
			out[src][field] = result(c)
		}
	}
	return out
}

func applyOp(c *cell, v float64) {
	switch c.op {
	case OpCount:
		c.count++
	case OpSum:
		c.value += v
	case OpMin:
		if c.count == 0 || v < c.value {
			c.value = v
		}
		c.count++
	case OpMax:
		if c.count == 0 || v > c.value {
			c.value = v
		}
		c.count++
	}
}

func result(c *cell) float64 {
	if c.op == OpCount {
		return float64(c.count)
	}
	return c.value
}

func toFloat(v any) (float64, bool) {
	switch t := v.(type) {
	case float64:
		return t, true
	case int:
		return float64(t), true
	case int64:
		return float64(t), true
	}
	return 0, false
}
