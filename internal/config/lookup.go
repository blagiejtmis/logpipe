package config

// LookupRule defines a single table-lookup enrichment step.
type LookupRule struct {
	// KeyField is the record field whose value is used as the lookup key.
	KeyField string `yaml:"key_field"`

	// Table maps key values to a map of fields that will be merged into the
	// record when the key matches.
	Table map[string]map[string]any `yaml:"table"`

	// DestField is written with the matched value when the table returns a
	// single scalar. When the table row contains multiple keys all of them
	// are merged and DestField is ignored.
	DestField string `yaml:"dest_field"`

	// OnMiss controls behaviour when the key is not found in the table.
	// Accepted values: "keep" (leave record unchanged) or "drop" (drop the
	// record). Defaults to "keep".
	OnMiss string `yaml:"on_miss"`
}

// LookupConfig groups global and per-source lookup rules.
type LookupConfig struct {
	// Global rules are applied to every record regardless of source.
	Global []LookupRule `yaml:"global"`

	// Sources maps a source name to a slice of rules that override Global for
	// that source.
	Sources map[string][]LookupRule `yaml:"sources"`
}
