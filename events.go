// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

package gpucontext

import "sync"

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
	//
	// All mouse/pointer coordinates (x, y) are in logical DIP (device-independent
	// pixels), consistent with WindowProvider.Size(). Do NOT divide by ScaleFactor —
	// the framework applies DPI scaling internally on all platforms.

	// OnMouseMove registers a callback for mouse movement.
	// x, y are in logical DIP relative to the window content area.
	OnMouseMove(func(x, y float64))

	// OnMousePress registers a callback for mouse button press.
	// x, y are in logical DIP relative to the window content area.
	OnMousePress(func(button MouseButton, x, y float64))

	// OnMouseRelease registers a callback for mouse button release.
	// x, y are in logical DIP relative to the window content area.
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
//
// Values use grouped ranges with explicit base offsets (the net/http pattern).
// Each group reserves room for future expansion without shifting existing values.
// This is a platform-independent virtual key code scheme; platform code maps
// native scan codes / virtual keys to these values.
//
// Groups and ranges:
//
//	KeyUnknown    = 0
//	Letters       [1..31]    — A-Z (26 used, 5 reserved)
//	Digits        [33..47]   — 0-9 (10 used, 5 reserved)
//	Function      [49..80]   — F1-F24 (24 used, 8 reserved)
//	Navigation    [81..112]  — arrows, home, end, etc. (15 used, 17 reserved)
//	Modifiers     [113..128] — shift, ctrl, alt, super (8 used, 8 reserved)
//	Punctuation   [129..160] — brackets, operators, intl (13 used, 19 reserved)
//	Numpad        [161..192] — numpad digits + operators (18 used, 14 reserved)
//	Locks         [193..208] — caps, scroll, num lock (5 used, 11 reserved)
//	Media         [209..240] — playback controls (5 used, 27 reserved)
//	Volume        [241..248] — volume up/down/mute (3 used, 5 reserved)
//	Browser       [249..264] — navigation keys (5 used, 11 reserved)
//	System        [265..280] — context menu, cancel, launch (4 used, 12 reserved)
type Key uint16

// KeyUnknown represents an unrecognized or unmapped key.
const KeyUnknown Key = 0

// Letters [1..31] — 26 keys (A-Z), 5 reserved for future use.
const (
	KeyA Key = iota + 1
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
)

// Digits [33..47] — 10 keys (0-9), 5 reserved for future use.
const (
	Key0 Key = iota + 33
	Key1
	Key2
	Key3
	Key4
	Key5
	Key6
	Key7
	Key8
	Key9
)

// Function keys [49..80] — 24 keys (F1-F24), 8 reserved for future use.
const (
	KeyF1 Key = iota + 49
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
	KeyF13
	KeyF14
	KeyF15
	KeyF16
	KeyF17
	KeyF18
	KeyF19
	KeyF20
	KeyF21
	KeyF22
	KeyF23
	KeyF24
)

// Navigation [81..112] — 15 keys, 17 reserved for future use.
const (
	KeyEscape Key = iota + 81
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
)

// Modifiers as keys [113..128] — 8 keys, 8 reserved for future use.
// These represent physical modifier keys. For modifier state in event
// callbacks, use [Modifiers] flags instead.
const (
	KeyLeftShift Key = iota + 113
	KeyRightShift
	KeyLeftControl
	KeyRightControl
	KeyLeftAlt
	KeyRightAlt
	KeyLeftSuper
	KeyRightSuper
)

// Punctuation and symbols [129..160] — 13 keys, 19 reserved for future use.
const (
	KeyMinus Key = iota + 129
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
	KeyIntlBackslash // ISO 102nd key (between left Shift and Z on ISO layouts).
	KeyIntlYen       // JIS Yen key.
)

// Numpad [161..192] — 18 keys, 14 reserved for future use.
const (
	KeyNumpad0 Key = iota + 161
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
	KeyNumpadEqual // Numpad = (Mac keyboards, some international layouts).
	KeyNumpadComma // Numpad , (Brazilian ABNT2 layout).
)

// Lock keys [193..208] — 5 keys, 11 reserved for future use.
const (
	KeyCapsLock Key = iota + 193
	KeyScrollLock
	KeyNumLock
	KeyPrintScreen
	KeyPause
)

// Media keys [209..240] — 5 keys, 27 reserved for future use.
// Names follow the W3C UIEvents KeyboardEvent.code convention.
const (
	KeyMediaPlayPause Key = iota + 209
	KeyMediaStop
	KeyMediaTrackNext
	KeyMediaTrackPrevious
	KeyMediaRecord
)

// Volume keys [241..248] — 3 keys, 5 reserved for future use.
// Names follow the W3C UIEvents KeyboardEvent.code convention.
const (
	KeyAudioVolumeUp Key = iota + 241
	KeyAudioVolumeDown
	KeyAudioVolumeMute
)

// Browser keys [249..264] — 5 keys, 11 reserved for future use.
const (
	KeyBrowserBack Key = iota + 249
	KeyBrowserForward
	KeyBrowserRefresh
	KeyBrowserHome
	KeyBrowserSearch
)

// System keys [265..280] — 4 keys, 12 reserved for future use.
const (
	KeyContextMenu Key = iota + 265 // Application/context menu key (not VK_MENU/Alt).
	KeyCancel                       // Cancel key (Ctrl+Break on Windows).
	KeyLaunchApp1                   // Launch application 1 (typically My Computer).
	KeyLaunchApp2                   // Launch application 2 (typically Calculator).
)

// String returns a human-readable name for the key.
func (k Key) String() string {
	switch k {
	case KeyUnknown:
		return "Unknown"
	case KeyA:
		return "A"
	case KeyB:
		return "B"
	case KeyC:
		return "C"
	case KeyD:
		return "D"
	case KeyE:
		return "E"
	case KeyF:
		return "F"
	case KeyG:
		return "G"
	case KeyH:
		return "H"
	case KeyI:
		return "I"
	case KeyJ:
		return "J"
	case KeyK:
		return "K"
	case KeyL:
		return "L"
	case KeyM:
		return "M"
	case KeyN:
		return "N"
	case KeyO:
		return "O"
	case KeyP:
		return "P"
	case KeyQ:
		return "Q"
	case KeyR:
		return "R"
	case KeyS:
		return "S"
	case KeyT:
		return "T"
	case KeyU:
		return "U"
	case KeyV:
		return "V"
	case KeyW:
		return "W"
	case KeyX:
		return "X"
	case KeyY:
		return "Y"
	case KeyZ:
		return "Z"
	case Key0:
		return "0"
	case Key1:
		return "1"
	case Key2:
		return "2"
	case Key3:
		return "3"
	case Key4:
		return "4"
	case Key5:
		return "5"
	case Key6:
		return "6"
	case Key7:
		return "7"
	case Key8:
		return "8"
	case Key9:
		return "9"
	case KeyF1:
		return "F1"
	case KeyF2:
		return "F2"
	case KeyF3:
		return "F3"
	case KeyF4:
		return "F4"
	case KeyF5:
		return "F5"
	case KeyF6:
		return "F6"
	case KeyF7:
		return "F7"
	case KeyF8:
		return "F8"
	case KeyF9:
		return "F9"
	case KeyF10:
		return "F10"
	case KeyF11:
		return "F11"
	case KeyF12:
		return "F12"
	case KeyF13:
		return "F13"
	case KeyF14:
		return "F14"
	case KeyF15:
		return "F15"
	case KeyF16:
		return "F16"
	case KeyF17:
		return "F17"
	case KeyF18:
		return "F18"
	case KeyF19:
		return "F19"
	case KeyF20:
		return "F20"
	case KeyF21:
		return "F21"
	case KeyF22:
		return "F22"
	case KeyF23:
		return "F23"
	case KeyF24:
		return "F24"
	case KeyEscape:
		return "Escape"
	case KeyTab:
		return "Tab"
	case KeyBackspace:
		return "Backspace"
	case KeyEnter:
		return "Enter"
	case KeySpace:
		return "Space"
	case KeyInsert:
		return "Insert"
	case KeyDelete:
		return "Delete"
	case KeyHome:
		return "Home"
	case KeyEnd:
		return "End"
	case KeyPageUp:
		return "PageUp"
	case KeyPageDown:
		return "PageDown"
	case KeyLeft:
		return "Left"
	case KeyRight:
		return "Right"
	case KeyUp:
		return "Up"
	case KeyDown:
		return "Down"
	case KeyLeftShift:
		return "LeftShift"
	case KeyRightShift:
		return "RightShift"
	case KeyLeftControl:
		return "LeftControl"
	case KeyRightControl:
		return "RightControl"
	case KeyLeftAlt:
		return "LeftAlt"
	case KeyRightAlt:
		return "RightAlt"
	case KeyLeftSuper:
		return "LeftSuper"
	case KeyRightSuper:
		return "RightSuper"
	case KeyMinus:
		return "Minus"
	case KeyEqual:
		return "Equal"
	case KeyLeftBracket:
		return "LeftBracket"
	case KeyRightBracket:
		return "RightBracket"
	case KeyBackslash:
		return "Backslash"
	case KeySemicolon:
		return "Semicolon"
	case KeyApostrophe:
		return "Apostrophe"
	case KeyGrave:
		return "Grave"
	case KeyComma:
		return "Comma"
	case KeyPeriod:
		return "Period"
	case KeySlash:
		return "Slash"
	case KeyIntlBackslash:
		return "IntlBackslash"
	case KeyIntlYen:
		return "IntlYen"
	case KeyNumpad0:
		return "Numpad0"
	case KeyNumpad1:
		return "Numpad1"
	case KeyNumpad2:
		return "Numpad2"
	case KeyNumpad3:
		return "Numpad3"
	case KeyNumpad4:
		return "Numpad4"
	case KeyNumpad5:
		return "Numpad5"
	case KeyNumpad6:
		return "Numpad6"
	case KeyNumpad7:
		return "Numpad7"
	case KeyNumpad8:
		return "Numpad8"
	case KeyNumpad9:
		return "Numpad9"
	case KeyNumpadDecimal:
		return "NumpadDecimal"
	case KeyNumpadDivide:
		return "NumpadDivide"
	case KeyNumpadMultiply:
		return "NumpadMultiply"
	case KeyNumpadSubtract:
		return "NumpadSubtract"
	case KeyNumpadAdd:
		return "NumpadAdd"
	case KeyNumpadEnter:
		return "NumpadEnter"
	case KeyNumpadEqual:
		return "NumpadEqual"
	case KeyNumpadComma:
		return "NumpadComma"
	case KeyCapsLock:
		return "CapsLock"
	case KeyScrollLock:
		return "ScrollLock"
	case KeyNumLock:
		return "NumLock"
	case KeyPrintScreen:
		return "PrintScreen"
	case KeyPause:
		return "Pause"
	case KeyMediaPlayPause:
		return "MediaPlayPause"
	case KeyMediaStop:
		return "MediaStop"
	case KeyMediaTrackNext:
		return "MediaTrackNext"
	case KeyMediaTrackPrevious:
		return "MediaTrackPrevious"
	case KeyMediaRecord:
		return "MediaRecord"
	case KeyAudioVolumeUp:
		return "AudioVolumeUp"
	case KeyAudioVolumeDown:
		return "AudioVolumeDown"
	case KeyAudioVolumeMute:
		return "AudioVolumeMute"
	case KeyBrowserBack:
		return "BrowserBack"
	case KeyBrowserForward:
		return "BrowserForward"
	case KeyBrowserRefresh:
		return "BrowserRefresh"
	case KeyBrowserHome:
		return "BrowserHome"
	case KeyBrowserSearch:
		return "BrowserSearch"
	case KeyContextMenu:
		return "ContextMenu"
	case KeyCancel:
		return "Cancel"
	case KeyLaunchApp1:
		return "LaunchApp1"
	case KeyLaunchApp2:
		return "LaunchApp2"
	default:
		return "Key(" + uitoa(uint(k)) + ")"
	}
}

// uitoa converts a uint to its string representation without importing strconv.
func uitoa(val uint) string {
	if val == 0 {
		return "0"
	}
	var buf [20]byte // big enough for 64-bit uint
	i := len(buf) - 1
	for val > 0 {
		buf[i] = byte('0' + val%10)
		val /= 10
		i--
	}
	return string(buf[i+1:])
}

var (
	keyStringOnce sync.Once
	keyStringMap  map[string]Key
)

func buildKeyStringMap() {
	m := make(map[string]Key, 140)
	for k := Key(0); k <= KeyLaunchApp2; k++ {
		s := k.String()
		if s != "" && s[0] != 'K' {
			m[s] = k
		}
	}
	keyStringMap = m
}

// KeyFromString returns the Key for a given string name.
// The name must match the value returned by [Key.String] (e.g., "A", "F13",
// "MediaPlayPause", "ContextMenu"). Returns (KeyUnknown, false) if the name
// is not recognized.
//
// This function is safe for concurrent use and enables W3C
// KeyboardEvent.code compatibility: browser platforms can convert
// JavaScript event.code strings directly to Key values.
//
//	key, ok := gpucontext.KeyFromString("A")              // → KeyA, true
//	key, ok := gpucontext.KeyFromString("MediaPlayPause") // → KeyMediaPlayPause, true
//	key, ok := gpucontext.KeyFromString("nonexistent")    // → KeyUnknown, false
func KeyFromString(name string) (Key, bool) {
	keyStringOnce.Do(buildKeyStringMap)
	k, ok := keyStringMap[name]
	return k, ok
}

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
