// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

package gpucontext

// InputEvent represents a discrete input event from the platform.
//
// InputEvent is a sealed interface — only types defined in this package
// implement it. This enables exhaustive handling via linters and guarantees
// the set of event types is known at compile time.
//
// Use a Go type switch to handle specific event types:
//
//	for ev, ok := app.PollInputEvent(); ok; ev, ok = app.PollInputEvent() {
//	    switch e := ev.(type) {
//	    case gpucontext.KeyEvent:
//	        if e.Pressed && e.Key == gpucontext.KeyEscape {
//	            app.Quit()
//	        }
//	    case gpucontext.PointerEvent:
//	        if e.Type == gpucontext.PointerDown {
//	            handleClick(e.X, e.Y)
//	        }
//	    case gpucontext.ScrollEvent:
//	        handleScroll(e.DeltaX, e.DeltaY)
//	    }
//	}
//
// Three input models coexist in gogpu — use whichever fits your use case:
//   - Callbacks (EventSource): best for GUI frameworks
//   - State polling (app.Input()): best for simple games (Ebiten-style)
//   - Event queue (app.PollInputEvent()): best for complex games (SDL-style)
//
// Design references:
//   - Qt6: QEvent in QtCore (shared core layer)
//   - Bevy: bevy_input crate (separate from bevy_app)
//   - Gio: io/event.Event with ImplementsEvent() marker
//   - SDL3: SDL_PollEvent() + SDL_AppEvent() dual model
type InputEvent interface {
	inputEventTag()
}

// KeyEvent represents a keyboard key state change.
//
// A KeyEvent is emitted for each physical key press and release. For text
// input (after keyboard layout and IME processing), use CharEvent instead.
//
// All coordinates and dimensions in the event system use logical DIP
// (device-independent pixels), consistent with WindowProvider.Size().
type KeyEvent struct {
	Key     Key
	Mods    Modifiers
	Pressed bool // true = key down, false = key up
}

func (KeyEvent) inputEventTag() {}

// CharEvent represents committed text input.
//
// CharEvent is emitted after keyboard layout and input method processing.
// For CJK input, the composition preview is delivered via IME callbacks
// on EventSource; CharEvent contains only the final committed character.
//
// Use CharEvent for text fields. Use KeyEvent for keyboard shortcuts.
type CharEvent struct {
	Char rune
}

func (CharEvent) inputEventTag() {}

// FocusEvent represents a window focus state change.
type FocusEvent struct {
	Focused bool // true = window gained focus, false = lost focus
}

func (FocusEvent) inputEventTag() {}

// ResizeEvent represents a window content area size change.
//
// Width and Height are in logical DIP (device-independent pixels),
// consistent with WindowProvider.Size(). For physical pixel dimensions,
// multiply by ScaleFactor.
type ResizeEvent struct {
	Width  int
	Height int
}

func (ResizeEvent) inputEventTag() {}

// ScaleChangedEvent represents a DPI scale factor change.
//
// Emitted when a window moves between monitors with different DPI (e.g.,
// Retina 2x → FullHD 1x), or when the OS DPI setting changes at runtime.
//
// ScaleFactor is the new display scale (1.0 = standard, 2.0 = Retina/HiDPI).
// Width and Height are the new logical size in DIP — may change during DPI
// transition because the OS adjusts the window to preserve visual size.
//
// Event ordering: ScaleChangedEvent is emitted BEFORE ResizeEvent (cause
// before effect, winit pattern). Consumers should update their rendering
// scale before processing the new dimensions.
type ScaleChangedEvent struct {
	ScaleFactor float64
	Width       int
	Height      int
}

func (ScaleChangedEvent) inputEventTag() {}
