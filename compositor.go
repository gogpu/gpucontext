// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

package gpucontext

import "image"

// SurfaceCompositor is the interface for compositor-level decisions about
// surface rendering. The compositor (gogpu) owns the surface lifecycle
// and makes decisions about LoadOp, damage scissoring, and swapchain
// presentation. Content renderers (gg, g3d) use this interface to query
// compositor state when recording render passes.
//
// This cleanly separates compositor concerns (gogpu) from content rendering
// (gg). Before this interface, gg made compositor decisions internally
// (LoadOp, damage scissoring, MSAA overlay compositing). Now gogpu provides
// these decisions through the compositor, and gg records draws accordingly.
//
// ADR-067: surface composition architecture.
type SurfaceCompositor interface {
	// ShouldPreserveContent reports whether the current render pass should
	// use LoadOpLoad to preserve existing surface content. True when:
	// - An earlier renderer already drew to the surface this frame
	// - External content (g3d) has been rendered to the surface
	// Content renderers use this to decide between LoadOpClear and LoadOpLoad.
	ShouldPreserveContent() bool

	// DamageRects returns the damage rectangles for the current frame.
	// When non-empty, the compositor has determined that only these regions
	// need re-rendering. Content renderers should apply scissor rects
	// accordingly. Empty means full-surface render.
	DamageRects() []image.Rectangle

	// MarkContentRendered signals that a content renderer has drawn to the
	// surface in the current frame. Subsequent renderers will see
	// ShouldPreserveContent() == true.
	MarkContentRendered()

	// CompositeMSAAOverlay requests the compositor to alpha-blend an MSAA
	// overlay resolve texture onto the surface. The compositeView is the
	// single-sample texture containing the resolved MSAA overlay.
	// The compositor uses its own blit pipeline (independent of gg's pipelines)
	// to prevent bind group lifetime conflicts.
	//
	// encoder is the shared frame encoder. view is the swapchain target.
	// compositeView is the resolved MSAA overlay source.
	// w, h are surface dimensions.
	CompositeMSAAOverlay(encoder CommandEncoder, view TextureView, compositeView TextureView, w, h uint32) error
}
