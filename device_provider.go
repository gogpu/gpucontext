// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

package gpucontext

import "github.com/gogpu/gputypes"

// DeviceProvider provides access to GPU device, queue, and related resources.
// This interface enables dependency injection of GPU capabilities between
// packages without circular dependencies.
//
// Implementations:
//   - gogpu.App implements DeviceProvider via renderer
//   - born.Session implements DeviceProvider for ML compute
//
// Example usage in gg:
//
//	import (
//	    "github.com/gogpu/gpucontext"
//	    "github.com/gogpu/gputypes"
//	)
//
//	func NewGPUCanvas(provider gpucontext.DeviceProvider) *Canvas {
//	    format := provider.SurfaceFormat() // returns gputypes.TextureFormat
//	    return &Canvas{
//	        device: provider.Device(),
//	        queue:  provider.Queue(),
//	        format: format,
//	    }
//	}
type DeviceProvider interface {
	// Device returns the WebGPU device handle.
	// The device is used for creating GPU resources (buffers, textures, pipelines).
	Device() Device

	// Queue returns the WebGPU command queue.
	// The queue is used for submitting command buffers to the GPU.
	Queue() Queue

	// SurfaceFormat returns the preferred texture format for the surface.
	// May return gputypes.TextureFormatUndefined if no surface is attached (headless mode).
	// This is useful for creating render targets that match the surface format.
	SurfaceFormat() gputypes.TextureFormat

	// Adapter returns the WebGPU adapter (optional, may be nil).
	// The adapter provides information about the GPU capabilities.
	// Some implementations may not expose the adapter.
	Adapter() Adapter
}