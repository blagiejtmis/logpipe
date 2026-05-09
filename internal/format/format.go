// Package format provides log record formatting utilities for rendering
// structured records into human-readable or machine-readable strings.
package format

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
)

// Record represents a structured log entry passed through the pipeline.
type Record = map[string]any

// Formatter converts a Record into a string representation.
type Formatter interface {
	Format(r Record) (string, error)
}

// New returns a Formatter for the given format name.
// Supported formats: "json", "text", "logfmt".
func New(format string) (Formatter, error) {
	switch strings.ToLower(format) {
	case "json":
		return &jsonFormatter{}, nil
	case "text":
		return &textFormatter{}, nil
	case "logfmt":
		return &logfmtFormatter{}, nil
	default:
		return nil, fmt.Errorf("format: unknown format %q", format)
	}
}

// jsonFormatter renders records as compact JSON.
type jsonFormatter struct{}

func (f *jsonFormatter) Format(r Record) (string, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return "", fmt.Errorf("format: json marshal: %w", err)
	}
	return string(b), nil
}

// textFormatter renders records as a human-readable line:
// TIMESTAMP LEVEL message key=value ...
type textFormatter struct{}

func (f *textFormatter) Format(r Record) (string, error) {
	ts := ""
	if v, ok := r["time"]; ok {
		switch t := v.(type) {
		case time.Time:
			ts = t.Format(time.RFC3339)
		case string:
			ts = t
		}
	}
	level, _ := r["level"].(string)
	msg, _ := r["message"].(string)

	var extras []string
	skip := map[string]bool{"time": true, "level": true, "message": true}
	keys := make([]string, 0, len(r))
	for k := range r {
		if !skip[k] {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	for _, k := range keys {
		extras = append(extras, fmt.Sprintf("%s=%v", k, r[k]))
	}

	parts := []string{ts, strings.ToUpper(level), msg}
	if len(extras) > 0 {
		parts = append(parts, extras...)
	}
	return strings.TrimSpace(strings.Join(parts, " ")), nil
}

// logfmtFormatter renders records in logfmt style: key=value pairs.
type logfmtFormatter struct{}

func (f *logfmtFormatter) Format(r Record) (string, error) {
	keys := make([]string, 0, len(r))
	for k := range r {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var pairs []string
	for _, k := range keys {
		v := fmt.Sprintf("%v", r[k])
		if strings.ContainsAny(v, " \t\n") {
			v = fmt.Sprintf("%q", v)
		}
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(pairs, " "), nil
}
