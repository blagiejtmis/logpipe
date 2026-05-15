package unwrap_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/unwrap"
)

func base() map[string]any {
	return map[string]any{
		"message": "hello",
		"meta": map[string]any{
			"host": "srv1",
			"env":  "prod",
		},
	}
}

func TestNew_NoRules_ReturnsError(t *testing.T) {
	_, err := unwrap.New(nil)
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestNew_BlankField_ReturnsError(t *testing.T) {
	_, err := unwrap.New([]unwrap.Rule{{Field: "  "}})
	if err == nil {
		t.Fatal("expected error for blank field")
	}
}

func TestNew_ValidRules_NoError(t *testing.T) {
	_, err := unwrap.New([]unwrap.Rule{{Field: "meta"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApply_PromotesKeys(t *testing.T) {
	u, _ := unwrap.New([]unwrap.Rule{{Field: "meta", Remove: true}})
	rec := base()
	out := u.Apply(rec)
	if out["host"] != "srv1" {
		t.Errorf("expected host=srv1, got %v", out["host"])
	}
	if _, ok := out["meta"]; ok {
		t.Error("expected meta to be removed")
	}
}

func TestApply_WithPrefix(t *testing.T) {
	u, _ := unwrap.New([]unwrap.Rule{{Field: "meta", Prefix: "meta_", Remove: false}})
	out := u.Apply(base())
	if out["meta_host"] != "srv1" {
		t.Errorf("expected meta_host=srv1, got %v", out["meta_host"])
	}
	if _, ok := out["meta"]; !ok {
		t.Error("expected meta to be preserved when Remove=false")
	}
}

func TestApply_MissingField_IsSkipped(t *testing.T) {
	u, _ := unwrap.New([]unwrap.Rule{{Field: "missing"}})
	out := u.Apply(base())
	if out["message"] != "hello" {
		t.Error("record should be unchanged when field is absent")
	}
}

func TestApply_NonMapField_IsSkipped(t *testing.T) {
	u, _ := unwrap.New([]unwrap.Rule{{Field: "message"}})
	out := u.Apply(base())
	if out["message"] != "hello" {
		t.Error("record should be unchanged when field is not a map")
	}
}
