package router

import (
	"testing"
)

func TestRouter_DefaultSinks_WhenNoRules(t *testing.T) {
	r, err := New(nil, []string{"stdout"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sinks := r.Resolve("app.log", map[string]string{"level": "info"})
	if len(sinks) != 1 || sinks[0] != "stdout" {
		t.Fatalf("expected [stdout], got %v", sinks)
	}
}

func TestRouter_MatchBySource(t *testing.T) {
	rules := []Rule{
		{SourcePattern: "^app\\.log$", Sinks: []string{"file-app"}},
	}
	r, err := New(rules, []string{"stdout"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sinks := r.Resolve("app.log", nil); len(sinks) != 1 || sinks[0] != "file-app" {
		t.Fatalf("expected [file-app], got %v", sinks)
	}
	if sinks := r.Resolve("other.log", nil); len(sinks) != 1 || sinks[0] != "stdout" {
		t.Fatalf("expected [stdout], got %v", sinks)
	}
}

func TestRouter_MatchByField(t *testing.T) {
	rules := []Rule{
		{SourcePattern: ".*", FieldKey: "level", FieldPattern: "^error$", Sinks: []string{"errors-sink"}},
	}
	r, err := New(rules, []string{"stdout"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fields := map[string]string{"level": "error"}
	if sinks := r.Resolve("any.log", fields); len(sinks) != 1 || sinks[0] != "errors-sink" {
		t.Fatalf("expected [errors-sink], got %v", sinks)
	}
	fields["level"] = "info"
	if sinks := r.Resolve("any.log", fields); len(sinks) != 1 || sinks[0] != "stdout" {
		t.Fatalf("expected [stdout] for info, got %v", sinks)
	}
}

func TestRouter_MultipleSinks(t *testing.T) {
	rules := []Rule{
		{SourcePattern: "svc", Sinks: []string{"file-svc", "stdout"}},
	}
	r, err := New(rules, []string{"stdout"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sinks := r.Resolve("svc.log", nil)
	if len(sinks) != 2 {
		t.Fatalf("expected 2 sinks, got %v", sinks)
	}
}

func TestRouter_InvalidSourcePattern(t *testing.T) {
	rules := []Rule{
		{SourcePattern: "[invalid", Sinks: []string{"stdout"}},
	}
	_, err := New(rules, nil)
	if err == nil {
		t.Fatal("expected error for invalid source pattern")
	}
}

func TestRouter_InvalidFieldPattern(t *testing.T) {
	rules := []Rule{
		{FieldKey: "level", FieldPattern: "[bad", Sinks: []string{"stdout"}},
	}
	_, err := New(rules, nil)
	if err == nil {
		t.Fatal("expected error for invalid field pattern")
	}
}
