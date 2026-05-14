package tag_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/tag"
)

func base() map[string]any {
	return map[string]any{"message": "hello", "level": "info"}
}

func TestNew_NoRules_ReturnsError(t *testing.T) {
	_, err := tag.New(nil)
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestNew_EmptyField_ReturnsError(t *testing.T) {
	_, err := tag.New([]tag.Rule{{Field: "", Value: "v"}})
	if err == nil {
		t.Fatal("expected error for blank field")
	}
}

func TestNew_ValidRules_NoError(t *testing.T) {
	_, err := tag.New([]tag.Rule{{Field: "env", Value: "prod"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApply_AddsField(t *testing.T) {
	tgr, _ := tag.New([]tag.Rule{{Field: "env", Value: "prod"}})
	rec := base()
	out := tgr.Apply(rec)
	if out["env"] != "prod" {
		t.Fatalf("expected env=prod, got %v", out["env"])
	}
}

func TestApply_DoesNotOverwriteByDefault(t *testing.T) {
	tgr, _ := tag.New([]tag.Rule{{Field: "level", Value: "debug", Overwrite: false}})
	rec := base()
	out := tgr.Apply(rec)
	if out["level"] != "info" {
		t.Fatalf("expected level to remain 'info', got %v", out["level"])
	}
}

func TestApply_OverwriteExistingField(t *testing.T) {
	tgr, _ := tag.New([]tag.Rule{{Field: "level", Value: "warn", Overwrite: true}})
	rec := base()
	out := tgr.Apply(rec)
	if out["level"] != "warn" {
		t.Fatalf("expected level=warn, got %v", out["level"])
	}
}

func TestApply_MultipleRules(t *testing.T) {
	tgr, _ := tag.New([]tag.Rule{
		{Field: "env", Value: "staging"},
		{Field: "region", Value: "us-east-1"},
	})
	rec := base()
	out := tgr.Apply(rec)
	if out["env"] != "staging" {
		t.Errorf("expected env=staging")
	}
	if out["region"] != "us-east-1" {
		t.Errorf("expected region=us-east-1")
	}
}

func TestApply_OriginalRecordMutated(t *testing.T) {
	// Apply mutates in place and returns the same map.
	tgr, _ := tag.New([]tag.Rule{{Field: "dc", Value: "dc1"}})
	rec := base()
	out := tgr.Apply(rec)
	if out["dc"] != rec["dc"] {
		t.Fatal("expected Apply to mutate and return the same map")
	}
}
