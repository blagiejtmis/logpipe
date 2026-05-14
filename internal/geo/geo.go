// Package geo provides IP-to-geographic-location enrichment for log records.
// Rules specify a source field containing an IP address and a destination field
// where a map of location attributes (country, region, city) will be written.
package geo

import (
	"errors"
	"fmt"
	"net"
)

// Lookup maps an IP string to location attributes.
type Lookup func(ip string) (map[string]string, bool)

// Rule describes a single geo-enrichment operation.
type Rule struct {
	// SrcField is the record field that holds the IP address string.
	SrcField string
	// DstField is the record field where location attributes are written.
	DstField string
	// OnMiss controls behaviour when the IP is not found: "skip" (default) or "empty".
	OnMiss string
}

// Enricher applies geo-lookup rules to log records.
type Enricher struct {
	rules  []Rule
	lookup Lookup
}

// New creates an Enricher. lookup must not be nil and at least one rule is required.
func New(rules []Rule, lookup Lookup) (*Enricher, error) {
	if lookup == nil {
		return nil, errors.New("geo: lookup function must not be nil")
	}
	if len(rules) == 0 {
		return nil, errors.New("geo: at least one rule is required")
	}
	for i, r := range rules {
		if r.SrcField == "" {
			return nil, fmt.Errorf("geo: rule[%d]: src_field must not be empty", i)
		}
		if r.DstField == "" {
			return nil, fmt.Errorf("geo: rule[%d]: dst_field must not be empty", i)
		}
		if r.OnMiss != "" && r.OnMiss != "skip" && r.OnMiss != "empty" {
			return nil, fmt.Errorf("geo: rule[%d]: on_miss must be \"skip\" or \"empty\"", i)
		}
	}
	return &Enricher{rules: rules, lookup: lookup}, nil
}

// Apply enriches rec in-place according to the configured rules.
func (e *Enricher) Apply(rec map[string]any) map[string]any {
	for _, r := range e.rules {
		ipRaw, ok := rec[r.SrcField]
		if !ok {
			continue
		}
		ipStr, ok := ipRaw.(string)
		if !ok {
			continue
		}
		if net.ParseIP(ipStr) == nil {
			continue
		}
		attrs, found := e.lookup(ipStr)
		if !found {
			onMiss := r.OnMiss
			if onMiss == "" {
				onMiss = "skip"
			}
			if onMiss == "empty" {
				rec[r.DstField] = map[string]string{}
			}
			continue
		}
		out := make(map[string]string, len(attrs))
		for k, v := range attrs {
			out[k] = v
		}
		rec[r.DstField] = out
	}
	return rec
}
