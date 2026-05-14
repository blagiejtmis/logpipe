// Package multiline provides multi-line log record assembly.
// It buffers consecutive log lines that belong to the same logical record
// (e.g. stack traces, multi-line exceptions) and emits them as a single entry.
package multiline

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

// Rule describes how to detect continuation lines.
type Rule struct {
	// StartPattern marks the beginning of a new record.
	StartPattern string
	// ContinuePattern marks a continuation line (mutually exclusive with StartPattern).
	ContinuePattern string
	// MaxLines is the maximum number of lines to buffer (0 = unlimited).
	MaxLines int
	// Timeout is the maximum time to wait for more lines before flushing.
	Timeout time.Duration
	// Field is the log field to assemble into (default: "message").
	Field string
}

// Assembler buffers lines and emits complete multi-line records.
type Assembler struct {
	rule    Rule
	start   *regexp.Regexp
	cont    *regexp.Regexp
	buf     []string
	lastAt  time.Time
}

// New creates an Assembler from the given Rule.
func New(r Rule) (*Assembler, error) {
	if r.StartPattern == "" && r.ContinuePattern == "" {
		return nil, errors.New("multiline: at least one of StartPattern or ContinuePattern must be set")
	}
	if r.StartPattern != "" && r.ContinuePattern != "" {
		return nil, errors.New("multiline: StartPattern and ContinuePattern are mutually exclusive")
	}
	if r.Field == "" {
		r.Field = "message"
	}
	if r.MaxLines < 0 {
		return nil, errors.New("multiline: MaxLines must be >= 0")
	}
	a := &Assembler{rule: r}
	var err error
	if r.StartPattern != "" {
		if a.start, err = regexp.Compile(r.StartPattern); err != nil {
			return nil, err
		}
	}
	if r.ContinuePattern != "" {
		if a.cont, err = regexp.Compile(r.ContinuePattern); err != nil {
			return nil, err
		}
	}
	return a, nil
}

// Add feeds a raw line into the assembler.
// It returns a flushed record (field -> joined text) when a complete record is
// ready, or nil if more lines are needed.
func (a *Assembler) Add(line string) map[string]string {
	now := time.Now()
	var flush bool

	if a.start != nil {
		// start-pattern mode: new record begins when pattern matches
		if a.start.MatchString(line) && len(a.buf) > 0 {
			flush = true
		}
	} else {
		// continue-pattern mode: flush when line does NOT match
		if len(a.buf) > 0 && !a.cont.MatchString(line) {
			flush = true
		}
	}

	if a.rule.MaxLines > 0 && len(a.buf) >= a.rule.MaxLines {
		flush = true
	}
	if a.rule.Timeout > 0 && !a.lastAt.IsZero() && now.Sub(a.lastAt) >= a.rule.Timeout {
		flush = true
	}

	var out map[string]string
	if flush {
		out = a.emit()
	}
	a.buf = append(a.buf, line)
	a.lastAt = now
	return out
}

// Flush forces emission of whatever is buffered.
func (a *Assembler) Flush() map[string]string {
	if len(a.buf) == 0 {
		return nil
	}
	return a.emit()
}

func (a *Assembler) emit() map[string]string {
	joined := strings.Join(a.buf, "\n")
	a.buf = a.buf[:0]
	a.lastAt = time.Time{}
	return map[string]string{a.rule.Field: joined}
}
