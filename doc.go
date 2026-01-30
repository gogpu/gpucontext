// Package gpucontext provides shared GPU infrastructure for the gogpu ecosystem.
//
// This package defines interfaces and utilities used across multiple gogpu
// projects to enable GPU resource sharing without circular dependencies:
//
//   - DeviceProvider: Interface for providing GPU device and queue
//   - EventSource: Interface for window/input events (keyboard, mouse)
//   - TouchEventSource: Interface for touch input (multi-touch, gestures)
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
// # Example Usage
//
//	// In gogpu/gogpu - implements DeviceProvider
//	func (app *App) Device() gpucontext.Device { return app.renderer.device }
//	func (app *App) Queue() gpucontext.Queue { return app.renderer.queue }
//
//	// In gogpu/gg - uses DeviceProvider
//	func NewGPUCanvas(provider gpucontext.DeviceProvider) *Canvas {
//	    return &Canvas{
//	        device: provider.Device(),
//	        queue:  provider.Queue(),
//	    }
//	}
//
// Reference: https://github.com/gogpu/gpucontext
package gpucontext
