// Package join provides a record joiner that correlates log records
// from multiple sources on a shared key field within a time window.
package join

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Record is a map of string keys to arbitrary values.
type Record map[string]interface{}

// Rule describes how to join records from two sources.
type Rule struct {
	// LeftSource is the name of the primary source.
	LeftSource string
	// RightSource is the name of the secondary source.
	RightSource string
	// KeyField is the field name used to correlate records.
	KeyField string
	// Window is how long to wait for the matching record.
	Window time.Duration
	// OutputField is the field under which merged right-side fields are placed.
	// If empty the right-side fields are merged directly into the output record.
	OutputField string
}

type pending struct {
	rec    Record
	expiry time.Time
}

// Joiner correlates records from two sources on a shared key.
type Joiner struct {
	rule Rule
	mu   sync.Mutex
	// left holds unmatched left records keyed by the join key value.
	left map[string]pending
	// right holds unmatched right records keyed by the join key value.
	right map[string]pending
}

// New creates a Joiner from the given Rule.
func New(r Rule) (*Joiner, error) {
	if r.LeftSource == "" {
		return nil, errors.New("join: LeftSource must not be empty")
	}
	if r.RightSource == "" {
		return nil, errors.New("join: RightSource must not be empty")
	}
	if r.KeyField == "" {
		return nil, errors.New("join: KeyField must not be empty")
	}
	if r.Window <= 0 {
		return nil, errors.New("join: Window must be positive")
	}
	return &Joiner{
		rule:  r,
		left:  make(map[string]pending),
		right: make(map[string]pending),
	}, nil
}

// Add accepts a record from the named source. It returns the merged record and
// true when a join is completed, or nil and false when the record is buffered.
func (j *Joiner) Add(source string, rec Record) (Record, bool) {
	key, ok := rec[j.rule.KeyField]
	if !ok {
		return nil, false
	}
	keyStr := fmt.Sprintf("%v", key)
	now := time.Now()

	j.mu.Lock()
	defer j.mu.Unlock()

	j.evict(now)

	switch source {
	case j.rule.LeftSource:
		if rp, found := j.right[keyStr]; found {
			delete(j.right, keyStr)
			return j.merge(rec, rp.rec), true
		}
		j.left[keyStr] = pending{rec: rec, expiry: now.Add(j.rule.Window)}
	case j.rule.RightSource:
		if lp, found := j.left[keyStr]; found {
			delete(j.left, keyStr)
			return j.merge(lp.rec, rec), true
		}
		j.right[keyStr] = pending{rec: rec, expiry: now.Add(j.rule.Window)}
	}
	return nil, false
}

func (j *Joiner) merge(left, right Record) Record {
	out := make(Record, len(left))
	for k, v := range left {
		out[k] = v
	}
	if j.rule.OutputField != "" {
		out[j.rule.OutputField] = right
	} else {
		for k, v := range right {
			out[k] = v
		}
	}
	return out
}

func (j *Joiner) evict(now time.Time) {
	for k, p := range j.left {
		if now.After(p.expiry) {
			delete(j.left, k)
		}
	}
	for k, p := range j.right {
		if now.After(p.expiry) {
			delete(j.right, k)
		}
	}
}
