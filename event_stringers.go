// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

package gpucontext

import "strings"

// String returns the mouse button name for debugging.
//
// The names for the extra buttons match [Button.String]: X1 and X2. Values
// outside the defined buttons are reported as Unknown.
func (b MouseButton) String() string {
	switch b {
	case MouseButtonLeft:
		return stringLeft
	case MouseButtonRight:
		return stringRight
	case MouseButtonMiddle:
		return stringMiddle
	case MouseButton4:
		return "X1"
	case MouseButton5:
		return "X2"
	default:
		return "Unknown"
	}
}

// String returns a human-readable representation of the modifier keys.
//
// Combined modifiers are ordered as Ctrl, Alt, Shift, Super, CapsLock, and
// NumLock, matching the gogpu/ui event convention. Unknown bits are ignored,
// also matching the gogpu/ui event convention; a value containing only unknown
// bits therefore produces an empty string.
func (m Modifiers) String() string {
	if m == 0 {
		return stringNone
	}

	parts := make([]string, 0, 6)
	if m&ModControl != 0 {
		parts = append(parts, "Ctrl")
	}
	if m&ModAlt != 0 {
		parts = append(parts, "Alt")
	}
	if m&ModShift != 0 {
		parts = append(parts, "Shift")
	}
	if m&ModSuper != 0 {
		parts = append(parts, "Super")
	}
	if m&ModCapsLock != 0 {
		parts = append(parts, "CapsLock")
	}
	if m&ModNumLock != 0 {
		parts = append(parts, "NumLock")
	}
	return strings.Join(parts, "+")
}
