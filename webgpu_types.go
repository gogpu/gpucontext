// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

package gpucontext

import "unsafe"

// WebGPU Opaque Handle Types for Cross-Package Sharing
//
// GPU resource handles (Device, Queue, Adapter, Surface, Instance) are
// struct tokens wrapping unsafe.Pointer — the same pattern as TextureView
// and CommandEncoder in handle.go.
//
// This design provides:
//   - Cross-package type safety: Device, Queue, Adapter are distinct types
//   - GC safety: unsafe.Pointer in struct fields is traced by GC (Go spec §Safety)
//   - Zero allocations: 8-byte value type, no interface boxing
//   - Compile-time protection: ptr field is unexported, only NewDevice() etc. can construct
//
// Precedent: reflect.Value uses the identical pattern (struct with unsafe.Pointer field).
//
// Consumers extract the concrete type via Pointer():
//
//	dev := provider.Device()
//	wgpuDev := (*wgpu.Device)(dev.Pointer())
//
// Or via helper in wgpu package:
//
//	wgpuDev := wgpu.DeviceFromHandle(dev)
//
// Types (TextureFormat, BufferUsage, etc.) are in the gputypes package.

// Device is a type-safe opaque handle to a logical GPU device.
// 8 bytes, value type, zero allocations. GC-safe.
type Device struct{ ptr unsafe.Pointer }

// NewDevice creates a Device handle from an unsafe.Pointer to a concrete
// GPU device (e.g., *wgpu.Device).
func NewDevice(ptr unsafe.Pointer) Device { return Device{ptr: ptr} }

// Pointer returns the underlying unsafe.Pointer.
func (d Device) Pointer() unsafe.Pointer { return d.ptr }

// IsNil reports whether the handle holds no resource (zero value).
func (d Device) IsNil() bool { return d.ptr == nil }

// Queue is a type-safe opaque handle to a GPU command queue.
// 8 bytes, value type, zero allocations. GC-safe.
type Queue struct{ ptr unsafe.Pointer }

// NewQueue creates a Queue handle.
func NewQueue(ptr unsafe.Pointer) Queue { return Queue{ptr: ptr} }

// Pointer returns the underlying unsafe.Pointer.
func (q Queue) Pointer() unsafe.Pointer { return q.ptr }

// IsNil reports whether the handle holds no resource.
func (q Queue) IsNil() bool { return q.ptr == nil }

// Adapter is a type-safe opaque handle to a physical GPU adapter.
// 8 bytes, value type, zero allocations. GC-safe.
type Adapter struct{ ptr unsafe.Pointer }

// NewAdapter creates an Adapter handle.
func NewAdapter(ptr unsafe.Pointer) Adapter { return Adapter{ptr: ptr} }

// Pointer returns the underlying unsafe.Pointer.
func (a Adapter) Pointer() unsafe.Pointer { return a.ptr }

// IsNil reports whether the handle holds no resource.
func (a Adapter) IsNil() bool { return a.ptr == nil }

// Surface is a type-safe opaque handle to a rendering surface (window).
// 8 bytes, value type, zero allocations. GC-safe.
type Surface struct{ ptr unsafe.Pointer }

// NewSurface creates a Surface handle.
func NewSurface(ptr unsafe.Pointer) Surface { return Surface{ptr: ptr} }

// Pointer returns the underlying unsafe.Pointer.
func (s Surface) Pointer() unsafe.Pointer { return s.ptr }

// IsNil reports whether the handle holds no resource.
func (s Surface) IsNil() bool { return s.ptr == nil }

// Instance is a type-safe opaque handle to the GPU instance entry point.
// 8 bytes, value type, zero allocations. GC-safe.
type Instance struct{ ptr unsafe.Pointer }

// NewInstance creates an Instance handle.
func NewInstance(ptr unsafe.Pointer) Instance { return Instance{ptr: ptr} }

// Pointer returns the underlying unsafe.Pointer.
func (i Instance) Pointer() unsafe.Pointer { return i.ptr }

// IsNil reports whether the handle holds no resource.
func (i Instance) IsNil() bool { return i.ptr == nil }

// OpenDevice bundles a device and queue together.
type OpenDevice struct {
	Device Device
	Queue  Queue
}
