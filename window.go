// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

package gpucontext

// WindowProvider provides window geometry and DPI information.
//
// This interface enables UI frameworks (like gogpu/ui) and rendering
// libraries (like gg/ggcanvas) to query window dimensions and scale factor
// for HiDPI-aware layout and rendering.
//
// Size() returns logical points (DIP — density-independent pixels).
// To get physical pixel dimensions, multiply by ScaleFactor():
//
//	physW := int(float64(w) * scale)
//
// Implementations:
//   - gogpu.App implements WindowProvider via platform window
//   - gogpu.gpuContextAdapter implements WindowProvider (returned by GPUContextProvider)
//   - NullWindowProvider provides configurable defaults for testing
//
// Example usage:
//
//	func (ui *UI) Layout(wp gpucontext.WindowProvider) {
//	    w, h := wp.Size()           // logical points
//	    scale := wp.ScaleFactor()   // 2.0 on Retina
//	    ui.root.Layout(w, h, scale)
//	}
type WindowProvider interface {
	// Size returns the current window client area dimensions in logical points (DIP).
	// On HiDPI displays, multiply by ScaleFactor() to get physical pixel dimensions.
	Size() (width, height int)

	// ScaleFactor returns the DPI scale factor.
	// 1.0 = standard (96 DPI on Windows, 72 on macOS), 2.0 = Retina/HiDPI.
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
//	w, h := wp.Size()       // 800, 600 (logical points)
//	scale := wp.ScaleFactor() // 2.0
type NullWindowProvider struct {
	// W is the window width in logical points (DIP).
	W int

	// H is the window height in logical points (DIP).
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
