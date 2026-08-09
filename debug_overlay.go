// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

package gpucontext

// Pluggable debug overlay system (ADR-066), inspired by GTK4's Inspector
// overlay architecture. Multiple debug overlays register with the compositor
// and are drawn in registration order after all content renderers, before
// present.
//
// GTK4 defines GtkInspectorOverlay with snapshot + queue_draw, supporting
// 8 concrete overlay types (updates, fps, layout, focus, a11y, etc.).
// Chromium uses a monolithic HeadsUpDisplayLayerImpl. We follow GTK4's
// pluggable model because it scales to N overlay types without monolith growth.
//
// Env vars control activation: GOGPU_DEBUG_DAMAGE=overlay, GOGPU_DEBUG_FPS=overlay,
// GOGPU_DEBUG_DIRTY=overlay. Each overlay is independent — any combination
// is valid simultaneously.

// DebugOverlay is a pluggable debug visualization layer.
//
// Multiple overlays register with the compositor via
// Context.RegisterDebugOverlay. They are drawn in registration order after
// all content renderers have finished, before present. Later overlays render
// on top of earlier ones (e.g., FPS counter on top of damage rects).
//
// This interface is frozen at 2 methods for v1.0+ stability. Future overlay
// capabilities go through DebugOverlayContext struct fields (non-breaking).
//
// Concrete overlays shipped with the ecosystem:
//   - Damage rects overlay (gogpu built-in + gg text-enhanced override)
//   - FPS counter overlay (gogpu built-in + gg text-enhanced override)
//   - Dirty widget overlay (ui, registers with gogpu)
//
// Example implementation:
//
//	type fpsOverlay struct {
//	    history [120]float64
//	    idx     int
//	}
//
//	func (f *fpsOverlay) Name() string { return "fps" }
//
//	func (f *fpsOverlay) Draw(ctx gpucontext.DebugOverlayContext) bool {
//	    // Render FPS counter using ctx.Encoder and ctx.SurfaceView
//	    return true // needs another frame for continuous update
//	}
//
// Registration:
//
//	dc.RegisterDebugOverlay(&fpsOverlay{})
type DebugOverlay interface {
	// Name identifies the overlay for logging, env var filtering, and
	// removal via RemoveDebugOverlay. Must be unique among registered
	// overlays. Examples: "damage", "fps", "dirty_widgets".
	Name() string

	// Draw renders the overlay on top of content.
	//
	// Returns true if the overlay needs another frame to complete its
	// visualization (e.g., fade animation in progress, FPS counter
	// updating). The compositor calls RequestRedraw when any overlay
	// returns true, creating a self-sustaining render loop that
	// automatically stops when all overlays return false.
	//
	// The overlay should use ctx.Encoder and ctx.SurfaceView to submit
	// GPU commands. All overlays use alpha blending (SrcAlpha /
	// OneMinusSrcAlpha) with LoadOp::Load to composite on top of content.
	Draw(ctx DebugOverlayContext) bool
}

// DebugOverlayContext provides GPU resources and frame metadata to a debug
// overlay's Draw method. The compositor (gogpu) populates this struct each
// frame before calling Draw on registered overlays.
//
// All fields are read-only from the overlay's perspective — the compositor
// owns the encoder and surface view lifetime.
type DebugOverlayContext struct {
	// SurfaceWidth is the surface width in physical pixels.
	SurfaceWidth uint32

	// SurfaceHeight is the surface height in physical pixels.
	SurfaceHeight uint32

	// Encoder is the GPU command encoder for the current frame.
	// The overlay uses this to begin render passes and submit draw commands.
	Encoder CommandEncoder

	// SurfaceView is the surface texture view for the current frame.
	// The overlay uses this as the render pass color attachment.
	SurfaceView TextureView

	// FrameNumber is the monotonic frame counter, useful for logging and
	// frame-based statistics (e.g., rolling average FPS calculation).
	FrameNumber uint64
}
