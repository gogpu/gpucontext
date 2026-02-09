// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

package gpucontext

// HalProvider is an optional interface for DeviceProviders that can expose
// low-level HAL types directly. This enables GPU accelerators to share
// devices without creating their own.
//
// The returned any values are hal.Device and hal.Queue from wgpu/hal.
// Consumers type-assert to the concrete hal types they need.
//
// Example usage:
//
//	if hp, ok := provider.(gpucontext.HalProvider); ok {
//	    device := hp.HalDevice().(hal.Device)
//	    queue := hp.HalQueue().(hal.Queue)
//	}
type HalProvider interface {
	// HalDevice returns the underlying HAL device for direct GPU access.
	// Returns nil if HAL access is not available.
	HalDevice() any

	// HalQueue returns the underlying HAL queue for direct GPU access.
	// Returns nil if HAL access is not available.
	HalQueue() any
}
