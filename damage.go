// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

package gpucontext

import (
	"image"
	"image/color"
)

// Damage reporting types and interfaces for the multi-renderer damage tracking
// system (ADR-065). When multiple renderers (gg, g3d, video, compose) share a
// single surface, each registers as a damage source and reports per-frame
// damage rectangles. The compositor (gogpu) unions all sources at present time.
//
// This design follows Chromium's cc DamageTracker pattern where all layers
// contribute damage and the compositor unions them — adapted for our explicit
// registration model where independent libraries register themselves rather
// than being discovered via a layer tree.

// DamageCategory classifies WHY damage occurred.
//
// Used for programmatic filtering, statistics, and overlay color
// differentiation (ADR-066). Chromium uses DamageReasonSet (bitfield of
// internal layer reasons). We use a category enum because our sources are
// independent libraries with diverse domain-specific reasons — a bitfield
// cannot cover all domains cleanly.
//
// Zero value is DamageCategoryContent, the most common category.
type DamageCategory uint8

const (
	// DamageCategoryContent means content changed (text, image, path, mesh).
	// This is the default category when no specific reason is provided.
	DamageCategoryContent DamageCategory = iota

	// DamageCategoryLayout means layout or position changed (widget moved, resized).
	DamageCategoryLayout

	// DamageCategoryAnimation means an animation tick occurred (spinner, transition, camera).
	DamageCategoryAnimation

	// DamageCategoryResize means the surface or window was resized.
	DamageCategoryResize

	// DamageCategoryFull means a full redraw is required (initial render, theme change).
	DamageCategoryFull

	// DamageCategoryExternal means damage from an external source (video frame, compose child).
	DamageCategoryExternal
)

// String returns the damage category name for debugging.
func (c DamageCategory) String() string {
	switch c {
	case DamageCategoryContent:
		return "Content"
	case DamageCategoryLayout:
		return "Layout"
	case DamageCategoryAnimation:
		return "Animation"
	case DamageCategoryResize:
		return "Resize"
	case DamageCategoryFull:
		return "Full"
	case DamageCategoryExternal:
		return "External"
	default:
		return "Unknown"
	}
}

// DamageReason combines a typed category (for filtering and optimization) with
// a human-readable detail string (for overlay labels and structured logging).
//
// Category is always meaningful — the zero value DamageCategoryContent is the
// most common case. Detail is optional — empty string means no additional
// context beyond the category.
//
// Examples:
//
//	DamageReason{}                                                        // content changed, no detail
//	DamageReason{Category: DamageCategoryAnimation, Detail: "camera rotation"} // animation with detail
//	DamageReason{Category: DamageCategoryFull, Detail: "theme change"}         // full redraw with detail
type DamageReason struct {
	// Category classifies the damage type for programmatic use.
	// Zero value (DamageCategoryContent) is the default for content changes.
	Category DamageCategory

	// Detail is a human-readable string for overlay labels and logging.
	// Empty string is valid and means no additional context beyond the category.
	// Examples: "camera rotation", "spinner tick", "theme change", "scene loaded".
	Detail string
}

// DamageReporter reports damage rectangles for a registered damage source.
//
// Each frame, the source calls ReportDamage or ReportDamageWithReason with the
// rectangles that changed. The compositor (gogpu) unions damage from all
// registered sources at present time to determine the final damage region
// sent to the Wayland compositor or other presentation backend.
//
// Calling ReportDamage with no rectangles signals full-surface damage for this
// source — the compositor will present the entire surface.
//
// Both methods reset after present — the source must report damage every frame
// that content changes.
//
// This interface is frozen at 2 methods for v1.0+ stability. Future metadata
// (per-rect labels, opaque regions, damage subtraction) goes through concrete
// type methods and DamageSourceSnapshot struct fields — both are non-breaking
// additions in Go.
//
// Implementations:
//   - gogpu.DamageSource implements DamageReporter (registered via Context.RegisterDamageSource)
//
// Example usage:
//
//	// gg (via ggcanvas) — reports partial damage:
//	ggSource.ReportDamage(dirtyRect1, dirtyRect2)
//
//	// g3d — reports full viewport with reason:
//	g3dSource.ReportDamageWithReason(
//	    gpucontext.DamageReason{Category: gpucontext.DamageCategoryAnimation, Detail: "camera rotation"},
//	    viewportRect,
//	)
//
//	// Full surface damage (no rects):
//	source.ReportDamage()
type DamageReporter interface {
	// ReportDamage reports damage rectangles for this frame.
	// No rects signals full-surface damage for this source.
	// Rects are in physical pixels (surface coordinates).
	// Damage is reset after present — report again next frame if content changes.
	ReportDamage(rects ...image.Rectangle)

	// ReportDamageWithReason reports damage with a typed reason for debug
	// overlay labels, structured logging, and future optimization heuristics.
	// No rects signals full-surface damage. Reason is reset after present.
	//
	// Use this method when the damage source can provide meaningful context
	// (e.g., "camera rotation", "spinner tick"). For common content changes
	// where no detail is needed, use ReportDamage instead.
	ReportDamageWithReason(reason DamageReason, rects ...image.Rectangle)
}

// DamageOverlayRenderer provides custom rendering for the damage debug overlay.
//
// The default overlay in gogpu renders flat-color quads (no text). Libraries
// with text rendering capability (gg) can register a custom renderer that adds
// anti-aliased borders, text labels per source, and richer visuals.
//
// Registration: gogpu.Context.SetDamageOverlayRenderer(renderer)
//
// The renderer receives a DamageOverlayInfo snapshot each frame when the
// damage overlay is active (GOGPU_DEBUG_DAMAGE=overlay). The renderer draws
// on top of content using the provided encoder and surface view.
//
// Implementations:
//   - gg/integration/ggcanvas provides a text-enhanced overlay renderer
//
// Example usage:
//
//	type ggDamageOverlay struct { /* ... */ }
//
//	func (o *ggDamageOverlay) RenderDamageOverlay(info gpucontext.DamageOverlayInfo) {
//	    for _, src := range info.Sources {
//	        // Draw colored rects with text labels using gg.Context
//	    }
//	}
type DamageOverlayRenderer interface {
	// RenderDamageOverlay draws the damage debug overlay for the current frame.
	// Called by the compositor after all content renderers have finished and
	// before present, when GOGPU_DEBUG_DAMAGE=overlay is active.
	//
	// The renderer should use info.Encoder and info.SurfaceView to submit GPU
	// commands. Fade tracking and flash state are the renderer's responsibility.
	RenderDamageOverlay(info DamageOverlayInfo)
}

// DamageOverlayInfo provides structured per-source damage data to the overlay
// renderer. The compositor (gogpu) constructs this from its internal damage
// sources — the renderer receives ready-to-use snapshots without needing to
// aggregate or group data.
//
// All fields are populated by gogpu before calling RenderDamageOverlay.
type DamageOverlayInfo struct {
	// Sources contains per-source snapshots for the current frame.
	// Each source corresponds to a registered damage source (gg, g3d, video, etc.).
	// Order matches registration order.
	Sources []DamageSourceSnapshot

	// FrameNumber is the monotonic frame counter, useful for logging
	// (e.g., "frame=1234") and frame-based statistics.
	FrameNumber uint64

	// SurfaceWidth is the surface width in physical pixels.
	SurfaceWidth uint32

	// SurfaceHeight is the surface height in physical pixels.
	SurfaceHeight uint32

	// Encoder is the GPU command encoder for the current frame.
	// The renderer uses this to begin render passes and submit draw commands.
	Encoder CommandEncoder

	// SurfaceView is the surface texture view for the current frame.
	// The renderer uses this as the render pass color attachment.
	SurfaceView TextureView
}

// DamageSourceSnapshot is a per-source damage snapshot for one frame.
//
// Name and Color are stable across frames — set once at RegisterDamageSource.
// Rects, Full, and Reason change every frame based on ReportDamage calls and
// are reset after present.
//
// Zero value: no damage reported (empty Rects, Full=false, zero Reason).
type DamageSourceSnapshot struct {
	// Name identifies the damage source, set at registration time.
	// Examples: "gg", "g3d", "video", "compose".
	Name string

	// Color is the overlay color assigned by the compositor from a fixed palette
	// at registration time. The color uses the stdlib image/color.RGBA type to
	// avoid adding dependencies to gpucontext.
	Color color.RGBA

	// Rects contains the damage rectangles reported this frame, in physical
	// pixels (surface coordinates). Empty when Full is true or no damage was
	// reported.
	Rects []image.Rectangle

	// Full is true when ReportDamage was called with no arguments, indicating
	// the entire surface is damaged for this source.
	Full bool

	// Reason is the damage reason from ReportDamageWithReason. Zero value
	// (DamageReason{}) means ReportDamage was used without a reason, which
	// is the common case for content changes.
	Reason DamageReason
}
