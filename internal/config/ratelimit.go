package config

// RateLimitConfig holds the top-level rate-limiting configuration block.
type RateLimitConfig struct {
	// Default applies to any source that does not have its own entry.
	// A zero Rate disables the default limit.
	Default RateLimitEntry `yaml:"default"`

	// Sources lists per-source overrides.
	Sources []SourceRateLimit `yaml:"sources"`
}

// RateLimitEntry is a single rate/window pair.
type RateLimitEntry struct {
	// Rate is the maximum number of log lines allowed per WindowSecs.
	Rate int `yaml:"rate"`

	// WindowSecs is the sliding-window size in seconds (must be >= 1).
	WindowSecs int `yaml:"window_secs"`
}

// SourceRateLimit binds a rate limit to a specific source path or glob.
type SourceRateLimit struct {
	// Source is the file path (or glob) this limit applies to.
	Source string `yaml:"source"`

	RateLimitEntry `yaml:",inline"`
}
