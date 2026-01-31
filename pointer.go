// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

package gpucontext

import "time"

// PointerEvent represents a unified pointer event following W3C Pointer Events Level 3.
//
// This event unifies mouse, touch, and pen input into a single abstraction,
// enabling applications to handle all pointing devices uniformly.
//
// W3C Pointer Events Level 3 Specification:
// https://www.w3.org/TR/pointerevents3/
//
// Key concepts:
//   - PointerID: Unique identifier for each active pointer
//   - PointerType: Distinguishes mouse, touch, and pen
//   - IsPrimary: Identifies the "main" pointer in multi-pointer scenarios
//   - Pressure/Tilt: Extended properties for pen/touch input
//
// Example usage:
//
//	source.OnPointer(func(ev gpucontext.PointerEvent) {
//	    switch ev.Type {
//	    case gpucontext.PointerDown:
//	        startDrag(ev.PointerID, ev.X, ev.Y)
//	    case gpucontext.PointerMove:
//	        if ev.Pressure > 0 {
//	            updateDrag(ev.PointerID, ev.X, ev.Y, ev.Pressure)
//	        }
//	    case gpucontext.PointerUp:
//	        endDrag(ev.PointerID)
//	    }
//	})
type PointerEvent struct {
	// Type indicates the type of pointer event (down, up, move, etc.).
	Type PointerEventType

	// PointerID uniquely identifies the pointer causing this event.
	// For mouse, this is typically 1. For touch/pen, each contact has a unique ID.
	// The ID remains constant from PointerDown through PointerUp/PointerCancel.
	PointerID int

	// X is the horizontal position relative to the window content area.
	// Uses logical pixels (CSS pixels equivalent).
	X float64

	// Y is the vertical position relative to the window content area.
	// Uses logical pixels (CSS pixels equivalent).
	Y float64

	// Pressure indicates the normalized pressure of the pointer input.
	// Range: 0.0 (no pressure) to 1.0 (maximum pressure).
	// For devices without pressure support (e.g., mouse), this is:
	//   - 0.5 when buttons are pressed
	//   - 0.0 when no buttons are pressed
	Pressure float32

	// TiltX is the plane angle between the Y-Z plane and the plane containing
	// the pointer axis and the Y axis, in degrees.
	// Range: -90 to 90 degrees.
	// Positive values tilt toward the right.
	// 0 when the pointer is perpendicular to the surface or not supported.
	TiltX float32

	// TiltY is the plane angle between the X-Z plane and the plane containing
	// the pointer axis and the X axis, in degrees.
	// Range: -90 to 90 degrees.
	// Positive values tilt toward the user.
	// 0 when the pointer is perpendicular to the surface or not supported.
	TiltY float32

	// Twist is the clockwise rotation of the pointer around its major axis,
	// in degrees.
	// Range: 0 to 359 degrees.
	// 0 when not supported.
	Twist float32

	// Width is the width of the contact geometry in logical pixels.
	// For devices that don't support contact geometry, this is 1.
	Width float32

	// Height is the height of the contact geometry in logical pixels.
	// For devices that don't support contact geometry, this is 1.
	Height float32

	// PointerType indicates the type of pointing device.
	PointerType PointerType

	// IsPrimary indicates if this is the primary pointer of its type.
	// For single-pointer devices (mouse), this is always true.
	// For multi-touch, the first finger down is primary.
	IsPrimary bool

	// Button indicates which button triggered this event.
	// Only meaningful for PointerDown and PointerUp events.
	// For PointerMove, this is ButtonNone.
	Button Button

	// Buttons is a bitmask of all currently pressed buttons.
	// This allows tracking multiple button states during movement.
	Buttons Buttons

	// Modifiers contains the keyboard modifier state at event time.
	Modifiers Modifiers

	// Timestamp is the event time as duration since an arbitrary reference.
	// Useful for calculating velocities and detecting double-clicks.
	// Zero if timestamps are not available on the platform.
	Timestamp time.Duration
}

// PointerEventType indicates the type of pointer event.
type PointerEventType uint8

const (
	// PointerDown is fired when a pointer becomes active.
	// For mouse: button press.
	// For touch: contact starts.
	// For pen: contact with or hover above the digitizer.
	PointerDown PointerEventType = iota

	// PointerUp is fired when a pointer is no longer active.
	// For mouse: button release.
	// For touch: contact ends.
	// For pen: leaves the digitizer detection range.
	PointerUp

	// PointerMove is fired when a pointer changes coordinates.
	// Also fired when pressure, tilt, or other properties change.
	PointerMove

	// PointerEnter is fired when a pointer enters the window bounds.
	// Does not bubble in W3C spec, but we deliver it directly.
	PointerEnter

	// PointerLeave is fired when a pointer leaves the window bounds.
	// Does not bubble in W3C spec, but we deliver it directly.
	PointerLeave

	// PointerCancel is fired when the system cancels the pointer.
	// This happens when:
	//   - The browser decides the pointer is unlikely to produce more events
	//   - A device orientation change occurs
	//   - The application loses focus during an active pointer
	// Always handle cancellation to reset state properly.
	PointerCancel
)

// String returns the event type name for debugging.
func (t PointerEventType) String() string {
	switch t {
	case PointerDown:
		return "PointerDown"
	case PointerUp:
		return "PointerUp"
	case PointerMove:
		return "PointerMove"
	case PointerEnter:
		return "PointerEnter"
	case PointerLeave:
		return "PointerLeave"
	case PointerCancel:
		return "PointerCancel"
	default:
		return "Unknown"
	}
}

// PointerType indicates the type of pointing device.
type PointerType uint8

const (
	// PointerTypeMouse indicates a mouse or similar device.
	// Includes trackpads when they emulate mouse behavior.
	PointerTypeMouse PointerType = iota

	// PointerTypeTouch indicates direct touch input.
	// Includes touchscreens and touch-enabled trackpads.
	PointerTypeTouch

	// PointerTypePen indicates a stylus or pen input.
	// Includes graphics tablets and pen-enabled displays.
	PointerTypePen
)

// String returns the pointer type name for debugging.
func (t PointerType) String() string {
	switch t {
	case PointerTypeMouse:
		return "Mouse"
	case PointerTypeTouch:
		return "Touch"
	case PointerTypePen:
		return "Pen"
	default:
		return "Unknown"
	}
}

// Button identifies which button triggered a pointer event.
// This follows the W3C Pointer Events specification button values.
type Button int8

const (
	// ButtonNone indicates no button or no change in button state.
	// Used for move events where no button triggered the event.
	ButtonNone Button = -1

	// ButtonLeft is the primary button (usually left mouse button).
	ButtonLeft Button = 0

	// ButtonMiddle is the auxiliary button (usually middle mouse button or wheel click).
	ButtonMiddle Button = 1

	// ButtonRight is the secondary button (usually right mouse button).
	ButtonRight Button = 2

	// ButtonX1 is the first extra button (usually "back" button).
	ButtonX1 Button = 3

	// ButtonX2 is the second extra button (usually "forward" button).
	ButtonX2 Button = 4

	// ButtonEraser is the eraser button on a pen (if available).
	ButtonEraser Button = 5
)

// String returns the button name for debugging.
func (b Button) String() string {
	switch b {
	case ButtonNone:
		return "None"
	case ButtonLeft:
		return "Left"
	case ButtonMiddle:
		return "Middle"
	case ButtonRight:
		return "Right"
	case ButtonX1:
		return "X1"
	case ButtonX2:
		return "X2"
	case ButtonEraser:
		return "Eraser"
	default:
		return "Unknown"
	}
}

// Buttons is a bitmask representing currently pressed buttons.
// This allows tracking multiple button states simultaneously.
type Buttons uint8

const (
	// ButtonsNone indicates no buttons are pressed.
	ButtonsNone Buttons = 0

	// ButtonsLeft indicates the left button is pressed.
	ButtonsLeft Buttons = 1 << 0

	// ButtonsRight indicates the right button is pressed.
	ButtonsRight Buttons = 1 << 1

	// ButtonsMiddle indicates the middle button is pressed.
	ButtonsMiddle Buttons = 1 << 2

	// ButtonsX1 indicates the X1 (back) button is pressed.
	ButtonsX1 Buttons = 1 << 3

	// ButtonsX2 indicates the X2 (forward) button is pressed.
	ButtonsX2 Buttons = 1 << 4

	// ButtonsEraser indicates the eraser button is pressed.
	ButtonsEraser Buttons = 1 << 5
)

// HasLeft returns true if the left button is pressed.
func (b Buttons) HasLeft() bool {
	return b&ButtonsLeft != 0
}

// HasRight returns true if the right button is pressed.
func (b Buttons) HasRight() bool {
	return b&ButtonsRight != 0
}

// HasMiddle returns true if the middle button is pressed.
func (b Buttons) HasMiddle() bool {
	return b&ButtonsMiddle != 0
}

// HasX1 returns true if the X1 (back) button is pressed.
func (b Buttons) HasX1() bool {
	return b&ButtonsX1 != 0
}

// HasX2 returns true if the X2 (forward) button is pressed.
func (b Buttons) HasX2() bool {
	return b&ButtonsX2 != 0
}

// HasEraser returns true if the eraser button is pressed.
func (b Buttons) HasEraser() bool {
	return b&ButtonsEraser != 0
}

// Count returns the number of pressed buttons.
func (b Buttons) Count() int {
	count := 0
	for v := b; v != 0; v &= v - 1 {
		count++
	}
	return count
}

// PointerEventSource extends EventSource with unified pointer event capabilities.
//
// This interface provides W3C Pointer Events Level 3 compliant pointer input,
// unifying mouse, touch, and pen input into a single event stream.
//
// Implementation note: Rather than adding to EventSource directly,
// we use a separate interface to maintain backward compatibility
// and allow type assertion:
//
//	if pes, ok := eventSource.(gpucontext.PointerEventSource); ok {
//	    pes.OnPointer(handlePointerEvent)
//	}
//
// For applications that need only unified pointer events:
//
//	pes.OnPointer(func(ev gpucontext.PointerEvent) {
//	    // Handle all pointer types uniformly
//	})
type PointerEventSource interface {
	// OnPointer registers a callback for pointer events.
	// The callback receives a PointerEvent containing all pointer information.
	//
	// Callback threading: Called on the main/UI thread.
	// Callbacks should be fast and non-blocking.
	//
	// Pointer events are delivered in order:
	//   PointerEnter -> PointerDown -> PointerMove* -> PointerUp/PointerCancel -> PointerLeave
	OnPointer(fn func(PointerEvent))
}

// NullPointerEventSource implements PointerEventSource by ignoring all registrations.
// Useful for platforms or configurations where pointer input is not available.
type NullPointerEventSource struct{}

// OnPointer does nothing.
func (NullPointerEventSource) OnPointer(func(PointerEvent)) {}

// Ensure NullPointerEventSource implements PointerEventSource.
var _ PointerEventSource = NullPointerEventSource{}
