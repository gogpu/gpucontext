// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: BSD-3-Clause

package gpucontext

// WebGPU Type Definitions for Cross-Package Sharing
//
// This file defines interfaces and types that are implemented by wgpu/hal
// and used by consumers like gg, gogpu, and born-ml.
//
// IMPORTANT: This package has ZERO dependencies to avoid circular imports.
// All types here are interfaces that backends implement.

// TextureFormat specifies the format of texture data.
// Values match the WebGPU specification.
type TextureFormat uint32

// Common texture formats matching WebGPU spec.
const (
	TextureFormatUndefined      TextureFormat = 0
	TextureFormatRGBA8Unorm     TextureFormat = 1
	TextureFormatRGBA8UnormSrgb TextureFormat = 2
	TextureFormatBGRA8Unorm     TextureFormat = 3
	TextureFormatBGRA8UnormSrgb TextureFormat = 4
	// Add more as needed - these are the most common
)

// Device represents a logical GPU device.
// Implemented by wgpu/hal.Device.
type Device interface {
	// CreateBuffer creates a GPU buffer.
	// CreateBuffer(descriptor BufferDescriptor) Buffer

	// CreateTexture creates a GPU texture.
	// CreateTexture(descriptor TextureDescriptor) Texture

	// CreateShaderModule creates a shader module from source.
	// CreateShaderModule(descriptor ShaderModuleDescriptor) ShaderModule

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
	// WriteTexture(destination ImageCopyTexture, data []byte, layout TextureDataLayout, size Extent3D)
}

// Adapter represents a physical GPU.
// Implemented by wgpu/hal.Adapter.
type Adapter interface {
	// RequestDevice requests a logical device from this adapter.
	// RequestDevice(descriptor DeviceDescriptor) (Device, Queue, error)

	// GetInfo returns information about this adapter.
	// GetInfo() AdapterInfo

	// Features returns the features supported by this adapter.
	// Features() Features

	// Limits returns the limits of this adapter.
	// Limits() Limits
}

// Surface represents a rendering surface (window).
// Implemented by wgpu/hal.Surface.
type Surface interface {
	// Configure configures the surface for rendering.
	// Configure(device Device, config SurfaceConfiguration)

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
	// RequestAdapter(options RequestAdapterOptions) (Adapter, error)

	// EnumerateAdapters returns all available adapters.
	// EnumerateAdapters() []Adapter
}

// OpenDevice bundles a device and queue together.
// This is a convenience type for initialization.
type OpenDevice struct {
	Device Device
	Queue  Queue
}
