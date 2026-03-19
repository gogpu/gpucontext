// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

package gpucontext

// HitTestResult represents what region of the window the cursor is over.
//
// When a window is frameless (no OS title bar), the application must tell
// the OS what part of the window the cursor is in, so the OS can handle
// dragging, resizing, and window button interactions.
//
// These values map directly to platform-specific hit test constants:
//   - Windows: WM_NCHITTEST return values (HTCLIENT, HTCAPTION, etc.)
//   - macOS: NSWindow regions
//   - Linux: xdg-toplevel resize edges
type HitTestResult int

const (
	// HitTestClient indicates the cursor is over the client (content) area.
	// The application handles all input normally.
	HitTestClient HitTestResult = iota

	// HitTestCaption indicates the cursor is over the title bar / drag area.
	// The OS handles window dragging on mouse down.
	HitTestCaption

	// HitTestClose indicates the cursor is over the close button region.
	HitTestClose

	// HitTestMaximize indicates the cursor is over the maximize button region.
	HitTestMaximize

	// HitTestMinimize indicates the cursor is over the minimize button region.
	HitTestMinimize

	// HitTestResizeN indicates the cursor is over the top resize edge.
	HitTestResizeN

	// HitTestResizeS indicates the cursor is over the bottom resize edge.
	HitTestResizeS

	// HitTestResizeW indicates the cursor is over the left resize edge.
	HitTestResizeW

	// HitTestResizeE indicates the cursor is over the right resize edge.
	HitTestResizeE

	// HitTestResizeNW indicates the cursor is over the top-left resize corner.
	HitTestResizeNW

	// HitTestResizeNE indicates the cursor is over the top-right resize corner.
	HitTestResizeNE

	// HitTestResizeSW indicates the cursor is over the bottom-left resize corner.
	HitTestResizeSW

	// HitTestResizeSE indicates the cursor is over the bottom-right resize corner.
	HitTestResizeSE
)

// String returns the hit test result name for debugging.
func (h HitTestResult) String() string {
	switch h {
	case HitTestClient:
		return "Client"
	case HitTestCaption:
		return "Caption"
	case HitTestClose:
		return "Close"
	case HitTestMaximize:
		return "Maximize"
	case HitTestMinimize:
		return "Minimize"
	case HitTestResizeN:
		return "ResizeN"
	case HitTestResizeS:
		return "ResizeS"
	case HitTestResizeW:
		return "ResizeW"
	case HitTestResizeE:
		return "ResizeE"
	case HitTestResizeNW:
		return "ResizeNW"
	case HitTestResizeNE:
		return "ResizeNE"
	case HitTestResizeSW:
		return "ResizeSW"
	case HitTestResizeSE:
		return "ResizeSE"
	default:
		return "Unknown"
	}
}

// HitTestCallback is called by the platform layer to determine what region
// of the window the cursor is in. The coordinates (x, y) are in logical
// points (DIP) relative to the window's top-left corner.
//
// The callback is invoked during mouse move events when the window is frameless.
// It must return quickly to avoid input lag.
type HitTestCallback func(x, y float64) HitTestResult

// WindowChrome provides control over window chrome (title bar, borders).
//
// This interface enables replacing the OS window chrome with a custom
// GPU-rendered title bar. When frameless mode is enabled, the OS removes
// the title bar and borders, and the application takes responsibility for:
//   - Rendering its own title bar
//   - Providing hit-test regions (drag area, buttons, resize edges)
//   - Handling minimize/maximize/close actions
//
// Implementations:
//   - gogpu.App implements WindowChrome via platform-specific code
//   - NullWindowChrome provides no-op defaults for testing
//
// WindowChrome is optional. Use type assertion to check availability:
//
//	if wc, ok := provider.(gpucontext.WindowChrome); ok {
//	    wc.SetFrameless(true)
//	    wc.SetHitTestCallback(myHitTest)
//	}
type WindowChrome interface {
	// SetFrameless enables or disables frameless (borderless) window mode.
	// When true, the OS title bar and borders are removed.
	// The application must provide its own title bar via SetHitTestCallback.
	SetFrameless(frameless bool)

	// IsFrameless returns true if the window is in frameless mode.
	IsFrameless() bool

	// SetHitTestCallback sets the callback that determines what region
	// of the window the cursor is over. This is used by the platform layer
	// to route mouse events to the OS for dragging, resizing, etc.
	//
	// Pass nil to clear the callback (all areas become HitTestClient).
	SetHitTestCallback(callback HitTestCallback)

	// Minimize minimizes the window to the taskbar/dock.
	Minimize()

	// Maximize toggles between maximized and restored window state.
	// If the window is currently maximized, it is restored to its previous size.
	Maximize()

	// IsMaximized returns true if the window is currently maximized.
	IsMaximized() bool

	// Close requests the window to close.
	// This triggers the normal close flow (close events, cleanup).
	Close()
}

// NullWindowChrome implements WindowChrome with no-op behavior.
// Used for testing and platforms without window chrome support.
//
// Default return values:
//   - IsFrameless: false
//   - IsMaximized: false
//   - All actions: no-op
type NullWindowChrome struct{}

// SetFrameless does nothing.
func (NullWindowChrome) SetFrameless(bool) {}

// IsFrameless returns false.
func (NullWindowChrome) IsFrameless() bool { return false }

// SetHitTestCallback does nothing.
func (NullWindowChrome) SetHitTestCallback(HitTestCallback) {}

// Minimize does nothing.
func (NullWindowChrome) Minimize() {}

// Maximize does nothing.
func (NullWindowChrome) Maximize() {}

// IsMaximized returns false.
func (NullWindowChrome) IsMaximized() bool { return false }

// Close does nothing.
func (NullWindowChrome) Close() {}

// Ensure NullWindowChrome implements WindowChrome.
var _ WindowChrome = NullWindowChrome{}
