package aggregate

import (
	"testing"
	"time"
)

func TestNew_InvalidField_ReturnsError(t *testing.T) {
	_, err := New([]Rule{{Field: "", Op: OpCount, Window: time.Second}})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_InvalidWindow_ReturnsError(t *testing.T) {
	_, err := New([]Rule{{Field: "bytes", Op: OpSum, Window: 0}})
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNew_UnknownOp_ReturnsError(t *testing.T) {
	_, err := New([]Rule{{Field: "bytes", Op: Op("avg"), Window: time.Second}})
	if err == nil {
		t.Fatal("expected error for unknown op")
	}
}

func TestNew_ValidRules_NoError(t *testing.T) {
	_, err := New([]Rule{
		{Field: "bytes", Op: OpSum, Window: time.Second},
		{Field: "errors", Op: OpCount, Window: time.Minute},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAdd_Count(t *testing.T) {
	a, _ := New([]Rule{{Field: "msg", Op: OpCount, Window: time.Minute}})
	for i := 0; i < 5; i++ {
		a.Add("src", map[string]any{"msg": "hello"})
	}
	snap := a.Snapshot()
	if snap["src"]["msg"] != 5 {
		t.Fatalf("expected count 5, got %v", snap["src"]["msg"])
	}
}

func TestAdd_Sum(t *testing.T) {
	a, _ := New([]Rule{{Field: "bytes", Op: OpSum, Window: time.Minute}})
	a.Add("src", map[string]any{"bytes": float64(100)})
	a.Add("src", map[string]any{"bytes": float64(200)})
	snap := a.Snapshot()
	if snap["src"]["bytes"] != 300 {
		t.Fatalf("expected sum 300, got %v", snap["src"]["bytes"])
	}
}

func TestAdd_Min(t *testing.T) {
	a, _ := New([]Rule{{Field: "lat", Op: OpMin, Window: time.Minute}})
	a.Add("src", map[string]any{"lat": float64(50)})
	a.Add("src", map[string]any{"lat": float64(10)})
	a.Add("src", map[string]any{"lat": float64(30)})
	snap := a.Snapshot()
	if snap["src"]["lat"] != 10 {
		t.Fatalf("expected min 10, got %v", snap["src"]["lat"])
	}
}

func TestAdd_Max(t *testing.T) {
	a, _ := New([]Rule{{Field: "lat", Op: OpMax, Window: time.Minute}})
	a.Add("src", map[string]any{"lat": float64(50)})
	a.Add("src", map[string]any{"lat": float64(10)})
	a.Add("src", map[string]any{"lat": float64(99)})
	snap := a.Snapshot()
	if snap["src"]["lat"] != 99 {
		t.Fatalf("expected max 99, got %v", snap["src"]["lat"])
	}
}

func TestAdd_IndependentSources(t *testing.T) {
	a, _ := New([]Rule{{Field: "n", Op: OpCount, Window: time.Minute}})
	a.Add("a", map[string]any{"n": 1})
	a.Add("a", map[string]any{"n": 1})
	a.Add("b", map[string]any{"n": 1})
	snap := a.Snapshot()
	if snap["a"]["n"] != 2 || snap["b"]["n"] != 1 {
		t.Fatalf("sources not independent: %v", snap)
	}
}
