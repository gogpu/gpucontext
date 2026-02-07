// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

package gpucontext

// PlatformProvider provides OS integration features.
//
// This interface enables UI frameworks (like gogpu/ui) to access platform
// capabilities such as clipboard, cursor management, and system accessibility
// preferences.
//
// Implementations:
//   - gogpu.App implements PlatformProvider via platform-specific code
//   - NullPlatformProvider provides no-op defaults for testing
//
// PlatformProvider is optional. Not all WindowProviders support platform
// integration (e.g., headless or embedded systems).
// Use type assertion to check availability:
//
//	if pp, ok := provider.(gpucontext.PlatformProvider); ok {
//	    pp.SetCursor(gpucontext.CursorPointer)
//	}
//
// Note: This interface is designed for gogpu <-> ui integration.
// The rendering library (gg) does NOT use this interface.
type PlatformProvider interface {
	// ClipboardRead reads text content from the system clipboard.
	// Returns empty string and nil error if clipboard is empty or not text.
	ClipboardRead() (string, error)

	// ClipboardWrite writes text content to the system clipboard.
	ClipboardWrite(text string) error

	// SetCursor changes the mouse cursor shape.
	// The cursor is typically reset to CursorDefault at the start of each frame.
	SetCursor(cursor CursorShape)

	// DarkMode returns true if the system dark mode is active.
	// Used for automatic theme switching.
	DarkMode() bool

	// ReduceMotion returns true if the user prefers reduced animation.
	// Used to disable or simplify animations for accessibility.
	ReduceMotion() bool

	// HighContrast returns true if the user needs high contrast mode.
	// Used to adjust colors and borders for accessibility.
	HighContrast() bool

	// FontScale returns the user's font size preference multiplier.
	// 1.0 = default system font size. Used to scale Sp (scale-independent pixels).
	FontScale() float32
}

// CursorShape represents the mouse cursor shape.
//
// These values cover the most common cursor shapes across platforms
// (Windows, macOS, Linux). They map directly to platform-specific
// cursor constants.
//
// For applications that need cursor changes:
//
//	if pp, ok := provider.(gpucontext.PlatformProvider); ok {
//	    pp.SetCursor(gpucontext.CursorText) // I-beam for text input
//	}
type CursorShape int

const (
	// CursorDefault is the standard arrow cursor.
	CursorDefault CursorShape = iota

	// CursorPointer is the hand cursor for clickable elements.
	CursorPointer

	// CursorText is the I-beam cursor for text input areas.
	CursorText

	// CursorCrosshair is the crosshair cursor for precise selection.
	CursorCrosshair

	// CursorMove is the four-arrow cursor for movable elements.
	CursorMove

	// CursorResizeNS is the north-south resize cursor.
	CursorResizeNS

	// CursorResizeEW is the east-west resize cursor.
	CursorResizeEW

	// CursorResizeNWSE is the NW-SE diagonal resize cursor.
	CursorResizeNWSE

	// CursorResizeNESW is the NE-SW diagonal resize cursor.
	CursorResizeNESW

	// CursorNotAllowed is the circle-with-line cursor for forbidden actions.
	CursorNotAllowed

	// CursorWait is the busy/wait cursor.
	CursorWait

	// CursorNone hides the cursor.
	CursorNone
)

// String returns the cursor shape name for debugging.
func (c CursorShape) String() string {
	switch c {
	case CursorDefault:
		return "Default"
	case CursorPointer:
		return "Pointer"
	case CursorText:
		return "Text"
	case CursorCrosshair:
		return "Crosshair"
	case CursorMove:
		return "Move"
	case CursorResizeNS:
		return "ResizeNS"
	case CursorResizeEW:
		return "ResizeEW"
	case CursorResizeNWSE:
		return "ResizeNWSE"
	case CursorResizeNESW:
		return "ResizeNESW"
	case CursorNotAllowed:
		return "NotAllowed"
	case CursorWait:
		return "Wait"
	case CursorNone:
		return "None"
	default:
		return "Unknown"
	}
}

// NullPlatformProvider implements PlatformProvider with no-op behavior.
// Used for testing and platforms without OS integration.
//
// Default return values:
//   - ClipboardRead: "", nil
//   - ClipboardWrite: nil
//   - SetCursor: no-op
//   - DarkMode: false
//   - ReduceMotion: false
//   - HighContrast: false
//   - FontScale: 1.0
type NullPlatformProvider struct{}

// ClipboardRead returns empty string and nil error.
func (NullPlatformProvider) ClipboardRead() (string, error) { return "", nil }

// ClipboardWrite does nothing and returns nil.
func (NullPlatformProvider) ClipboardWrite(string) error { return nil }

// SetCursor does nothing.
func (NullPlatformProvider) SetCursor(CursorShape) {}

// DarkMode returns false.
func (NullPlatformProvider) DarkMode() bool { return false }

// ReduceMotion returns false.
func (NullPlatformProvider) ReduceMotion() bool { return false }

// HighContrast returns false.
func (NullPlatformProvider) HighContrast() bool { return false }

// FontScale returns 1.0.
func (NullPlatformProvider) FontScale() float32 { return 1.0 }

// Ensure NullPlatformProvider implements PlatformProvider.
var _ PlatformProvider = NullPlatformProvider{}
