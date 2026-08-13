// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

package gpucontext

import "testing"

func TestMouseButtonString(t *testing.T) {
	tests := []struct {
		name   string
		button MouseButton
		want   string
	}{
		{name: "left", button: MouseButtonLeft, want: "Left"},
		{name: "right", button: MouseButtonRight, want: "Right"},
		{name: "middle", button: MouseButtonMiddle, want: "Middle"},
		{name: "button four", button: MouseButton4, want: "X1"},
		{name: "button five", button: MouseButton5, want: "X2"},
		{name: "invalid", button: MouseButton(255), want: "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.button.String(); got != tt.want {
				t.Errorf("MouseButton(%d).String() = %q, want %q", tt.button, got, tt.want)
			}
		})
	}
}

func TestModifiersString(t *testing.T) {
	const unknown Modifiers = 1 << 6

	tests := []struct {
		name string
		mods Modifiers
		want string
	}{
		{name: "zero", mods: 0, want: "None"},
		{name: "control", mods: ModControl, want: "Ctrl"},
		{name: "alt", mods: ModAlt, want: "Alt"},
		{name: "shift", mods: ModShift, want: "Shift"},
		{name: "super", mods: ModSuper, want: "Super"},
		{name: "caps lock", mods: ModCapsLock, want: "CapsLock"},
		{name: "num lock", mods: ModNumLock, want: "NumLock"},
		{name: "control and shift", mods: ModControl | ModShift, want: "Ctrl+Shift"},
		{name: "control alt shift", mods: ModControl | ModAlt | ModShift, want: "Ctrl+Alt+Shift"},
		{name: "all known", mods: ModControl | ModAlt | ModShift | ModSuper | ModCapsLock | ModNumLock, want: "Ctrl+Alt+Shift+Super+CapsLock+NumLock"},
		// Unknown bits are ignored to stay compatible with gogpu/ui's
		// event.Modifiers.String implementation.
		{name: "unknown only", mods: unknown, want: ""},
		{name: "unknown mixed", mods: ModControl | unknown, want: "Ctrl"},
		{name: "unknown and all known", mods: ModControl | ModAlt | ModShift | ModSuper | ModCapsLock | ModNumLock | unknown, want: "Ctrl+Alt+Shift+Super+CapsLock+NumLock"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mods.String(); got != tt.want {
				t.Errorf("Modifiers(%d).String() = %q, want %q", tt.mods, got, tt.want)
			}
		})
	}
}
