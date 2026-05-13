package config

// SequenceConfig configures monotonic sequence stamping for log records.
// A global Default rule applies to all sources unless overridden by a
// source-specific entry in Sources.
type SequenceConfig struct {
	// Default is applied to any source not listed in Sources.
	Default *SequenceRule `yaml:"default,omitempty"`

	// Sources maps source names to per-source sequencing rules.
	Sources map[string]*SequenceRule `yaml:"sources,omitempty"`
}

// SequenceRule configures a single sequencer instance.
type SequenceRule struct {
	// Field is the record key that will receive the sequence number.
	// Defaults to "seq" if empty.
	Field string `yaml:"field"`

	// Start is the first sequence number emitted. Defaults to 1.
	Start uint64 `yaml:"start"`
}
