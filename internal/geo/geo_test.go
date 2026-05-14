package geo

import (
	"testing"
)

func stubLookup(db map[string]map[string]string) Lookup {
	return func(ip string) (map[string]string, bool) {
		v, ok := db[ip]
		return v, ok
	}
}

var testDB = map[string]map[string]string{
	"1.2.3.4": {"country": "US", "city": "New York"},
}

func base() map[string]any {
	return map[string]any{"ip": "1.2.3.4", "msg": "hello"}
}

func TestNew_NilLookup_ReturnsError(t *testing.T) {
	_, err := New([]Rule{{SrcField: "ip", DstField: "geo"}}, nil)
	if err == nil {
		t.Fatal("expected error for nil lookup")
	}
}

func TestNew_NoRules_ReturnsError(t *testing.T) {
	_, err := New(nil, stubLookup(testDB))
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestNew_EmptySrcField_ReturnsError(t *testing.T) {
	_, err := New([]Rule{{SrcField: "", DstField: "geo"}}, stubLookup(testDB))
	if err == nil {
		t.Fatal("expected error for blank src_field")
	}
}

func TestNew_InvalidOnMiss_ReturnsError(t *testing.T) {
	_, err := New([]Rule{{SrcField: "ip", DstField: "geo", OnMiss: "panic"}}, stubLookup(testDB))
	if err == nil {
		t.Fatal("expected error for unknown on_miss value")
	}
}

func TestApply_KnownIP_EnrichesRecord(t *testing.T) {
	e, err := New([]Rule{{SrcField: "ip", DstField: "geo"}}, stubLookup(testDB))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rec := base()
	e.Apply(rec)
	geo, ok := rec["geo"]
	if !ok {
		t.Fatal("expected geo field to be set")
	}
	attrs := geo.(map[string]string)
	if attrs["country"] != "US" {
		t.Errorf("expected country=US, got %s", attrs["country"])
	}
}

func TestApply_UnknownIP_SkipByDefault(t *testing.T) {
	e, _ := New([]Rule{{SrcField: "ip", DstField: "geo"}}, stubLookup(testDB))
	rec := map[string]any{"ip": "9.9.9.9"}
	e.Apply(rec)
	if _, ok := rec["geo"]; ok {
		t.Fatal("expected geo field to be absent on miss")
	}
}

func TestApply_UnknownIP_EmptyOnMiss(t *testing.T) {
	e, _ := New([]Rule{{SrcField: "ip", DstField: "geo", OnMiss: "empty"}}, stubLookup(testDB))
	rec := map[string]any{"ip": "9.9.9.9"}
	e.Apply(rec)
	geo, ok := rec["geo"]
	if !ok {
		t.Fatal("expected geo field with empty map")
	}
	if len(geo.(map[string]string)) != 0 {
		t.Error("expected empty map")
	}
}

func TestApply_InvalidIP_SkipsRule(t *testing.T) {
	e, _ := New([]Rule{{SrcField: "ip", DstField: "geo"}}, stubLookup(testDB))
	rec := map[string]any{"ip": "not-an-ip"}
	e.Apply(rec)
	if _, ok := rec["geo"]; ok {
		t.Fatal("expected no geo field for invalid IP")
	}
}

func TestApply_MissingSrcField_SkipsRule(t *testing.T) {
	e, _ := New([]Rule{{SrcField: "ip", DstField: "geo"}}, stubLookup(testDB))
	rec := map[string]any{"msg": "no ip here"}
	e.Apply(rec)
	if _, ok := rec["geo"]; ok {
		t.Fatal("expected no geo field when src field absent")
	}
}
