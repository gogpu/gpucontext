# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.1] - 2026-01-27

### Added

- **DeviceProvider** interface for GPU device/queue injection
  - `Device()` returns WebGPU device
  - `Queue()` returns command queue
  - `Adapter()` returns GPU adapter
  - `SurfaceFormat()` returns preferred texture format

- **EventSource** interface for input events
  - Keyboard: `OnKeyPress`, `OnKeyRelease`, `OnTextInput`
  - Mouse: `OnMouseMove`, `OnMousePress`, `OnMouseRelease`, `OnScroll`
  - Window: `OnResize`, `OnFocus`
  - `Key`, `Modifiers`, `MouseButton` types
  - `NullEventSource` no-op implementation

- **Registry[T]** generic backend registry
  - Thread-safe registration with `sync.RWMutex`
  - Priority-based selection via `Best()`
  - `Register`, `Unregister`, `Get`, `Has`, `Available`, `Count`

- **WebGPU Types** (zero dependencies)
  - `Device`, `Queue`, `Adapter`, `Surface`, `Instance` interfaces
  - `TextureFormat` enum with common formats
  - `OpenDevice` convenience struct

### Notes

- This package has **zero external dependencies** by design
- All interfaces are minimal to allow diverse implementations
- Part of the [gogpu](https://github.com/gogpu) ecosystem

[0.1.1]: https://github.com/gogpu/gpucontext/releases/tag/v0.1.1
