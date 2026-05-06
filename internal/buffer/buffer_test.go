package buffer

import (
	"sync"
	"testing"
)

func rec(src string) Record { return Record{Source: src, Fields: map[string]string{"msg": src}} }

func TestNew_InvalidCapacity(t *testing.T) {
	_, err := New(0, false)
	if err == nil {
		t.Fatal("expected error for zero capacity")
	}
}

func TestPushPop_FIFO(t *testing.T) {
	b, _ := New(4, false)
	b.Push(rec("a"))
	b.Push(rec("b"))
	b.Push(rec("c"))

	r, ok := b.Pop()
	if !ok || r.Source != "a" {
		t.Fatalf("expected 'a', got %q ok=%v", r.Source, ok)
	}
	r, _ = b.Pop()
	if r.Source != "b" {
		t.Fatalf("expected 'b', got %q", r.Source)
	}
}

func TestPop_EmptyBuffer(t *testing.T) {
	b, _ := New(2, false)
	_, ok := b.Pop()
	if ok {
		t.Fatal("expected false pop on empty buffer")
	}
}

func TestPush_DropNew_ReturnsFull(t *testing.T) {
	b, _ := New(2, true)
	b.Push(rec("a"))
	b.Push(rec("b"))

	err := b.Push(rec("c"))
	if err != ErrFull {
		t.Fatalf("expected ErrFull, got %v", err)
	}
	if b.Dropped() != 1 {
		t.Fatalf("expected 1 dropped, got %d", b.Dropped())
	}
	if b.Len() != 2 {
		t.Fatalf("len should remain 2, got %d", b.Len())
	}
}

func TestPush_EvictOldest_WhenFull(t *testing.T) {
	b, _ := New(2, false)
	b.Push(rec("a"))
	b.Push(rec("b"))
	b.Push(rec("c")) // evicts "a"

	r, _ := b.Pop()
	if r.Source != "b" {
		t.Fatalf("expected 'b' after eviction, got %q", r.Source)
	}
	if b.Dropped() != 1 {
		t.Fatalf("expected 1 dropped, got %d", b.Dropped())
	}
}

func TestBuffer_ConcurrentAccess(t *testing.T) {
	b, _ := New(64, false)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			b.Push(rec("x"))
			b.Pop()
		}()
	}
	wg.Wait()
}
