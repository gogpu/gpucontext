// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: BSD-3-Clause

package gpucontext

// EventSource provides input events from the host application to UI frameworks.
//
// This interface enables UI frameworks (like gogpu/ui) to receive user input
// events from the host window system. The host application (e.g., gogpu.App)
// implements EventSource and passes it to the UI layer.
//
// Event callbacks are invoked on the main thread during the event loop.
// Callback functions should be fast and non-blocking.
//
// Example usage in a UI framework:
//
//	func (ui *UI) AttachEvents(source gpucontext.EventSource) {
//	    source.OnMousePress(func(button MouseButton, x, y float64) {
//	        widget := ui.hitTest(x, y)
//	        if widget != nil {
//	            widget.HandleMouseDown(button, x, y)
//	        }
//	    })
//
//	    source.OnKeyPress(func(key Key, mods Modifiers) {
//	        ui.focused.HandleKeyDown(key, mods)
//	    })
//	}
//
// Note: This interface is designed for gogpu ↔ ui integration.
// The rendering library (gg) does NOT use this interface.
type EventSource interface {
	// Keyboard events

	// OnKeyPress registers a callback for key press events.
	OnKeyPress(func(key Key, mods Modifiers))

	// OnKeyRelease registers a callback for key release events.
	OnKeyRelease(func(key Key, mods Modifiers))

	// OnTextInput registers a callback for text input events.
	// Text input is the result of key presses after applying keyboard layouts
	// and input methods. This is the preferred way to handle text entry.
	OnTextInput(func(text string))

	// Mouse events

	// OnMouseMove registers a callback for mouse movement.
	OnMouseMove(func(x, y float64))

	// OnMousePress registers a callback for mouse button press.
	OnMousePress(func(button MouseButton, x, y float64))

	// OnMouseRelease registers a callback for mouse button release.
	OnMouseRelease(func(button MouseButton, x, y float64))

	// OnScroll registers a callback for scroll wheel events.
	// dx and dy are the scroll deltas (positive = right/down).
	OnScroll(func(dx, dy float64))

	// Window events

	// OnResize registers a callback for window resize.
	OnResize(func(width, height int))

	// OnFocus registers a callback for focus change.
	OnFocus(func(focused bool))

	// IME events for international text input

	// OnIMECompositionStart registers a callback for when IME composition begins.
	// This is called when the user starts typing in an IME (e.g., for CJK input).
	OnIMECompositionStart(fn func())

	// OnIMECompositionUpdate registers a callback for IME composition updates.
	// Called during composition with the current state (preview text, cursor).
	OnIMECompositionUpdate(fn func(state IMEState))

	// OnIMECompositionEnd registers a callback for when IME composition ends.
	// The committed parameter contains the final text that should be inserted.
	OnIMECompositionEnd(fn func(committed string))
}

// IMEState represents the current state of the Input Method Editor.
// This is used for CJK (Chinese, Japanese, Korean) and other complex text input.
//
// During IME composition, the user types phonetic characters that are converted
// to ideographic characters. The IMEState contains the current preview text
// and cursor information for rendering the composition inline.
type IMEState struct {
	// Composing indicates whether IME is currently in composition mode.
	Composing bool

	// CompositionText is the text currently being composed (e.g., pinyin for Chinese).
	// This should be displayed inline at the cursor position with special styling.
	CompositionText string

	// CursorPos is the cursor position within the composition text.
	CursorPos int

	// SelectionStart is the start of the selection within the composition text.
	// This is used for candidate selection in some IME systems.
	SelectionStart int

	// SelectionEnd is the end of the selection within the composition text.
	SelectionEnd int
}

// IMEController allows widgets to control IME behavior.
// This interface is typically implemented by the host window system.
type IMEController interface {
	// SetIMEPosition tells the platform where to show the IME candidate window.
	// The coordinates are in screen pixels relative to the window.
	SetIMEPosition(x, y int)

	// SetIMEEnabled enables or disables IME for the current input context.
	// When disabled, key presses are delivered directly without IME processing.
	// This is useful for password fields or non-text inputs.
	SetIMEEnabled(enabled bool)
}

// Key represents a keyboard key.
// Values follow a platform-independent virtual key code scheme.
type Key uint16

// Common key codes.
// These match typical USB HID usage codes for cross-platform compatibility.
const (
	KeyUnknown Key = iota

	// Letters
	KeyA
	KeyB
	KeyC
	KeyD
	KeyE
	KeyF
	KeyG
	KeyH
	KeyI
	KeyJ
	KeyK
	KeyL
	KeyM
	KeyN
	KeyO
	KeyP
	KeyQ
	KeyR
	KeyS
	KeyT
	KeyU
	KeyV
	KeyW
	KeyX
	KeyY
	KeyZ

	// Numbers
	Key0
	Key1
	Key2
	Key3
	Key4
	Key5
	Key6
	Key7
	Key8
	Key9

	// Function keys
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12

	// Navigation
	KeyEscape
	KeyTab
	KeyBackspace
	KeyEnter
	KeySpace
	KeyInsert
	KeyDelete
	KeyHome
	KeyEnd
	KeyPageUp
	KeyPageDown
	KeyLeft
	KeyRight
	KeyUp
	KeyDown

	// Modifiers (as keys, not modifiers)
	KeyLeftShift
	KeyRightShift
	KeyLeftControl
	KeyRightControl
	KeyLeftAlt
	KeyRightAlt
	KeyLeftSuper
	KeyRightSuper

	// Punctuation
	KeyMinus
	KeyEqual
	KeyLeftBracket
	KeyRightBracket
	KeyBackslash
	KeySemicolon
	KeyApostrophe
	KeyGrave
	KeyComma
	KeyPeriod
	KeySlash

	// Numpad
	KeyNumpad0
	KeyNumpad1
	KeyNumpad2
	KeyNumpad3
	KeyNumpad4
	KeyNumpad5
	KeyNumpad6
	KeyNumpad7
	KeyNumpad8
	KeyNumpad9
	KeyNumpadDecimal
	KeyNumpadDivide
	KeyNumpadMultiply
	KeyNumpadSubtract
	KeyNumpadAdd
	KeyNumpadEnter

	// Other
	KeyCapsLock
	KeyScrollLock
	KeyNumLock
	KeyPrintScreen
	KeyPause
)

// Modifiers represents keyboard modifier keys.
type Modifiers uint8

const (
	// ModShift indicates the Shift key is pressed.
	ModShift Modifiers = 1 << iota

	// ModControl indicates the Control key is pressed.
	ModControl

	// ModAlt indicates the Alt key is pressed (Option on macOS).
	ModAlt

	// ModSuper indicates the Super key is pressed (Windows/Command).
	ModSuper

	// ModCapsLock indicates Caps Lock is active.
	ModCapsLock

	// ModNumLock indicates Num Lock is active.
	ModNumLock
)

// HasShift returns true if the Shift modifier is set.
func (m Modifiers) HasShift() bool {
	return m&ModShift != 0
}

// HasControl returns true if the Control modifier is set.
func (m Modifiers) HasControl() bool {
	return m&ModControl != 0
}

// HasAlt returns true if the Alt modifier is set.
func (m Modifiers) HasAlt() bool {
	return m&ModAlt != 0
}

// HasSuper returns true if the Super modifier is set.
func (m Modifiers) HasSuper() bool {
	return m&ModSuper != 0
}

// MouseButton represents a mouse button.
type MouseButton uint8

const (
	// MouseButtonLeft is the primary mouse button.
	MouseButtonLeft MouseButton = iota

	// MouseButtonRight is the secondary mouse button.
	MouseButtonRight

	// MouseButtonMiddle is the middle mouse button (scroll wheel click).
	MouseButtonMiddle

	// MouseButton4 is an extra mouse button.
	MouseButton4

	// MouseButton5 is an extra mouse button.
	MouseButton5
)

// NullEventSource is an EventSource that ignores all event registrations.
// Used when events are not needed.
type NullEventSource struct{}

// OnKeyPress does nothing.
func (NullEventSource) OnKeyPress(func(Key, Modifiers)) {}

// OnKeyRelease does nothing.
func (NullEventSource) OnKeyRelease(func(Key, Modifiers)) {}

// OnTextInput does nothing.
func (NullEventSource) OnTextInput(func(string)) {}

// OnMouseMove does nothing.
func (NullEventSource) OnMouseMove(func(float64, float64)) {}

// OnMousePress does nothing.
func (NullEventSource) OnMousePress(func(MouseButton, float64, float64)) {}

// OnMouseRelease does nothing.
func (NullEventSource) OnMouseRelease(func(MouseButton, float64, float64)) {}

// OnScroll does nothing.
func (NullEventSource) OnScroll(func(float64, float64)) {}

// OnResize does nothing.
func (NullEventSource) OnResize(func(int, int)) {}

// OnFocus does nothing.
func (NullEventSource) OnFocus(func(bool)) {}

// OnIMECompositionStart does nothing.
func (NullEventSource) OnIMECompositionStart(func()) {}

// OnIMECompositionUpdate does nothing.
func (NullEventSource) OnIMECompositionUpdate(func(IMEState)) {}

// OnIMECompositionEnd does nothing.
func (NullEventSource) OnIMECompositionEnd(func(string)) {}

// Ensure NullEventSource implements EventSource.
var _ EventSource = NullEventSource{}
