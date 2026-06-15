// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

package gpucontext

// WebGPU Type Token Interfaces for Cross-Package Sharing
//
// This file defines type token interfaces for GPU objects (Device, Queue, etc.)
// that enable type-safe dependency injection between packages without coupling
// them to a specific GPU implementation.
//
// Concrete types (e.g., *wgpu.Device) satisfy these empty interfaces implicitly.
// Consumers type-assert to the concrete type when they need the full API.
//
// Note: TextureView and CommandEncoder are opaque handle structs (not interfaces),
// defined in handle.go. They use unsafe.Pointer for compile-time type safety
// without boxing allocations (ADR-018).
//
// Types (TextureFormat, BufferUsage, etc.) are in the gputypes package.
//
// Usage:
//
//	import (
//	    "github.com/gogpu/gpucontext"
//	    "github.com/gogpu/gputypes"
//	)

// Device is a type-safe token for a logical GPU device.
//
// The unexported sentinel method prevents nil or arbitrary values from
// satisfying this interface — only concrete GPU implementations (e.g.,
// *wgpu.Device) can implement it. Consumers type-assert to the concrete
// type when they need the full API:
//
//	dev := provider.Device()
//	wgpuDev, ok := dev.(*wgpu.Device)
type Device interface {
	gpuDevice() // sentinel — only wgpu.Device implements this
}

// Queue is a type-safe token for a GPU command queue.
//
// See Device for the sentinel method pattern.
type Queue interface {
	gpuQueue() // sentinel — only wgpu.Queue implements this
}

// Adapter is a type-safe token for a physical GPU adapter.
type Adapter interface {
	gpuAdapter() // sentinel — only wgpu.Adapter implements this
}

// Surface is a type-safe token for a rendering surface (window).
type Surface interface {
	gpuSurface() // sentinel — only wgpu.Surface implements this
}

// TextureView is now a type-safe opaque handle struct defined in handle.go.
// See NewTextureView, TextureView.Pointer, TextureView.IsNil.

// Instance is a type-safe token for the GPU instance entry point.
type Instance interface {
	gpuInstance() // sentinel — only wgpu.Instance implements this
}

// OpenDevice bundles a device and queue together.
// This is a convenience type for initialization.
type OpenDevice struct {
	Device Device
	Queue  Queue
}
