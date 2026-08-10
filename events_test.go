// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

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
		mods    Modifiers
		shift   bool
		control bool
		alt     bool
		super   bool
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

func TestKeyGroupedBases(t *testing.T) {
	tests := []struct {
		name string
		key  Key
		want Key
	}{
		{"KeyUnknown", KeyUnknown, 0},
		{"KeyA (letters base)", KeyA, 1},
		{"KeyZ (letters end)", KeyZ, 26},
		{"Key0 (digits base)", Key0, 33},
		{"Key9 (digits end)", Key9, 42},
		{"KeyF1 (function base)", KeyF1, 49},
		{"KeyF12", KeyF12, 60},
		{"KeyF13", KeyF13, 61},
		{"KeyF24 (function end)", KeyF24, 72},
		{"KeyEscape (navigation base)", KeyEscape, 81},
		{"KeyDown (navigation end)", KeyDown, 95},
		{"KeyLeftShift (modifiers base)", KeyLeftShift, 113},
		{"KeyRightSuper (modifiers end)", KeyRightSuper, 120},
		{"KeyMinus (punctuation base)", KeyMinus, 129},
		{"KeyIntlYen (punctuation end)", KeyIntlYen, 141},
		{"KeyNumpad0 (numpad base)", KeyNumpad0, 161},
		{"KeyNumpadComma (numpad end)", KeyNumpadComma, 178},
		{"KeyCapsLock (locks base)", KeyCapsLock, 193},
		{"KeyPause (locks end)", KeyPause, 197},
		{"KeyMediaPlayPause (media base)", KeyMediaPlayPause, 209},
		{"KeyMediaRecord (media end)", KeyMediaRecord, 213},
		{"KeyAudioVolumeUp (volume base)", KeyAudioVolumeUp, 241},
		{"KeyAudioVolumeMute (volume end)", KeyAudioVolumeMute, 243},
		{"KeyBrowserBack (browser base)", KeyBrowserBack, 249},
		{"KeyBrowserSearch (browser end)", KeyBrowserSearch, 253},
		{"KeyContextMenu (system base)", KeyContextMenu, 265},
		{"KeyLaunchApp2 (system end)", KeyLaunchApp2, 268},
	}
	for _, tt := range tests {
		if tt.key != tt.want {
			t.Errorf("%s = %d, want %d", tt.name, tt.key, tt.want)
		}
	}
}

func TestKeyGroupsNoOverlap(t *testing.T) {
	groups := []struct {
		name     string
		first    Key
		last     Key
		maxRange Key
	}{
		{"Letters", KeyA, KeyZ, 31},
		{"Digits", Key0, Key9, 47},
		{"Function", KeyF1, KeyF24, 80},
		{"Navigation", KeyEscape, KeyDown, 112},
		{"Modifiers", KeyLeftShift, KeyRightSuper, 128},
		{"Punctuation", KeyMinus, KeyIntlYen, 160},
		{"Numpad", KeyNumpad0, KeyNumpadComma, 192},
		{"Locks", KeyCapsLock, KeyPause, 208},
		{"Media", KeyMediaPlayPause, KeyMediaRecord, 240},
		{"Volume", KeyAudioVolumeUp, KeyAudioVolumeMute, 248},
		{"Browser", KeyBrowserBack, KeyBrowserSearch, 264},
		{"System", KeyContextMenu, KeyLaunchApp2, 280},
	}
	for i, g := range groups {
		if g.last > g.maxRange {
			t.Errorf("%s: last key %d exceeds group max range %d", g.name, g.last, g.maxRange)
		}
		if i > 0 {
			prev := groups[i-1]
			if g.first <= prev.maxRange {
				t.Errorf("%s (base %d) overlaps with %s (max %d)", g.name, g.first, prev.name, prev.maxRange)
			}
		}
	}
}

func TestKeyStringRoundTrip(t *testing.T) {
	keys := []Key{
		KeyA, KeyZ, Key0, Key9, KeyF1, KeyF12, KeyF13, KeyF24,
		KeyEscape, KeyEnter, KeySpace, KeyLeft, KeyDown,
		KeyLeftShift, KeyRightSuper,
		KeyMinus, KeySlash, KeyIntlBackslash, KeyIntlYen,
		KeyNumpad0, KeyNumpadEnter, KeyNumpadEqual, KeyNumpadComma,
		KeyCapsLock, KeyPause,
		KeyMediaPlayPause, KeyMediaStop, KeyMediaTrackNext, KeyMediaTrackPrevious, KeyMediaRecord,
		KeyAudioVolumeUp, KeyAudioVolumeDown, KeyAudioVolumeMute,
		KeyBrowserBack, KeyBrowserSearch,
		KeyContextMenu, KeyCancel, KeyLaunchApp1, KeyLaunchApp2,
	}
	for _, k := range keys {
		name := k.String()
		got, ok := KeyFromString(name)
		if !ok {
			t.Errorf("KeyFromString(%q) returned false for valid key %d", name, k)
			continue
		}
		if got != k {
			t.Errorf("KeyFromString(%q) = %d, want %d", name, got, k)
		}
	}
}

func TestKeyFromStringUnknown(t *testing.T) {
	tests := []string{"", "nonexistent", "KEY_A", "keyA", "ArrowLeft"}
	for _, name := range tests {
		k, ok := KeyFromString(name)
		if ok {
			t.Errorf("KeyFromString(%q) = (%d, true), want (_, false)", name, k)
		}
		if k != KeyUnknown {
			t.Errorf("KeyFromString(%q) = %d, want KeyUnknown (0)", name, k)
		}
	}
}

func TestKeyStringNotEmpty(t *testing.T) {
	keys := []Key{
		KeyA, KeyZ, Key0, Key9, KeyF1, KeyF24, KeyEscape,
		KeyContextMenu, KeyCancel, KeyMediaPlayPause, KeyAudioVolumeMute,
	}
	for _, k := range keys {
		if s := k.String(); s == "" {
			t.Errorf("Key(%d).String() is empty", k)
		}
	}
}

func TestKeyConstants(t *testing.T) {
	keys := []Key{
		KeyA, KeyB, KeyC, KeyZ,
		Key0, Key1, Key9,
		KeyF1, KeyF12, KeyF13, KeyF24,
		KeyEscape, KeyEnter, KeySpace,
		KeyContextMenu, KeyCancel,
		KeyMediaPlayPause, KeyAudioVolumeUp,
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
