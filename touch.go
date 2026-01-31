// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

package gpucontext

import "time"

// TouchID uniquely identifies a touch point within a touch session.
// The ID remains constant from TouchBegan through TouchEnded/TouchCanceled.
// IDs may be reused after a touch ends.
//
// Design note: Using int for compatibility with both int32 (Android) and int64 (iOS).
// Matches Ebitengine pattern for familiarity.
type TouchID int

// TouchPhase represents the lifecycle stage of a touch point.
// These phases align with platform conventions:
//   - Android: ACTION_DOWN/MOVE/UP/CANCEL
//   - iOS: touchesBegan/Moved/Ended/Canceled
//   - W3C: touchstart/move/end/cancel
type TouchPhase uint8

const (
	// TouchBegan indicates first contact with the touch surface.
	// Sent once per touch point at the start of interaction.
	TouchBegan TouchPhase = iota

	// TouchMoved indicates the touch point has moved.
	// Sent multiple times during drag/pan gestures.
	TouchMoved

	// TouchEnded indicates the touch point was lifted normally.
	// Sent once per touch point at the end of interaction.
	TouchEnded

	// TouchCanceled indicates the system interrupted the touch.
	// This can happen when:
	//   - The app loses focus
	//   - A system gesture is recognized (e.g., swipe to home)
	//   - The touch hardware reports an error
	// Always handle cancellation to reset UI state properly.
	TouchCanceled
)

// String returns the phase name for debugging.
func (p TouchPhase) String() string {
	switch p {
	case TouchBegan:
		return "Began"
	case TouchMoved:
		return "Moved"
	case TouchEnded:
		return "Ended"
	case TouchCanceled:
		return "Canceled"
	default:
		return "Unknown"
	}
}

// TouchPoint represents a single point of contact on a touch surface.
//
// Position coordinates are in logical pixels relative to the window's
// content area. The coordinate system matches the graphics system:
//   - Origin (0, 0) is at top-left
//   - X increases rightward
//   - Y increases downward
//
// Design decisions:
//   - Using float64 for sub-pixel precision (matches mouse events in EventSource)
//   - Pressure/Radius are pointers to indicate optional support
//   - No coordinate transformation - that's the UI layer's responsibility
type TouchPoint struct {
	// ID uniquely identifies this touch point within the session.
	// Track touches by ID, not by array index (indices can change).
	ID TouchID

	// X is the horizontal position in logical pixels.
	X float64

	// Y is the vertical position in logical pixels.
	Y float64

	// Pressure is the contact pressure, if supported by hardware.
	// Range: 0.0 (no pressure) to 1.0 (maximum pressure).
	// nil if pressure sensing is not available.
	//
	// Use case: Drawing apps, pressure-sensitive UI elements.
	Pressure *float32

	// Radius is the approximate contact radius in logical pixels.
	// Represents a circular approximation of the contact area.
	// nil if radius detection is not available.
	//
	// Use case: Distinguishing finger vs knuckle touches,
	// accessibility features for users with larger contact areas.
	Radius *float32
}

// TouchEvent represents a touch input event containing one or more touch points.
//
// Multi-touch handling:
//   - TouchBegan: Changed contains new touches, All contains all active touches
//   - TouchMoved: Changed contains moved touches, All contains all active touches
//   - TouchEnded: Changed contains lifted touches, All contains remaining touches
//   - TouchCanceled: Changed contains canceled touches, All may be empty
//
// Example multi-touch pinch gesture processing:
//
//	func handleTouch(ev gpucontext.TouchEvent) {
//	    if ev.Phase == gpucontext.TouchMoved && len(ev.All) == 2 {
//	        // Calculate distance between two fingers for pinch
//	        dx := ev.All[0].X - ev.All[1].X
//	        dy := ev.All[0].Y - ev.All[1].Y
//	        distance := math.Sqrt(dx*dx + dy*dy)
//	        // Use distance for zoom...
//	    }
//	}
type TouchEvent struct {
	// Phase indicates the lifecycle stage of the touches in Changed.
	Phase TouchPhase

	// Changed contains the touch points that triggered this event.
	// For TouchBegan: newly added touches
	// For TouchMoved: touches that moved
	// For TouchEnded: touches that were lifted
	// For TouchCanceled: touches that were interrupted
	Changed []TouchPoint

	// All contains all currently active touch points.
	// Useful for multi-touch gestures that need to track all contacts.
	// For TouchEnded/TouchCanceled, this excludes the Changed touches.
	All []TouchPoint

	// Modifiers contains keyboard modifier state at the time of the event.
	// Useful for modifier+touch combinations (e.g., Ctrl+drag for zoom).
	Modifiers Modifiers

	// Timestamp is the event time as duration since an arbitrary reference.
	// Useful for calculating velocities in fling gestures.
	// Zero if timestamps are not available on the platform.
	Timestamp time.Duration
}

// TouchEventSource extends EventSource with touch input capabilities.
// This interface is optional - not all EventSource implementations
// support touch input (e.g., desktop-only applications).
//
// Implementation note: Rather than adding to EventSource directly,
// we use a separate interface to maintain backward compatibility
// and allow type assertion:
//
//	if tes, ok := eventSource.(gpucontext.TouchEventSource); ok {
//	    tes.OnTouch(handleTouchEvent)
//	}
type TouchEventSource interface {
	// OnTouch registers a callback for touch events.
	// The callback receives a TouchEvent containing all touch information.
	//
	// Callback threading: Called on the main/UI thread.
	// Callbacks should be fast and non-blocking.
	//
	// Touch events are delivered in order: Began -> Moved* -> Ended/Canceled
	// Multi-touch events for simultaneous contacts are coalesced into single events.
	OnTouch(fn func(TouchEvent))
}

// NullTouchEventSource implements TouchEventSource by ignoring all registrations.
// Useful for platforms or configurations where touch input is not available.
type NullTouchEventSource struct{}

// OnTouch does nothing.
func (NullTouchEventSource) OnTouch(func(TouchEvent)) {}

// Ensure NullTouchEventSource implements TouchEventSource.
var _ TouchEventSource = NullTouchEventSource{}
