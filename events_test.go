// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: BSD-3-Clause

package gpucontext

import "testing"

func TestNullEventSource(t *testing.T) {
	// NullEventSource should implement EventSource
	var es EventSource = NullEventSource{}

	// All methods should be callable without panic
	es.OnKeyPress(func(Key, Modifiers) {})
	es.OnKeyRelease(func(Key, Modifiers) {})
	es.OnTextInput(func(string) {})
	es.OnMouseMove(func(float64, float64) {})
	es.OnMousePress(func(MouseButton, float64, float64) {})
	es.OnMouseRelease(func(MouseButton, float64, float64) {})
	es.OnScroll(func(float64, float64) {})
	es.OnResize(func(int, int) {})
	es.OnFocus(func(bool) {})
}

func TestModifiers(t *testing.T) {
	tests := []struct {
		mods     Modifiers
		shift    bool
		control  bool
		alt      bool
		super    bool
	}{
		{0, false, false, false, false},
		{ModShift, true, false, false, false},
		{ModControl, false, true, false, false},
		{ModAlt, false, false, true, false},
		{ModSuper, false, false, false, true},
		{ModShift | ModControl, true, true, false, false},
		{ModShift | ModControl | ModAlt | ModSuper, true, true, true, true},
	}

	for _, tt := range tests {
		if got := tt.mods.HasShift(); got != tt.shift {
			t.Errorf("Modifiers(%d).HasShift() = %v, want %v", tt.mods, got, tt.shift)
		}
		if got := tt.mods.HasControl(); got != tt.control {
			t.Errorf("Modifiers(%d).HasControl() = %v, want %v", tt.mods, got, tt.control)
		}
		if got := tt.mods.HasAlt(); got != tt.alt {
			t.Errorf("Modifiers(%d).HasAlt() = %v, want %v", tt.mods, got, tt.alt)
		}
		if got := tt.mods.HasSuper(); got != tt.super {
			t.Errorf("Modifiers(%d).HasSuper() = %v, want %v", tt.mods, got, tt.super)
		}
	}
}

func TestKeyConstants(t *testing.T) {
	// Verify key codes are unique and sequential
	keys := []Key{
		KeyA, KeyB, KeyC, KeyZ,
		Key0, Key1, Key9,
		KeyF1, KeyF12,
		KeyEscape, KeyEnter, KeySpace,
	}

	seen := make(map[Key]bool)
	for _, k := range keys {
		if seen[k] {
			t.Errorf("Duplicate key code: %d", k)
		}
		seen[k] = true
	}
}

func TestMouseButtonConstants(t *testing.T) {
	// Verify mouse button codes
	if MouseButtonLeft != 0 {
		t.Error("MouseButtonLeft should be 0")
	}
	if MouseButtonRight != 1 {
		t.Error("MouseButtonRight should be 1")
	}
	if MouseButtonMiddle != 2 {
		t.Error("MouseButtonMiddle should be 2")
	}
}
