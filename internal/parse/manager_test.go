package parse

import (
	"testing"

	"github.com/logpipe/logpipe/internal/config"
)

func makeParseCfg(defaultFmt string, sources map[string]string) *config.ParseConfig {
	return &config.ParseConfig{
		DefaultFormat: defaultFmt,
		Sources:       sources,
	}
}

func TestNewManager_NilConfig_UsesJSONDefault(t *testing.T) {
	m, err := NewManager(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p := m.ParserFor("any-source")
	if p == nil {
		t.Fatal("expected non-nil parser")
	}
	// JSON parser should handle a valid JSON line
	rec, err := p.Parse(`{"msg":"hello"}`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if rec["msg"] != "hello" {
		t.Errorf("expected msg=hello, got %v", rec["msg"])
	}
}

func TestNewManager_DefaultFormat_Applied(t *testing.T) {
	m, err := NewManager(makeParseCfg("logfmt", nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p := m.ParserFor("src")
	rec, err := p.Parse(`level=info msg=started`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if rec["level"] != "info" {
		t.Errorf("expected level=info, got %v", rec["level"])
	}
}

func TestNewManager_SourceSpecific_OverridesDefault(t *testing.T) {
	m, err := NewManager(makeParseCfg("json", map[string]string{
		"app.log": "logfmt",
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// source-specific parser should be logfmt
	p := m.ParserFor("app.log")
	rec, err := p.Parse(`level=warn msg=oops`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if rec["msg"] != "oops" {
		t.Errorf("expected msg=oops, got %v", rec["msg"])
	}

	// unknown source falls back to default (json)
	pDefault := m.ParserFor("other.log")
	rec2, err := pDefault.Parse(`{"x":"y"}`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if rec2["x"] != "y" {
		t.Errorf("expected x=y, got %v", rec2["x"])
	}
}

func TestNewManager_InvalidDefaultFormat_ReturnsError(t *testing.T) {
	_, err := NewManager(makeParseCfg("xml", nil))
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestNewManager_InvalidSourceFormat_ReturnsError(t *testing.T) {
	_, err := NewManager(makeParseCfg("json", map[string]string{
		"bad.log": "csv",
	}))
	if err == nil {
		t.Fatal("expected error for unknown source format")
	}
}
