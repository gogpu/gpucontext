// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

package gpucontext

// WindowProvider provides window geometry and DPI information.
//
// This interface enables UI frameworks (like gogpu/ui) to query window
// dimensions and scale factor for layout calculations in density-independent
// pixels (Dp), and to request redraws for on-demand rendering.
//
// Implementations:
//   - gogpu.App implements WindowProvider via platform window
//   - NullWindowProvider provides configurable defaults for testing
//
// Example usage in a UI framework:
//
//	func (ui *UI) Layout(wp gpucontext.WindowProvider) {
//	    w, h := wp.Size()
//	    scale := wp.ScaleFactor()
//	    dpW := float64(w) / scale
//	    dpH := float64(h) / scale
//	    ui.root.Layout(dpW, dpH)
//	}
//
// Note: This interface is designed for gogpu <-> ui integration.
// The rendering library (gg) does NOT use this interface.
type WindowProvider interface {
	// Size returns the current window client area dimensions in physical pixels.
	Size() (width, height int)

	// ScaleFactor returns the DPI scale factor.
	// 1.0 = standard (96 DPI on Windows, 72 on macOS), 2.0 = Retina/HiDPI.
	// Used to convert between physical pixels and density-independent pixels (Dp).
	ScaleFactor() float64

	// RequestRedraw requests the host to render a new frame.
	// In on-demand rendering mode, this triggers a single frame render.
	// In continuous mode, this is a no-op.
	RequestRedraw()
}

// NullWindowProvider implements WindowProvider with configurable defaults.
// Used for testing and headless operation.
//
// When SF is zero (the default), ScaleFactor returns 1.0.
//
// Example:
//
//	wp := gpucontext.NullWindowProvider{W: 800, H: 600, SF: 2.0}
//	w, h := wp.Size()       // 800, 600
//	scale := wp.ScaleFactor() // 2.0
type NullWindowProvider struct {
	// W is the window width in physical pixels.
	W int

	// H is the window height in physical pixels.
	H int

	// SF is the DPI scale factor. When zero, ScaleFactor returns 1.0.
	SF float64
}

// Size returns the configured window dimensions.
func (n NullWindowProvider) Size() (int, int) { return n.W, n.H }

// ScaleFactor returns the configured scale factor, defaulting to 1.0 when zero.
func (n NullWindowProvider) ScaleFactor() float64 {
	if n.SF == 0 {
		return 1.0
	}
	return n.SF
}

// RequestRedraw does nothing.
func (n NullWindowProvider) RequestRedraw() {}

// Ensure NullWindowProvider implements WindowProvider.
var _ WindowProvider = NullWindowProvider{}
