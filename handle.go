// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

package gpucontext

import "unsafe"

// TextureView is a type-safe opaque handle to a GPU texture view.
//
// Provides compile-time type safety: TextureView cannot be confused with
// CommandEncoder or other GPU resource types (unlike the previous interface{}
// token approach). Following the Vulkan/Ebitengine/Go Protobuf Opaque pattern.
//
// The ptr field is unexported, preventing construction outside this package.
// Use NewTextureView to create, Pointer to extract, IsNil to check.
// 8 bytes, value type, zero allocations.
//
// GC safety: unsafe.Pointer keeps the underlying object alive (Go spec §Safety).
type TextureView struct {
	ptr unsafe.Pointer
}

// NewTextureView creates a TextureView from an unsafe.Pointer to a concrete
// GPU texture view (e.g., *wgpu.TextureView). The caller must ensure the
// pointer remains valid for the lifetime of the returned handle.
func NewTextureView(ptr unsafe.Pointer) TextureView {
	return TextureView{ptr: ptr}
}

// Pointer returns the underlying unsafe.Pointer. Consumers type-convert to
// the concrete type: (*wgpu.TextureView)(tv.Pointer()).
func (tv TextureView) Pointer() unsafe.Pointer { return tv.ptr }

// IsNil reports whether the handle holds no resource (zero value).
func (tv TextureView) IsNil() bool { return tv.ptr == nil }

// CommandEncoder is a type-safe opaque handle to a GPU command encoder.
//
// Same pattern as TextureView. Used for the shared encoder pipeline (ADR-017).
// 8 bytes, value type, zero allocations.
type CommandEncoder struct {
	ptr unsafe.Pointer
}

// NewCommandEncoder creates a CommandEncoder from an unsafe.Pointer to a
// concrete GPU command encoder (e.g., *wgpu.CommandEncoder).
func NewCommandEncoder(ptr unsafe.Pointer) CommandEncoder {
	return CommandEncoder{ptr: ptr}
}

// Pointer returns the underlying unsafe.Pointer. Consumers type-convert to
// the concrete type: (*wgpu.CommandEncoder)(ce.Pointer()).
func (ce CommandEncoder) Pointer() unsafe.Pointer { return ce.ptr }

// IsNil reports whether the handle holds no resource (zero value).
func (ce CommandEncoder) IsNil() bool { return ce.ptr == nil }
