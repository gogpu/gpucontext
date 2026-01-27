package gpucontext

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
//	func NewGPUCanvas(provider gpucontext.DeviceProvider) *Canvas {
//	    return &Canvas{
//	        device: provider.Device(),
//	        queue:  provider.Queue(),
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
	// May return TextureFormatUndefined if no surface is attached (headless mode).
	// This is useful for creating render targets that match the surface format.
	SurfaceFormat() TextureFormat

	// Adapter returns the WebGPU adapter (optional, may be nil).
	// The adapter provides information about the GPU capabilities.
	// Some implementations may not expose the adapter.
	Adapter() Adapter
}

// DeviceHandle is an alias for DeviceProvider for backward compatibility.
// Deprecated: Use DeviceProvider instead.
type DeviceHandle = DeviceProvider
