// Package router matches log records to configured sink destinations
// based on source name and optional field-based routing rules.
package router

import (
	"fmt"
	"regexp"
)

// Rule defines a single routing rule mapping a source pattern and
// optional field match to a list of sink names.
type Rule struct {
	SourcePattern string
	FieldKey      string
	FieldPattern  string
	Sinks         []string

	sourceRe *regexp.Regexp
	fieldRe  *regexp.Regexp
}

// Router holds compiled routing rules and resolves sink targets for records.
type Router struct {
	rules        []*Rule
	defaultSinks []string
}

// New compiles the provided rules and returns a Router.
// defaultSinks are used when no rule matches.
func New(rules []Rule, defaultSinks []string) (*Router, error) {
	compiled := make([]*Rule, 0, len(rules))
	for i, r := range rules {
		rc := r
		if rc.SourcePattern == "" {
			rc.SourcePattern = ".*"
		}
		sre, err := regexp.Compile(rc.SourcePattern)
		if err != nil {
			return nil, fmt.Errorf("router: rule %d invalid source pattern: %w", i, err)
		}
		rc.sourceRe = sre
		if rc.FieldPattern != "" {
			fre, err := regexp.Compile(rc.FieldPattern)
			if err != nil {
				return nil, fmt.Errorf("router: rule %d invalid field pattern: %w", i, err)
			}
			rc.fieldRe = fre
		}
		compiled = append(compiled, &rc)
	}
	return &Router{rules: compiled, defaultSinks: defaultSinks}, nil
}

// Resolve returns the sink names that should receive the record.
// source is the origin identifier; fields are the parsed log fields.
func (r *Router) Resolve(source string, fields map[string]string) []string {
	for _, rule := range r.rules {
		if !rule.sourceRe.MatchString(source) {
			continue
		}
		if rule.fieldRe != nil {
			val := fields[rule.FieldKey]
			if !rule.fieldRe.MatchString(val) {
				continue
			}
		}
		return rule.Sinks
	}
	return r.defaultSinks
}
