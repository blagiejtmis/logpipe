package metrics_test

import (
	"sync"
	"testing"

	"github.com/logpipe/logpipe/internal/metrics"
)

func TestCounter_IncAndValue(t *testing.T) {
	c := &metrics.Counter{}
	if c.Value() != 0 {
		t.Fatalf("expected 0, got %d", c.Value())
	}
	c.Inc()
	c.Inc()
	if c.Value() != 2 {
		t.Fatalf("expected 2, got %d", c.Value())
	}
}

func TestCounter_Add(t *testing.T) {
	c := &metrics.Counter{}
	c.Add(10)
	if c.Value() != 10 {
		t.Fatalf("expected 10, got %d", c.Value())
	}
}

func TestRegistry_Counter_SameName(t *testing.T) {
	r := metrics.NewRegistry()
	a := r.Counter("lines_in")
	b := r.Counter("lines_in")
	if a != b {
		t.Fatal("expected same counter instance for same name")
	}
}

func TestRegistry_Snapshot(t *testing.T) {
	r := metrics.NewRegistry()
	r.Counter("lines_in").Add(5)
	r.Counter("lines_out").Add(3)

	snap := r.Snapshot()
	if snap["lines_in"] != 5 {
		t.Fatalf("expected lines_in=5, got %d", snap["lines_in"])
	}
	if snap["lines_out"] != 3 {
		t.Fatalf("expected lines_out=3, got %d", snap["lines_out"])
	}
}

func TestRegistry_ConcurrentInc(t *testing.T) {
	r := metrics.NewRegistry()
	c := r.Counter("concurrent")

	var wg sync.WaitGroup
	const goroutines = 100
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			c.Inc()
		}()
	}
	wg.Wait()

	if c.Value() != goroutines {
		t.Fatalf("expected %d, got %d", goroutines, c.Value())
	}
}
