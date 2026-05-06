// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

package gpucontext

// AdapterType describes the physical GPU adapter type.
// Used by rendering libraries (gg) for render mode auto-selection.
type AdapterType int

const (
	// AdapterTypeDiscrete is a dedicated GPU (NVIDIA, AMD discrete).
	AdapterTypeDiscrete AdapterType = iota
	// AdapterTypeIntegrated is an integrated GPU (Intel Iris, AMD Vega iGPU).
	AdapterTypeIntegrated
	// AdapterTypeSoftware is a CPU-based software renderer (llvmpipe, SwiftShader, gogpu/wgpu software HAL).
	AdapterTypeSoftware
	// AdapterTypeUnknown means the adapter type could not be determined.
	AdapterTypeUnknown
)

// String returns the adapter type name.
func (t AdapterType) String() string {
	switch t {
	case AdapterTypeDiscrete:
		return "Discrete"
	case AdapterTypeIntegrated:
		return "Integrated"
	case AdapterTypeSoftware:
		return "Software"
	default:
		return "Unknown"
	}
}

// AdapterInfo describes the physical GPU adapter.
// Provided by DeviceProvider for render mode decisions.
type AdapterInfo struct {
	// Name is the adapter name (e.g., "NVIDIA GeForce RTX 4090", "Software Renderer").
	Name string
	// Type is the adapter type (Discrete, Integrated, Software, Unknown).
	Type AdapterType
}
