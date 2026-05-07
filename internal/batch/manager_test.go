package batch

import (
	"testing"
	"time"
)

func TestNewManager_NilConfig_UsesDefaults(t *testing.T) {
	m, err := NewManager(nil)
	if err != nil {
		t.Fatal(err)
	}
	if m.defaultMaxSize != 100 {
		t.Errorf("expected defaultMaxSize=100, got %d", m.defaultMaxSize)
	}
	if m.defaultInterval != 5*time.Second {
		t.Errorf("expected defaultInterval=5s, got %v", m.defaultInterval)
	}
}

func TestNewManager_CustomConfig(t *testing.T) {
	cfg := &Config{MaxSize: 50, FlushInterval: 2 * time.Second}
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if m.defaultMaxSize != 50 {
		t.Errorf("expected 50, got %d", m.defaultMaxSize)
	}
	if m.defaultInterval != 2*time.Second {
		t.Errorf("expected 2s, got %v", m.defaultInterval)
	}
}

func TestNewManager_InvalidMaxSize_ReturnsError(t *testing.T) {
	_, err := NewManager(&Config{MaxSize: 0, FlushInterval: time.Second})
	if err == nil {
		t.Fatal("expected error for MaxSize=0")
	}
}

func TestNewManager_InvalidInterval_ReturnsError(t *testing.T) {
	_, err := NewManager(&Config{MaxSize: 10, FlushInterval: -1})
	if err == nil {
		t.Fatal("expected error for negative interval")
	}
}

func TestManager_NewBatcher_Works(t *testing.T) {
	m, _ := NewManager(nil)
	var received []Record
	b, err := m.NewBatcher(func(batch []Record) {
		received = append(received, batch...)
	})
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 100; i++ {
		b.Add(Record{"i": "v"})
	}
	if len(received) != 100 {
		t.Errorf("expected 100 records flushed, got %d", len(received))
	}
}
