package offset

import (
	"os"
	"path/filepath"
	"testing"
)

func tmpPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "offsets.json")
}

func TestNew_EmptyPath_ReturnsError(t *testing.T) {
	_, err := New("")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestNew_NonExistentFile_CreatesStore(t *testing.T) {
	p := tmpPath(t)
	s, err := New(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := s.Get("src"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestSet_PersistsAcrossReopen(t *testing.T) {
	p := tmpPath(t)
	s, _ := New(p)
	if err := s.Set("src", 42); err != nil {
		t.Fatalf("Set error: %v", err)
	}
	s2, err := New(p)
	if err != nil {
		t.Fatalf("reopen error: %v", err)
	}
	if got := s2.Get("src"); got != 42 {
		t.Fatalf("expected 42, got %d", got)
	}
}

func TestSet_MultipleSourcesIndependent(t *testing.T) {
	p := tmpPath(t)
	s, _ := New(p)
	_ = s.Set("a", 10)
	_ = s.Set("b", 20)
	if s.Get("a") != 10 || s.Get("b") != 20 {
		t.Fatal("sources should be independent")
	}
}

func TestDelete_RemovesOffset(t *testing.T) {
	p := tmpPath(t)
	s, _ := New(p)
	_ = s.Set("src", 99)
	if err := s.Delete("src"); err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if got := s.Get("src"); got != 0 {
		t.Fatalf("expected 0 after delete, got %d", got)
	}
}

func TestNew_CorruptFile_ReturnsError(t *testing.T) {
	p := tmpPath(t)
	if err := os.WriteFile(p, []byte("not-json"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := New(p)
	if err == nil {
		t.Fatal("expected error for corrupt file")
	}
}
