package parse

import (
	"strings"
	"testing"
)

func TestNew_UnknownFormat_ReturnsError(t *testing.T) {
	_, err := New("xml")
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestNew_JSONFormat_ReturnsParser(t *testing.T) {
	p, err := New("json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil parser")
	}
}

func TestJSONParser_ValidLine(t *testing.T) {
	p, _ := New("json")
	rec, err := p.Parse(`{"level":"info","msg":"hello","time":"2024-01-01T00:00:00Z"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec["level"] != "info" {
		t.Errorf("expected level=info, got %q", rec["level"])
	}
	if rec["msg"] != "hello" {
		t.Errorf("expected msg=hello, got %q", rec["msg"])
	}
}

func TestJSONParser_InjectsTimeIfMissing(t *testing.T) {
	p, _ := New("json")
	rec, err := p.Parse(`{"msg":"no time here"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec["time"] == "" {
		t.Error("expected time to be injected")
	}
}

func TestJSONParser_InvalidLine_ReturnsError(t *testing.T) {
	p, _ := New("json")
	_, err := p.Parse("not json")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestLogfmtParser_ValidLine(t *testing.T) {
	p, _ := New("logfmt")
	rec, err := p.Parse(`level=info msg="hello world" code=200`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec["level"] != "info" {
		t.Errorf("expected level=info, got %q", rec["level"])
	}
	if rec["code"] != "200" {
		t.Errorf("expected code=200, got %q", rec["code"])
	}
}

func TestLogfmtParser_BareToken_StoredAsMsg(t *testing.T) {
	p, _ := New("logfmt")
	rec, err := p.Parse("ERROR")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(rec["msg"], "ERROR") {
		t.Errorf("expected msg to contain ERROR, got %q", rec["msg"])
	}
}

func TestLogfmtParser_EmptyLine_ReturnsError(t *testing.T) {
	p, _ := New("logfmt")
	_, err := p.Parse("")
	if err == nil {
		t.Fatal("expected error for empty line")
	}
}

func TestLogfmtParser_InjectsTimeIfMissing(t *testing.T) {
	p, _ := New("logfmt")
	rec, err := p.Parse("level=warn msg=test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec["time"] == "" {
		t.Error("expected time to be injected")
	}
}
