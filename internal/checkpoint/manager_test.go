package checkpoint

import (
	"path/filepath"
	"testing"
)

func TestNewManager_EmptyPath_ReturnsError(t *testing.T) {
	_, err := NewManager("")
	if err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

func TestNewManager_CreatesStore(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ckpt.db")
	m, err := NewManager(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer m.Close()
}

func TestManager_GetSet_RoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ckpt.db")
	m, err := NewManager(path)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	defer m.Close()

	_, ok := m.Get("source-a")
	if ok {
		t.Fatal("expected no offset before Set")
	}

	if err := m.Set("source-a", 42); err != nil {
		t.Fatalf("Set: %v", err)
	}

	offset, ok := m.Get("source-a")
	if !ok {
		t.Fatal("expected offset after Set")
	}
	if offset != 42 {
		t.Fatalf("expected 42, got %d", offset)
	}
}

func TestManager_MultipleSourcesIndependent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ckpt.db")
	m, err := NewManager(path)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	defer m.Close()

	sources := map[string]int64{"a": 10, "b": 20, "c": 30}
	for src, off := range sources {
		if err := m.Set(src, off); err != nil {
			t.Fatalf("Set %q: %v", src, err)
		}
	}
	for src, want := range sources {
		got, ok := m.Get(src)
		if !ok {
			t.Fatalf("missing offset for %q", src)
		}
		if got != want {
			t.Fatalf("source %q: expected %d, got %d", src, want, got)
		}
	}
}

func TestManager_PersistsAcrossReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ckpt.db")

	m1, err := NewManager(path)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	if err := m1.Set("src", 99); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if err := m1.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	m2, err := NewManager(path)
	if err != nil {
		t.Fatalf("reopen NewManager: %v", err)
	}
	defer m2.Close()

	offset, ok := m2.Get("src")
	if !ok {
		t.Fatal("expected persisted offset after reopen")
	}
	if offset != 99 {
		t.Fatalf("expected 99, got %d", offset)
	}
}
