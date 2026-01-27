package gpucontext

import (
	"testing"
)

func TestRegistry_RegisterAndGet(t *testing.T) {
	r := NewRegistry[string]()

	r.Register("foo", func() string { return "bar" })

	if got := r.Get("foo"); got != "bar" {
		t.Errorf("Get(foo) = %q, want %q", got, "bar")
	}
}

func TestRegistry_Get_NotFound(t *testing.T) {
	r := NewRegistry[string]()

	if got := r.Get("nonexistent"); got != "" {
		t.Errorf("Get(nonexistent) = %q, want empty string", got)
	}
}

func TestRegistry_Has(t *testing.T) {
	r := NewRegistry[int]()

	if r.Has("foo") {
		t.Error("Has(foo) = true before registration")
	}

	r.Register("foo", func() int { return 42 })

	if !r.Has("foo") {
		t.Error("Has(foo) = false after registration")
	}
}

func TestRegistry_Unregister(t *testing.T) {
	r := NewRegistry[int]()
	r.Register("foo", func() int { return 42 })

	r.Unregister("foo")

	if r.Has("foo") {
		t.Error("Has(foo) = true after unregister")
	}
}

func TestRegistry_Best_WithPriority(t *testing.T) {
	r := NewRegistry[string](WithPriority("high", "medium", "low"))

	r.Register("low", func() string { return "low-value" })
	r.Register("medium", func() string { return "medium-value" })

	// Should return medium (highest priority among registered)
	if got := r.Best(); got != "medium-value" {
		t.Errorf("Best() = %q, want %q", got, "medium-value")
	}

	// Register high priority
	r.Register("high", func() string { return "high-value" })

	if got := r.Best(); got != "high-value" {
		t.Errorf("Best() = %q, want %q", got, "high-value")
	}
}

func TestRegistry_Best_NoPriority(t *testing.T) {
	r := NewRegistry[string]()

	r.Register("a", func() string { return "a-value" })

	// Should return some value (order not guaranteed)
	if got := r.Best(); got == "" {
		t.Error("Best() returned empty string")
	}
}

func TestRegistry_Best_Empty(t *testing.T) {
	r := NewRegistry[string]()

	if got := r.Best(); got != "" {
		t.Errorf("Best() on empty registry = %q, want empty", got)
	}
}

func TestRegistry_BestName(t *testing.T) {
	r := NewRegistry[int](WithPriority("first", "second"))

	r.Register("second", func() int { return 2 })

	if got := r.BestName(); got != "second" {
		t.Errorf("BestName() = %q, want %q", got, "second")
	}

	r.Register("first", func() int { return 1 })

	if got := r.BestName(); got != "first" {
		t.Errorf("BestName() = %q, want %q", got, "first")
	}
}

func TestRegistry_Available(t *testing.T) {
	r := NewRegistry[int]()

	r.Register("a", func() int { return 1 })
	r.Register("b", func() int { return 2 })

	available := r.Available()
	if len(available) != 2 {
		t.Errorf("Available() length = %d, want 2", len(available))
	}

	// Check both are present (order not guaranteed)
	hasA, hasB := false, false
	for _, name := range available {
		if name == "a" {
			hasA = true
		}
		if name == "b" {
			hasB = true
		}
	}
	if !hasA || !hasB {
		t.Errorf("Available() = %v, want [a, b]", available)
	}
}

func TestRegistry_Count(t *testing.T) {
	r := NewRegistry[int]()

	if r.Count() != 0 {
		t.Errorf("Count() on empty = %d, want 0", r.Count())
	}

	r.Register("a", func() int { return 1 })
	r.Register("b", func() int { return 2 })

	if r.Count() != 2 {
		t.Errorf("Count() = %d, want 2", r.Count())
	}
}

func TestRegistry_Replace(t *testing.T) {
	r := NewRegistry[string]()

	r.Register("foo", func() string { return "first" })
	r.Register("foo", func() string { return "second" })

	if got := r.Get("foo"); got != "second" {
		t.Errorf("Get(foo) after replace = %q, want %q", got, "second")
	}
}
