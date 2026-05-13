package lookup

import (
	"testing"

	"github.com/yourorg/logpipe/internal/pipeline"
)

func base() pipeline.Record {
	return pipeline.Record{"level": "info", "env": "prod", "message": "ok"}
}

func TestNew_NoRules_ReturnsError(t *testing.T) {
	_, err := New(nil)
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestNew_EmptyKeyField_ReturnsError(t *testing.T) {
	_, err := New([]Rule{{KeyField: "", Table: map[string]map[string]string{"a": {"x": "1"}}}})
	if err == nil {
		t.Fatal("expected error for empty key_field")
	}
}

func TestNew_EmptyTable_ReturnsError(t *testing.T) {
	_, err := New([]Rule{{KeyField: "env", Table: nil}})
	if err == nil {
		t.Fatal("expected error for empty table")
	}
}

func TestNew_InvalidOnMiss_ReturnsError(t *testing.T) {
	_, err := New([]Rule{{KeyField: "env", Table: map[string]map[string]string{"prod": {"dc": "us-east"}}, OnMiss: "skip"}})
	if err == nil {
		t.Fatal("expected error for invalid on_miss value")
	}
}

func TestApply_MergesMatchedFields(t *testing.T) {
	l, err := New([]Rule{{
		KeyField: "env",
		Table:    map[string]map[string]string{"prod": {"dc": "us-east", "tier": "production"}},
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rec, keep := l.Apply(base())
	if !keep {
		t.Fatal("expected record to be kept")
	}
	if rec["dc"] != "us-east" {
		t.Errorf("expected dc=us-east, got %v", rec["dc"])
	}
	if rec["tier"] != "production" {
		t.Errorf("expected tier=production, got %v", rec["tier"])
	}
}

func TestApply_DoesNotOverwriteExisting(t *testing.T) {
	l, _ := New([]Rule{{
		KeyField: "env",
		Table:    map[string]map[string]string{"prod": {"level": "override"}},
	}})
	rec, _ := l.Apply(base())
	if rec["level"] != "info" {
		t.Errorf("existing field should not be overwritten, got %v", rec["level"])
	}
}

func TestApply_MissIgnore_KeepsRecord(t *testing.T) {
	l, _ := New([]Rule{{
		KeyField: "env",
		Table:    map[string]map[string]string{"staging": {"dc": "eu-west"}},
		OnMiss:   "ignore",
	}})
	_, keep := l.Apply(base()) // env=prod, not in table
	if !keep {
		t.Fatal("expected record to be kept on miss with ignore policy")
	}
}

func TestApply_MissDrop_DropsRecord(t *testing.T) {
	l, _ := New([]Rule{{
		KeyField: "env",
		Table:    map[string]map[string]string{"staging": {"dc": "eu-west"}},
		OnMiss:   "drop",
	}})
	_, keep := l.Apply(base()) // env=prod, not in table
	if keep {
		t.Fatal("expected record to be dropped on miss with drop policy")
	}
}

func TestApply_MissingKeyField_Ignore(t *testing.T) {
	l, _ := New([]Rule{{
		KeyField: "region",
		Table:    map[string]map[string]string{"us-east": {"cloud": "aws"}},
	}})
	_, keep := l.Apply(base()) // no "region" field
	if !keep {
		t.Fatal("expected record kept when key field absent and on_miss=ignore")
	}
}
