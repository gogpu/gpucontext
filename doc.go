// Package gpucontext provides shared GPU infrastructure for the gogpu ecosystem.
//
// This package defines interfaces and utilities used across multiple gogpu
// projects to enable GPU resource sharing without circular dependencies:
//
//   - DeviceProvider: Interface for providing GPU device and queue
//   - EventSource: Interface for window/input events (keyboard, mouse)
//   - PointerEventSource: Interface for unified pointer events (W3C Level 3, mouse+touch+pen)
//   - WindowProvider: Interface for window geometry, DPI, and redraw requests
//   - PlatformProvider: Interface for clipboard, cursor, dark mode, accessibility
//   - ScrollEventSource: Interface for detailed scroll events
//   - WindowChrome: Interface for custom window chrome (frameless windows)
//   - Texture: Minimal interface for GPU textures
//   - TextureDrawer: Interface for drawing textures (2D rendering)
//   - TextureCreator: Interface for creating textures from pixel data
//
// # Consumers
//
//   - gogpu/gogpu: Implements DeviceProvider via App/Renderer
//   - gogpu/gg: Uses DeviceProvider for GPU-accelerated 2D rendering
//   - born-ml/born: Implements and uses for GPU compute
//
// # Design Principles
//
// This package follows the wgpu ecosystem pattern where shared types
// are separated from implementation (cf. wgpu-types in Rust).
//
// The key insight is that GPU context (device + queue + related state)
// is a universal concept across Vulkan, CUDA, OpenGL, and WebGPU.
// By defining a minimal interface here, different packages can share
// GPU resources without depending on each other.
//
// # GPU Handle Types
//
// Device, Queue, Adapter, Surface, Instance are opaque struct tokens
// wrapping unsafe.Pointer — same pattern as TextureView and CommandEncoder.
// 8 bytes, value type, zero allocations, GC-safe (reflect.Value precedent).
//
// # Example Usage
//
//	// In gogpu/gogpu — wraps *wgpu.Device into opaque handle
//	func (a *adapter) Device() gpucontext.Device {
//	    return gpucontext.NewDevice(unsafe.Pointer(a.renderer.device))
//	}
//
//	// In gogpu/gg — extracts concrete type from handle
//	func initGPU(provider gpucontext.DeviceProvider) {
//	    dev := provider.Device()
//	    wgpuDev := (*wgpu.Device)(dev.Pointer())
//	}
//
// Reference: https://github.com/gogpu/gpucontext
package gpucontext

const stringNone = "None"
