// Package parse provides log line parsers that convert raw strings into
// structured records for the logpipe pipeline.
package parse

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Record is a structured log entry produced by a parser.
type Record map[string]string

// Parser converts a raw log line into a Record.
type Parser interface {
	Parse(line string) (Record, error)
}

// New returns a Parser for the given format.
// Supported formats: "json", "logfmt".
func New(format string) (Parser, error) {
	switch strings.ToLower(format) {
	case "json":
		return &jsonParser{}, nil
	case "logfmt":
		return &logfmtParser{}, nil
	default:
		return nil, fmt.Errorf("parse: unsupported format %q", format)
	}
}

// jsonParser parses JSON log lines.
type jsonParser struct{}

func (p *jsonParser) Parse(line string) (Record, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		return nil, fmt.Errorf("parse: invalid JSON: %w", err)
	}
	rec := make(Record, len(raw))
	for k, v := range raw {
		switch val := v.(type) {
		case string:
			rec[k] = val
		default:
			rec[k] = fmt.Sprintf("%v", v)
		}
	}
	if _, ok := rec["time"]; !ok {
		rec["time"] = time.Now().UTC().Format(time.RFC3339)
	}
	return rec, nil
}

// logfmtParser parses key=value log lines.
type logfmtParser struct{}

func (p *logfmtParser) Parse(line string) (Record, error) {
	rec := make(Record)
	for _, token := range strings.Fields(line) {
		parts := strings.SplitN(token, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			val := strings.Trim(parts[1], `"`)
			rec[key] = val
		} else {
			// bare token stored under "msg" if no key yet
			if _, ok := rec["msg"]; !ok {
				rec["msg"] = parts[0]
			}
		}
	}
	if len(rec) == 0 {
		return nil, fmt.Errorf("parse: empty logfmt line")
	}
	if _, ok := rec["time"]; !ok {
		rec["time"] = time.Now().UTC().Format(time.RFC3339)
	}
	return rec, nil
}
