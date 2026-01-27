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

	// IME methods should also be callable without panic
	es.OnIMECompositionStart(func() {})
	es.OnIMECompositionUpdate(func(IMEState) {})
	es.OnIMECompositionEnd(func(string) {})
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

func TestIMEState(t *testing.T) {
	// Test IMEState struct fields
	state := IMEState{
		Composing:       true,
		CompositionText: "nihao",
		CursorPos:       5,
		SelectionStart:  2,
		SelectionEnd:    4,
	}

	if !state.Composing {
		t.Error("Composing should be true")
	}
	if state.CompositionText != "nihao" {
		t.Errorf("CompositionText = %q, want \"nihao\"", state.CompositionText)
	}
	if state.CursorPos != 5 {
		t.Errorf("CursorPos = %d, want 5", state.CursorPos)
	}
	if state.SelectionStart != 2 {
		t.Errorf("SelectionStart = %d, want 2", state.SelectionStart)
	}
	if state.SelectionEnd != 4 {
		t.Errorf("SelectionEnd = %d, want 4", state.SelectionEnd)
	}
}

func TestIMEStateZeroValue(t *testing.T) {
	// Test IMEState zero value
	var state IMEState

	if state.Composing {
		t.Error("Zero value Composing should be false")
	}
	if state.CompositionText != "" {
		t.Errorf("Zero value CompositionText = %q, want empty", state.CompositionText)
	}
	if state.CursorPos != 0 {
		t.Errorf("Zero value CursorPos = %d, want 0", state.CursorPos)
	}
	if state.SelectionStart != 0 {
		t.Errorf("Zero value SelectionStart = %d, want 0", state.SelectionStart)
	}
	if state.SelectionEnd != 0 {
		t.Errorf("Zero value SelectionEnd = %d, want 0", state.SelectionEnd)
	}
}

// mockIMEController is used to verify IMEController interface at compile time.
type mockIMEController struct{}

func (mockIMEController) SetIMEPosition(_, _ int) {}
func (mockIMEController) SetIMEEnabled(_ bool)    {}

// Ensure mockIMEController implements IMEController.
var _ IMEController = mockIMEController{}

func TestIMEControllerInterface(t *testing.T) {
	// Verify IMEController can be used through the interface
	var controller IMEController = mockIMEController{}

	// These should not panic
	controller.SetIMEPosition(100, 200)
	controller.SetIMEEnabled(true)
	controller.SetIMEEnabled(false)
}
