package config

// RedactRule describes a single field-level redaction rule.
type RedactRule struct {
	// Field is the log record key to inspect.
	Field string `yaml:"field"`
	// Pattern is a regular expression matched against the field value.
	// A full match triggers redaction.
	Pattern string `yaml:"pattern"`
	// Placeholder replaces the matched value. Defaults to "***" when empty.
	Placeholder string `yaml:"placeholder,omitempty"`
}

// RedactConfig groups global and per-source redaction rules.
type RedactConfig struct {
	// Global rules apply to every source unless overridden.
	Global []RedactRule `yaml:"global,omitempty"`
	// Sources maps source identifiers to their specific rule sets.
	// A source-specific set fully replaces the global rules for that source.
	Sources map[string][]RedactRule `yaml:"sources,omitempty"`
}
