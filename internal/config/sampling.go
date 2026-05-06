package config

// SamplingConfig controls probabilistic sampling of log records.
// DefaultRate applies to all sources not listed in Sources.
// A rate of 1.0 passes all records; 0.1 passes ~10%.
// A zero DefaultRate means no global sampler is configured (all pass).
type SamplingConfig struct {
	// DefaultRate is the fallback sampling rate for unlisted sources.
	// Valid range: (0.0, 1.0]. Zero disables the global sampler.
	DefaultRate float64 `yaml:"default_rate"`

	// Sources maps source identifiers to their individual sampling rates.
	Sources map[string]float64 `yaml:"sources"`
}
