// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

package gpucontext

import "time"

// GestureEvent contains computed gesture deltas per frame.
//
// This event follows the Vello multi-touch pattern where gesture deltas
// are computed once per frame from the set of active pointers. This approach
// avoids jitter from individual pointer moves and provides smooth, predictable
// gesture values.
//
// The event is designed for multi-touch gestures (pinch-to-zoom, rotation, pan)
// but degrades gracefully with fewer pointers:
//   - 0-1 pointers: Empty event (no gesture possible)
//   - 2+ pointers: Full gesture with zoom, rotation, and translation
//
// Example usage:
//
//	source.OnGesture(func(ev gpucontext.GestureEvent) {
//	    if ev.NumPointers >= 2 {
//	        camera.Zoom(ev.ZoomDelta)
//	        camera.Rotate(ev.RotationDelta)
//	        camera.Pan(ev.TranslationDelta)
//	    }
//	})
type GestureEvent struct {
	// NumPointers is the number of active touch points.
	// Gestures require at least 2 pointers.
	NumPointers int

	// ZoomDelta is the proportional zoom factor for this frame.
	// 1.0 = no change, >1.0 = zoom in, <1.0 = zoom out.
	// Computed from change in average distance from centroid.
	ZoomDelta float64

	// ZoomDelta2D provides non-proportional zoom (stretch) deltas.
	// This allows independent X and Y scaling for non-uniform zoom.
	// For most use cases, use ZoomDelta instead.
	ZoomDelta2D Point

	// RotationDelta is the rotation change in radians for this frame.
	// Positive = counter-clockwise, negative = clockwise.
	// Computed from angle change of first pointer relative to centroid.
	RotationDelta float64

	// TranslationDelta is the pan movement in logical pixels for this frame.
	// Computed from change in centroid position.
	TranslationDelta Point

	// PinchType classifies the pinch gesture based on finger geometry.
	// Useful for constraining zoom to one axis (e.g., timeline scrubbing).
	PinchType PinchType

	// Center is the centroid of all active touch points.
	// Use this as the zoom/rotation pivot point.
	Center Point

	// Timestamp is the event time as duration since an arbitrary reference.
	// Useful for velocity calculations or animation timing.
	// Zero if timestamps are not available.
	Timestamp time.Duration
}

// PinchType classifies a two-finger pinch gesture based on finger geometry.
type PinchType uint8

const (
	// PinchNone indicates no pinch gesture (fewer than 2 pointers).
	PinchNone PinchType = iota

	// PinchHorizontal indicates horizontal separation exceeds vertical by 3x.
	// The fingers are spread horizontally, suggesting horizontal zoom/scrub.
	PinchHorizontal

	// PinchVertical indicates vertical separation exceeds horizontal by 3x.
	// The fingers are spread vertically, suggesting vertical zoom.
	PinchVertical

	// PinchProportional indicates uniform pinch (default).
	// Neither axis dominates, suggesting proportional zoom.
	PinchProportional
)

// String returns the pinch type name for debugging.
func (p PinchType) String() string {
	switch p {
	case PinchNone:
		return "None"
	case PinchHorizontal:
		return "Horizontal"
	case PinchVertical:
		return "Vertical"
	case PinchProportional:
		return "Proportional"
	default:
		return "Unknown"
	}
}

// Point represents a 2D coordinate in logical pixels.
type Point struct {
	X, Y float64
}

// Add returns the sum of two points.
func (p Point) Add(other Point) Point {
	return Point{X: p.X + other.X, Y: p.Y + other.Y}
}

// Sub returns the difference of two points.
func (p Point) Sub(other Point) Point {
	return Point{X: p.X - other.X, Y: p.Y - other.Y}
}

// Scale returns the point scaled by a factor.
func (p Point) Scale(factor float64) Point {
	return Point{X: p.X * factor, Y: p.Y * factor}
}

// GestureEventSource provides gesture event callbacks.
//
// This interface extends EventSource with high-level gesture recognition.
// The gesture recognizer computes deltas once per frame from pointer events,
// following the Vello pattern for smooth, predictable gestures.
//
// Type assertion pattern:
//
//	if ges, ok := eventSource.(gpucontext.GestureEventSource); ok {
//	    ges.OnGesture(handleGestureEvent)
//	}
//
// For applications that need gesture support:
//
//	ges.OnGesture(func(ev gpucontext.GestureEvent) {
//	    if ev.NumPointers >= 2 {
//	        handlePinchZoom(ev.ZoomDelta, ev.Center)
//	    }
//	})
type GestureEventSource interface {
	// OnGesture registers a callback for gesture events.
	// The callback receives a GestureEvent containing computed deltas.
	//
	// Callback threading: Called on the main/UI thread at end of frame.
	// Callbacks should be fast and non-blocking.
	//
	// Gesture events are delivered once per frame when 2+ pointers are active.
	OnGesture(fn func(GestureEvent))
}

// NullGestureEventSource implements GestureEventSource by ignoring all registrations.
// Useful for platforms or configurations where gesture input is not available.
type NullGestureEventSource struct{}

// OnGesture does nothing.
func (NullGestureEventSource) OnGesture(func(GestureEvent)) {}

// Ensure NullGestureEventSource implements GestureEventSource.
var _ GestureEventSource = NullGestureEventSource{}
