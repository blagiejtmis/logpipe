package checkpoint_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/logpipe/internal/checkpoint"
)

func tmpPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "checkpoint.json")
}

func TestNew_EmptyPath_ReturnsError(t *testing.T) {
	_, err := checkpoint.New("")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestNew_NonExistentFile_CreatesStore(t *testing.T) {
	s, err := checkpoint.New(tmpPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := s.Get("/var/log/app.log"); got != 0 {
		t.Fatalf("expected 0 offset, got %d", got)
	}
}

func TestSet_PersistsAcrossReopen(t *testing.T) {
	p := tmpPath(t)
	s, err := checkpoint.New(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := s.Set("/var/log/app.log", 1024); err != nil {
		t.Fatalf("Set error: %v", err)
	}

	s2, err := checkpoint.New(p)
	if err != nil {
		t.Fatalf("reopen error: %v", err)
	}
	if got := s2.Get("/var/log/app.log"); got != 1024 {
		t.Fatalf("expected 1024, got %d", got)
	}
}

func TestSet_MultipleSourcesIndependent(t *testing.T) {
	p := tmpPath(t)
	s, _ := checkpoint.New(p)
	_ = s.Set("/a", 100)
	_ = s.Set("/b", 200)

	s2, _ := checkpoint.New(p)
	if s2.Get("/a") != 100 {
		t.Errorf("/a: expected 100")
	}
	if s2.Get("/b") != 200 {
		t.Errorf("/b: expected 200")
	}
}

func TestNew_CorruptFile_ReturnsError(t *testing.T) {
	p := tmpPath(t)
	_ = os.WriteFile(p, []byte("not json{"), 0o644)
	_, err := checkpoint.New(p)
	if err == nil {
		t.Fatal("expected error for corrupt checkpoint file")
	}
}

func TestGet_UnknownSource_ReturnsZero(t *testing.T) {
	s, _ := checkpoint.New(tmpPath(t))
	if got := s.Get("unknown"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}
