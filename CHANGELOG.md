# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.7.0] - 2026-02-05

### Added

- **TextureUpdater interface** for updating existing texture pixel data ([gg#79](https://github.com/gogpu/gg/issues/79))
  - `UpdateData(data []byte) error` — upload new pixel data to existing texture
  - Enables proper error handling for dynamic content (canvas rendering, video frames)
  - Implemented by `gogpu.Texture`

## [0.6.0] - 2026-01-31

### Added

- **Gesture Events** for multi-touch gesture recognition ([#6](https://github.com/gogpu/gpucontext/pull/6))
  - `GestureEvent` — Vello-style per-frame gesture deltas (zoom, rotation, translation)
  - `GestureEventSource` — interface for registering gesture callbacks
  - `NullGestureEventSource` — no-op implementation

## [0.5.0] - 2026-01-31

### Added

- **W3C Pointer Events Level 3** for unified pointer input
  - `PointerEvent` — unified mouse, touch, pen input with full W3C compliance
  - `PointerEventType` — Down, Up, Move, Enter, Leave, Cancel
  - `PointerType` — Mouse, Touch, Pen
  - `Button` — Left, Middle, Right, X1, X2, Eraser
  - `Buttons` — bitmask for tracking multiple pressed buttons
  - `PointerEventSource` — interface for registering pointer callbacks
  - `NullPointerEventSource` — no-op implementation

- **Scroll Events** for mouse wheel and trackpad
  - `ScrollEvent` — horizontal/vertical scroll with delta modes
  - `ScrollDeltaMode` — Pixel, Line, Page modes
  - `ScrollEventSource` — interface for registering scroll callbacks
  - `NullScrollEventSource` — no-op implementation

- **CI/CD Infrastructure**
  - GitHub Actions workflow (build, test, lint on Linux/macOS/Windows)
  - golangci-lint v2 configuration

### Changed

- **TouchCancelled → TouchCanceled** — US English spelling (misspell linter)
- Removed unused `DeviceHandle` alias

[0.7.0]: https://github.com/gogpu/gpucontext/releases/tag/v0.7.0
[0.6.0]: https://github.com/gogpu/gpucontext/releases/tag/v0.6.0
[0.5.0]: https://github.com/gogpu/gpucontext/releases/tag/v0.5.0

## [0.4.0] - 2026-01-30

### Added

- **Texture interfaces** for GPU texture sharing across packages
  - `Texture` — minimal interface with Width/Height
  - `TextureDrawer` — interface for drawing textures (DrawTexture, DrawTextureEx)
  - `TextureCreator` — interface for creating textures from pixel data
  - `TextureDrawOptions` — options for advanced texture rendering (position, scale, alpha, flip)

- **Touch input support** for mobile and tablet applications
  - `TouchID` — unique identifier for touch points
  - `TouchPhase` — lifecycle stages (Began, Moved, Ended, Cancelled)
  - `TouchPoint` — single touch contact with position, optional pressure/radius
  - `TouchEvent` — complete touch event with Changed/All points, modifiers, timestamp
  - `TouchEventSource` — interface for registering touch callbacks
  - `NullTouchEventSource` — no-op implementation for non-touch platforms

### Notes

- Touch interfaces follow platform conventions (iOS, Android, W3C Touch Events)
- Texture interfaces enable gg↔gogpu integration without circular dependencies
- Both are **contracts only** — implementations in host applications

[0.4.0]: https://github.com/gogpu/gpucontext/releases/tag/v0.4.0

## [0.3.1] - 2026-01-29

### Changed

- **Update gputypes to v0.2.0** for webgpu.h spec-compliant enum values

[0.3.1]: https://github.com/gogpu/gpucontext/releases/tag/v0.3.1

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
