package router

import (
	"fmt"

	"github.com/user/logpipe/internal/config"
)

// RuleConfig mirrors the config structure for a routing rule.
type RuleConfig struct {
	Source  string   `yaml:"source"`
	Field   string   `yaml:"field"`
	Match   string   `yaml:"match"`
	Sinks   []string `yaml:"sinks"`
}

// NewFromConfig builds a Router from the application configuration.
// cfg.Routing.Rules lists per-rule overrides; cfg.Routing.DefaultSinks
// provides the fallback sink list.
func NewFromConfig(cfg *config.Config) (*Router, error) {
	var rules []Rule
	for i, rc := range cfg.Routing.Rules {
		if len(rc.Sinks) == 0 {
			return nil, fmt.Errorf("router: rule %d has no sinks", i)
		}
		if rc.Source == "" && rc.Field == "" {
			return nil, fmt.Errorf("router: rule %d must specify at least one of 'source' or 'field'", i)
		}
		rules = append(rules, Rule{
			SourcePattern: rc.Source,
			FieldKey:      rc.Field,
			FieldPattern:  rc.Match,
			Sinks:         rc.Sinks,
		})
	}
	defaults := cfg.Routing.DefaultSinks
	if len(defaults) == 0 {
		// Fall back to all declared sink names.
		for _, s := range cfg.Sinks {
			defaults = append(defaults, s.Name)
		}
	}
	return New(rules, defaults)
}
