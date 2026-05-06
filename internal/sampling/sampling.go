// Package sampling provides probabilistic and rate-based log sampling.
// It allows reducing log volume by forwarding only a fraction of records
// matching configurable criteria.
package sampling

import (
	"fmt"
	"math/rand"
	"sync"
)

// Config holds sampling configuration for a single rule.
type Config struct {
	// Rate is the fraction of records to keep, in the range (0.0, 1.0].
	// A rate of 1.0 keeps all records; 0.5 keeps half.
	Rate float64 `yaml:"rate"`
	// Source restricts the rule to a specific source name.
	// An empty string applies the rule to all sources.
	Source string `yaml:"source,omitempty"`
}

// Sampler decides whether an individual log record should be forwarded.
type Sampler struct {
	mu      sync.Mutex
	rng     *rand.Rand
	default_ float64
	sources  map[string]float64
}

// New creates a Sampler from a slice of Config rules.
// Rules with a non-empty Source override the default rate for that source.
// Returns an error if any rate is outside (0.0, 1.0].
func New(rules []Config) (*Sampler, error) {
	s := &Sampler{
		rng:      rand.New(rand.NewSource(rand.Int63())),
		default_: 1.0,
		sources:  make(map[string]float64),
	}
	for _, r := range rules {
		if r.Rate <= 0 || r.Rate > 1.0 {
			return nil, fmt.Errorf("sampling: rate %.4f out of range (0, 1]", r.Rate)
		}
		if r.Source == "" {
			s.default_ = r.Rate
		} else {
			s.sources[r.Source] = r.Rate
		}
	}
	return s, nil
}

// Allow returns true if the record from the given source should be forwarded.
func (s *Sampler) Allow(source string) bool {
	rate := s.default_
	if r, ok := s.sources[source]; ok {
		rate = r
	}
	if rate >= 1.0 {
		return true
	}
	s.mu.Lock()
	v := s.rng.Float64()
	s.mu.Unlock()
	return v < rate
}
