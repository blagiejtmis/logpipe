package fingerprint_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/fingerprint"
)

func base() fingerprint.Record {
	return fingerprint.Record{
		"level":   "info",
		"message": "hello world",
		"host":    "web-01",
	}
}

func TestNew_BlankField_ReturnsError(t *testing.T) {
	_, err := fingerprint.New(fingerprint.Config{Fields: []string{"level", "  "}})
	if err == nil {
		t.Fatal("expected error for blank field name")
	}
}

func TestNew_ValidConfig_NoError(t *testing.T) {
	_, err := fingerprint.New(fingerprint.Config{Fields: []string{"level", "message"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCompute_DeterministicForSameRecord(t *testing.T) {
	fp, _ := fingerprint.New(fingerprint.Config{Fields: []string{"level", "message"}})
	rec := base()
	a := fp.Compute(rec)
	b := fp.Compute(rec)
	if a != b {
		t.Errorf("expected same fingerprint, got %s vs %s", a, b)
	}
}

func TestCompute_DifferentValues_DifferentHash(t *testing.T) {
	fp, _ := fingerprint.New(fingerprint.Config{Fields: []string{"message"}})
	r1 := fingerprint.Record{"message": "hello"}
	r2 := fingerprint.Record{"message": "world"}
	if fp.Compute(r1) == fp.Compute(r2) {
		t.Error("expected different fingerprints for different values")
	}
}

func TestCompute_NoFields_UsesAllKeysSorted(t *testing.T) {
	fp, _ := fingerprint.New(fingerprint.Config{})
	rec := base()
	h1 := fp.Compute(rec)
	if len(h1) != 64 {
		t.Errorf("expected 64-char hex string, got len %d", len(h1))
	}
	// Same record should always produce the same hash.
	if fp.Compute(rec) != h1 {
		t.Error("non-deterministic fingerprint with all-fields mode")
	}
}

func TestCompute_MissingField_TreatedAsNil(t *testing.T) {
	fp, _ := fingerprint.New(fingerprint.Config{Fields: []string{"missing"}})
	rec := base()
	// Should not panic; missing key yields "<nil>" via fmt.Sprintf.
	h := fp.Compute(rec)
	if h == "" {
		t.Error("expected non-empty fingerprint even for missing field")
	}
}

func TestCompute_CustomSeparator(t *testing.T) {
	fp1, _ := fingerprint.New(fingerprint.Config{Fields: []string{"level", "host"}, Separator: "|"})
	fp2, _ := fingerprint.New(fingerprint.Config{Fields: []string{"level", "host"}, Separator: ":"})
	rec := base()
	if fp1.Compute(rec) == fp2.Compute(rec) {
		t.Error("expected different fingerprints for different separators")
	}
}
