package multiline

import (
	"testing"
	"time"
)

func TestNew_NoPatterns_ReturnsError(t *testing.T) {
	_, err := New(Rule{})
	if err == nil {
		t.Fatal("expected error for empty patterns")
	}
}

func TestNew_BothPatterns_ReturnsError(t *testing.T) {
	_, err := New(Rule{StartPattern: "^START", ContinuePattern: `^\s`})
	if err == nil {
		t.Fatal("expected error when both patterns set")
	}
}

func TestNew_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := New(Rule{StartPattern: "[invalid"})
	if err == nil {
		t.Fatal("expected regexp compile error")
	}
}

func TestNew_NegativeMaxLines_ReturnsError(t *testing.T) {
	_, err := New(Rule{StartPattern: "^START", MaxLines: -1})
	if err == nil {
		t.Fatal("expected error for negative MaxLines")
	}
}

func TestAdd_StartPattern_BuffersUntilNextStart(t *testing.T) {
	a, err := New(Rule{StartPattern: `^\d{4}-`})
	if err != nil {
		t.Fatal(err)
	}
	if out := a.Add("2024-01-01 ERROR something"); out != nil {
		t.Fatalf("unexpected flush on first line: %v", out)
	}
	if out := a.Add("  at foo.go:10"); out != nil {
		t.Fatalf("unexpected flush on continuation: %v", out)
	}
	out := a.Add("2024-01-01 INFO next")
	if out == nil {
		t.Fatal("expected flush when new start line arrived")
	}
	if out["message"] != "2024-01-01 ERROR something\n  at foo.go:10" {
		t.Errorf("unexpected message: %q", out["message"])
	}
}

func TestAdd_ContinuePattern_FlushesOnNonMatch(t *testing.T) {
	a, err := New(Rule{ContinuePattern: `^\s`, Field: "msg"})
	if err != nil {
		t.Fatal(err)
	}
	a.Add("first line")
	a.Add("  indented")
	out := a.Add("new record")
	if out == nil {
		t.Fatal("expected flush")
	}
	if out["msg"] != "first line\n  indented" {
		t.Errorf("unexpected msg: %q", out["msg"])
	}
}

func TestAdd_MaxLines_ForcesFlush(t *testing.T) {
	a, err := New(Rule{ContinuePattern: `^\s`, MaxLines: 2})
	if err != nil {
		t.Fatal(err)
	}
	a.Add("line1")
	out := a.Add(" line2")
	if out == nil {
		t.Fatal("expected flush at MaxLines")
	}
}

func TestAdd_Timeout_ForcesFlush(t *testing.T) {
	a, err := New(Rule{ContinuePattern: `^\s`, Timeout: 1 * time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}
	a.Add("line1")
	time.Sleep(5 * time.Millisecond)
	out := a.Add(" line2")
	if out == nil {
		t.Fatal("expected flush after timeout")
	}
}

func TestFlush_EmitsBuffered(t *testing.T) {
	a, err := New(Rule{StartPattern: `^START`})
	if err != nil {
		t.Fatal(err)
	}
	a.Add("START record")
	a.Add("  detail")
	out := a.Flush()
	if out == nil {
		t.Fatal("expected flush")
	}
	if out["message"] != "START record\n  detail" {
		t.Errorf("unexpected message: %q", out["message"])
	}
}

func TestFlush_EmptyBuffer_ReturnsNil(t *testing.T) {
	a, _ := New(Rule{StartPattern: `^START`})
	if out := a.Flush(); out != nil {
		t.Fatalf("expected nil for empty buffer, got %v", out)
	}
}
