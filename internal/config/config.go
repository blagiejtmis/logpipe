package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// SourceConfig defines a log source to tail.
type SourceConfig struct {
	Name   string `yaml:"name"`
	Type   string `yaml:"type"`   // file, stdin, journald
	Path   string `yaml:"path"`   // for file sources
	Format string `yaml:"format"` // json, logfmt, plain
}

// SinkConfig defines a destination for log records.
type SinkConfig struct {
	Name   string            `yaml:"name"`
	Type   string            `yaml:"type"`   // stdout, file, http
	Target string            `yaml:"target"` // path or URL
	Opts   map[string]string `yaml:"opts"`
}

// Config is the top-level logpipe configuration.
type Config struct {
	Sources []SourceConfig `yaml:"sources"`
	Sinks   []SinkConfig   `yaml:"sinks"`
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	var cfg Config
	dec := yaml.NewDecoder(f)
	dec.KnownFields(true)
	if err := dec.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("config: decode %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config: validation: %w", err)
	}

	return &cfg, nil
}

// SourceByName returns the SourceConfig with the given name, or an error if not found.
func (c *Config) SourceByName(name string) (*SourceConfig, error) {
	for i := range c.Sources {
		if c.Sources[i].Name == name {
			return &c.Sources[i], nil
		}
	}
	return nil, fmt.Errorf("config: no source with name %q", name)
}

// SinkByName returns the SinkConfig with the given name, or an error if not found.
func (c *Config) SinkByName(name string) (*SinkConfig, error) {
	for i := range c.Sinks {
		if c.Sinks[i].Name == name {
			return &c.Sinks[i], nil
		}
	}
	return nil, fmt.Errorf("config: no sink with name %q", name)
}

func (c *Config) validate() error {
	if len(c.Sources) == 0 {
		return fmt.Errorf("at least one source is required")
	}
	if len(c.Sinks) == 0 {
		return fmt.Errorf("at least one sink is required")
	}
	for i, s := range c.Sources {
		if s.Name == "" {
			return fmt.Errorf("source[%d]: name is required", i)
		}
		if s.Type == "" {
			return fmt.Errorf("source %q: type is required", s.Name)
		}
	}
	for i, s := range c.Sinks {
		if s.Name == "" {
			return fmt.Errorf("sink[%d]: name is required", i)
		}
		if s.Type == "" {
			return fmt.Errorf("sink %q: type is required", s.Name)
		}
	}
	return nil
}
