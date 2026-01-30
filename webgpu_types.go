// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

package gpucontext

// WebGPU Interface Definitions for Cross-Package Sharing
//
// This file defines interfaces that are implemented by wgpu/hal
// and used by consumers like gg, gogpu, and born-ml.
//
// Types (TextureFormat, BufferUsage, etc.) are in gputypes package.
// Interfaces (Device, Queue, etc.) are defined here as behavioral contracts.
//
// Users should import both packages:
//
//	import (
//	    "github.com/gogpu/gpucontext"
//	    "github.com/gogpu/gputypes"
//	)

// Device represents a logical GPU device.
// Implemented by wgpu/hal.Device.
type Device interface {
	// Poll processes pending operations.
	Poll(wait bool)

	// Destroy releases the device resources.
	Destroy()
}

// Queue represents a GPU command queue.
// Implemented by wgpu/hal.Queue.
type Queue interface {
	// Submit submits command buffers for execution.
	// Submit(commandBuffers []CommandBuffer)

	// WriteBuffer writes data to a buffer.
	// WriteBuffer(buffer Buffer, offset uint64, data []byte)

	// WriteTexture writes data to a texture.
	// WriteTexture(destination gputypes.ImageCopyTexture, data []byte, layout gputypes.TextureDataLayout, size gputypes.Extent3D)
}

// Adapter represents a physical GPU.
// Implemented by wgpu/hal.Adapter.
type Adapter interface {
	// RequestDevice requests a logical device from this adapter.
	// RequestDevice(descriptor gputypes.DeviceDescriptor) (Device, Queue, error)

	// GetInfo returns information about this adapter.
	// GetInfo() gputypes.AdapterInfo

	// Features returns the features supported by this adapter.
	// Features() gputypes.Features

	// Limits returns the limits of this adapter.
	// Limits() gputypes.Limits
}

// Surface represents a rendering surface (window).
// Implemented by wgpu/hal.Surface.
type Surface interface {
	// Configure configures the surface for rendering.
	// Configure(device Device, config gputypes.SurfaceConfiguration)

	// GetCurrentTexture gets the current texture for rendering.
	// GetCurrentTexture() (SurfaceTexture, error)

	// Present presents the current frame.
	// Present()
}

// Instance is the entry point for GPU operations.
// Implemented by wgpu/hal.Instance.
type Instance interface {
	// CreateSurface creates a surface for a window.
	// CreateSurface(descriptor SurfaceDescriptor) (Surface, error)

	// RequestAdapter requests a GPU adapter.
	// RequestAdapter(options gputypes.RequestAdapterOptions) (Adapter, error)

	// EnumerateAdapters returns all available adapters.
	// EnumerateAdapters() []Adapter
}

// OpenDevice bundles a device and queue together.
// This is a convenience type for initialization.
type OpenDevice struct {
	Device Device
	Queue  Queue
}
