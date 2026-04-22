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
// Types (TextureFormat, BufferUsage, etc.) are in the gputypes package.
//
// Usage:
//
//	import (
//	    "github.com/gogpu/gpucontext"
//	    "github.com/gogpu/gputypes"
//	)

// Device is a type token for a logical GPU device.
//
// Concrete implementations (e.g., *wgpu.Device) satisfy this interface
// implicitly. Consumers that need the full device API should type-assert
// to the concrete type:
//
//	dev := provider.Device()
//	wgpuDev, ok := dev.(*wgpu.Device)
//	if ok {
//	    halDevice := wgpuDev.HalDevice()
//	}
//
// The interface is intentionally minimal to avoid coupling gpucontext
// to any specific GPU implementation.
type Device interface{}

// Queue is a type token for a GPU command queue.
//
// Concrete implementations (e.g., *wgpu.Queue) satisfy this interface
// implicitly. Consumers that need the full queue API should type-assert
// to the concrete type:
//
//	q := provider.Queue()
//	wgpuQueue, ok := q.(*wgpu.Queue)
type Queue interface{}

// Adapter is a type token for a physical GPU adapter.
// Consumers that need the full adapter API should type-assert
// to the concrete type (e.g., *wgpu.Adapter).
type Adapter interface{}

// Surface is a type token for a rendering surface (window).
// Consumers that need the full surface API should type-assert
// to the concrete type (e.g., *wgpu.Surface).
type Surface interface{}

// TextureView is a type token for a GPU texture view.
// Used as render target in render pass descriptors and for
// direct surface rendering. Consumers that need the full API
// should type-assert to the concrete type (e.g., *wgpu.TextureView).
type TextureView interface{}

// Instance is a type token for the GPU instance entry point.
// Consumers that need the full instance API should type-assert
// to the concrete type (e.g., *wgpu.Instance).
type Instance interface{}

// OpenDevice bundles a device and queue together.
// This is a convenience type for initialization.
type OpenDevice struct {
	Device Device
	Queue  Queue
}
