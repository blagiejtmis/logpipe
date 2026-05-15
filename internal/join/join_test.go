package join

import (
	"testing"
	"time"
)

func validRule() Rule {
	return Rule{
		LeftSource:  "app",
		RightSource: "db",
		KeyField:    "request_id",
		Window:      100 * time.Millisecond,
	}
}

func TestNew_EmptyLeftSource_ReturnsError(t *testing.T) {
	r := validRule()
	r.LeftSource = ""
	_, err := New(r)
	if err == nil {
		t.Fatal("expected error for empty LeftSource")
	}
}

func TestNew_EmptyRightSource_ReturnsError(t *testing.T) {
	r := validRule()
	r.RightSource = ""
	_, err := New(r)
	if err == nil {
		t.Fatal("expected error for empty RightSource")
	}
}

func TestNew_EmptyKeyField_ReturnsError(t *testing.T) {
	r := validRule()
	r.KeyField = ""
	_, err := New(r)
	if err == nil {
		t.Fatal("expected error for empty KeyField")
	}
}

func TestNew_ZeroWindow_ReturnsError(t *testing.T) {
	r := validRule()
	r.Window = 0
	_, err := New(r)
	if err == nil {
		t.Fatal("expected error for zero Window")
	}
}

func TestNew_ValidRule_NoError(t *testing.T) {
	_, err := New(validRule())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAdd_LeftThenRight_Joins(t *testing.T) {
	j, _ := New(validRule())
	left := Record{"request_id": "abc", "method": "GET"}
	right := Record{"request_id": "abc", "status": 200}

	_, ok := j.Add("app", left)
	if ok {
		t.Fatal("expected no join on first add")
	}
	out, ok := j.Add("db", right)
	if !ok {
		t.Fatal("expected join on second add")
	}
	if out["method"] != "GET" {
		t.Errorf("missing left field: %v", out)
	}
	if out["status"] != 200 {
		t.Errorf("missing right field: %v", out)
	}
}

func TestAdd_RightThenLeft_Joins(t *testing.T) {
	j, _ := New(validRule())
	left := Record{"request_id": "xyz", "user": "alice"}
	right := Record{"request_id": "xyz", "latency": 42}

	j.Add("db", right)
	out, ok := j.Add("app", left)
	if !ok {
		t.Fatal("expected join")
	}
	if out["user"] != "alice" || out["latency"] != 42 {
		t.Errorf("unexpected merged record: %v", out)
	}
}

func TestAdd_OutputField_NestsRightRecord(t *testing.T) {
	r := validRule()
	r.OutputField = "db_info"
	j, _ := New(r)

	j.Add("app", Record{"request_id": "1", "method": "POST"})
	out, ok := j.Add("db", Record{"request_id": "1", "rows": 5})
	if !ok {
		t.Fatal("expected join")
	}
	dbInfo, ok := out["db_info"].(Record)
	if !ok {
		t.Fatalf("expected db_info to be a Record, got %T", out["db_info"])
	}
	if dbInfo["rows"] != 5 {
		t.Errorf("expected rows=5, got %v", dbInfo["rows"])
	}
}

func TestAdd_MissingKeyField_Ignored(t *testing.T) {
	j, _ := New(validRule())
	_, ok := j.Add("app", Record{"method": "GET"})
	if ok {
		t.Fatal("expected no join when key field missing")
	}
}

func TestAdd_WindowExpiry_DoesNotJoin(t *testing.T) {
	r := validRule()
	r.Window = 1 * time.Millisecond
	j, _ := New(r)

	j.Add("app", Record{"request_id": "exp", "x": 1})
	time.Sleep(10 * time.Millisecond)
	_, ok := j.Add("db", Record{"request_id": "exp", "y": 2})
	if ok {
		t.Fatal("expected no join after window expiry")
	}
}

func TestAdd_DifferentKeys_NoJoin(t *testing.T) {
	j, _ := New(validRule())
	j.Add("app", Record{"request_id": "aaa", "a": 1})
	_, ok := j.Add("db", Record{"request_id": "bbb", "b": 2})
	if ok {
		t.Fatal("expected no join for different keys")
	}
}
