package batch

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestNew_InvalidArgs(t *testing.T) {
	noop := func([]Record) {}
	if _, err := New(0, time.Second, noop); err != ErrInvalidMaxSize {
		t.Fatalf("expected ErrInvalidMaxSize, got %v", err)
	}
	if _, err := New(1, 0, noop); err != ErrInvalidInterval {
		t.Fatalf("expected ErrInvalidInterval, got %v", err)
	}
	if _, err := New(1, time.Second, nil); err != ErrNilFlushFunc {
		t.Fatalf("expected ErrNilFlushFunc, got %v", err)
	}
}

func TestAdd_FlushesOnMaxSize(t *testing.T) {
	var mu sync.Mutex
	var got [][]Record
	b, err := New(3, time.Hour, func(batch []Record) {
		mu.Lock()
		got = append(got, batch)
		mu.Unlock()
	})
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		b.Add(Record{"n": "x"})
	}
	mu.Lock()
	defer mu.Unlock()
	if len(got) != 1 || len(got[0]) != 3 {
		t.Fatalf("expected 1 flush of 3 records, got %v", got)
	}
}

func TestRun_FlushesOnInterval(t *testing.T) {
	var mu sync.Mutex
	var got [][]Record
	b, err := New(100, 30*time.Millisecond, func(batch []Record) {
		mu.Lock()
		got = append(got, batch)
		mu.Unlock()
	})
	if err != nil {
		t.Fatal(err)
	}
	b.Add(Record{"k": "v"})
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()
	b.Run(ctx)
	mu.Lock()
	defer mu.Unlock()
	if len(got) == 0 {
		t.Fatal("expected at least one interval flush")
	}
}

func TestRun_FlushesOnCancel(t *testing.T) {
	flushed := make(chan []Record, 1)
	b, _ := New(100, time.Hour, func(batch []Record) {
		flushed <- batch
	})
	b.Add(Record{"a": "1"})
	b.Add(Record{"b": "2"})
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()
	b.Run(ctx)
	select {
	case batch := <-flushed:
		if len(batch) != 2 {
			t.Fatalf("expected 2 records on cancel flush, got %d", len(batch))
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for cancel flush")
	}
}

func TestAdd_NoFlush_WhenUnderMaxSize(t *testing.T) {
	flushCount := 0
	b, _ := New(5, time.Hour, func([]Record) { flushCount++ })
	b.Add(Record{"x": "1"})
	b.Add(Record{"x": "2"})
	if flushCount != 0 {
		t.Fatalf("expected no flush, got %d", flushCount)
	}
}
