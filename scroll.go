// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

package gpucontext

import "time"

// ScrollEvent represents a scroll wheel or touchpad scroll event.
//
// This event is separate from PointerEvent because scroll events have
// different semantics:
//   - They don't have a persistent pointer ID
//   - They provide delta values rather than absolute positions
//   - They may have different units (lines, pages, pixels)
//
// For touchpad gestures, consider using GestureEvent (if available)
// for pinch-to-zoom and other multi-finger gestures.
//
// Example usage:
//
//	source.OnScroll(func(ev gpucontext.ScrollEvent) {
//	    scrollY += ev.DeltaY * scrollSpeed
//	    if ev.Modifiers.HasControl() {
//	        // Ctrl+scroll for zoom
//	        zoom *= 1.0 + ev.DeltaY * 0.1
//	    }
//	})
type ScrollEvent struct {
	// X is the pointer horizontal position at the time of scrolling.
	// Uses logical pixels relative to the window content area.
	X float64

	// Y is the pointer vertical position at the time of scrolling.
	// Uses logical pixels relative to the window content area.
	Y float64

	// DeltaX is the horizontal scroll amount.
	// Positive values scroll content to the right (or indicate rightward intent).
	// The unit depends on DeltaMode.
	DeltaX float64

	// DeltaY is the vertical scroll amount.
	// Positive values scroll content down (or indicate downward intent).
	// The unit depends on DeltaMode.
	DeltaY float64

	// DeltaMode indicates the unit of the delta values.
	DeltaMode ScrollDeltaMode

	// Modifiers contains the keyboard modifier state at event time.
	// Commonly used for Ctrl+scroll zoom behavior.
	Modifiers Modifiers

	// Timestamp is the event time as duration since an arbitrary reference.
	// Useful for smooth scrolling animations.
	// Zero if timestamps are not available on the platform.
	Timestamp time.Duration
}

// ScrollDeltaMode indicates the unit of scroll delta values.
type ScrollDeltaMode uint8

const (
	// ScrollDeltaPixel indicates delta values are in logical pixels.
	// This is typical for touchpad scrolling.
	ScrollDeltaPixel ScrollDeltaMode = iota

	// ScrollDeltaLine indicates delta values are in lines.
	// This is typical for traditional mouse wheel scrolling.
	// Applications should convert to pixels using their line height.
	ScrollDeltaLine

	// ScrollDeltaPage indicates delta values are in pages.
	// This is rare but may occur for Page Up/Down emulation.
	// Applications should convert to pixels using their viewport height.
	ScrollDeltaPage
)

// String returns the delta mode name for debugging.
func (m ScrollDeltaMode) String() string {
	switch m {
	case ScrollDeltaPixel:
		return "Pixel"
	case ScrollDeltaLine:
		return "Line"
	case ScrollDeltaPage:
		return "Page"
	default:
		return "Unknown"
	}
}

// ScrollEventSource extends EventSource with scroll event capabilities.
//
// This interface provides detailed scroll events with position, delta mode,
// and timing information beyond what the basic EventSource.OnScroll provides.
//
// For basic scroll handling, EventSource.OnScroll(dx, dy) is sufficient.
// Use ScrollEventSource when you need:
//   - Pointer position at scroll time
//   - Delta mode (pixels vs lines vs pages)
//   - Timestamps for smooth scrolling
//
// Type assertion pattern:
//
//	if ses, ok := eventSource.(gpucontext.ScrollEventSource); ok {
//	    ses.OnScrollEvent(handleScrollEvent)
//	} else {
//	    // Fall back to basic scroll handler
//	    eventSource.OnScroll(handleBasicScroll)
//	}
type ScrollEventSource interface {
	// OnScrollEvent registers a callback for detailed scroll events.
	// The callback receives a ScrollEvent containing all scroll information.
	//
	// Callback threading: Called on the main/UI thread.
	// Callbacks should be fast and non-blocking.
	OnScrollEvent(fn func(ScrollEvent))
}

// NullScrollEventSource implements ScrollEventSource by ignoring all registrations.
// Useful for platforms or configurations where scroll input is not available.
type NullScrollEventSource struct{}

// OnScrollEvent does nothing.
func (NullScrollEventSource) OnScrollEvent(func(ScrollEvent)) {}

// Ensure NullScrollEventSource implements ScrollEventSource.
var _ ScrollEventSource = NullScrollEventSource{}
