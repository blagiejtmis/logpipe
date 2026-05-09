package format_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logpipe/internal/format"
)

func TestNew_UnknownFormat_ReturnsError(t *testing.T) {
	_, err := format.New("xml")
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestNew_KnownFormats_NoError(t *testing.T) {
	for _, f := range []string{"json", "text", "logfmt", "JSON", "TEXT"} {
		_, err := format.New(f)
		if err != nil {
			t.Errorf("format %q: unexpected error: %v", f, err)
		}
	}
}

func TestJSONFormatter_ValidRecord(t *testing.T) {
	f, _ := format.New("json")
	out, err := f.Format(format.Record{"level": "info", "message": "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "\"level\":\"info\"") {
		t.Errorf("expected level in output, got: %s", out)
	}
}

func TestTextFormatter_ContainsLevelAndMessage(t *testing.T) {
	f, _ := format.New("text")
	out, err := f.Format(format.Record{
		"level":   "warn",
		"message": "disk full",
		"host":    "srv1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "WARN") {
		t.Errorf("expected WARN in output, got: %s", out)
	}
	if !strings.Contains(out, "disk full") {
		t.Errorf("expected message in output, got: %s", out)
	}
	if !strings.Contains(out, "host=srv1") {
		t.Errorf("expected host=srv1 in output, got: %s", out)
	}
}

func TestTextFormatter_EmptyRecord(t *testing.T) {
	f, _ := format.New("text")
	out, err := f.Format(format.Record{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should not panic; output may be empty or whitespace only.
	_ = out
}

func TestLogfmtFormatter_KeyValuePairs(t *testing.T) {
	f, _ := format.New("logfmt")
	out, err := f.Format(format.Record{
		"level":   "error",
		"message": "oops",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "level=error") {
		t.Errorf("expected level=error in output, got: %s", out)
	}
	if !strings.Contains(out, "message=oops") {
		t.Errorf("expected message=oops in output, got: %s", out)
	}
}

func TestLogfmtFormatter_QuotesSpacedValues(t *testing.T) {
	f, _ := format.New("logfmt")
	out, err := f.Format(format.Record{"msg": "hello world"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"`) {
		t.Errorf("expected quoted value for spaced string, got: %s", out)
	}
}
