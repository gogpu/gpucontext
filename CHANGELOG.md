# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.3.0] - 2026-01-29

### Changed

- **Import gputypes for unified WebGPU types**
  - DeviceProvider.SurfaceFormat() now returns `gputypes.TextureFormat`
  - Removed local type re-exports in favor of gputypes
  - Single source of truth for WebGPU types across ecosystem

### Added

- CODE_OF_CONDUCT.md
- SECURITY.md

[0.3.0]: https://github.com/gogpu/gpucontext/releases/tag/v0.3.0

## [0.2.0] - 2026-01-27

### Added

- **IME Support** for CJK input (Chinese, Japanese, Korean)
  - `IMEState` struct with composition state tracking
  - `IMEController` interface for positioning IME window
  - Extended `EventSource` with `OnIMECompositionStart`, `OnIMECompositionUpdate`, `OnIMECompositionEnd`
  - Updated `NullEventSource` with no-op IME implementations

### Notes

- IME interfaces are **contracts only** — platform integration happens in host applications (gogpu)
- Required for enterprise UI frameworks supporting international users

[0.2.0]: https://github.com/gogpu/gpucontext/releases/tag/v0.2.0

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
