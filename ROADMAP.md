# Roadmap

## Vision

`gpucontext` is the shared foundation for the [gogpu](https://github.com/gogpu) ecosystem, providing interfaces and utilities for GPU resource sharing without circular dependencies.

## Current: v0.14.0

- CursorMode (Locked/Confined/Normal) for mouse grab / pointer lock
- PointerEvent.DeltaX/DeltaY for relative mouse movement
- WindowChrome interface for frameless window support
- WindowProvider, PlatformProvider, DeviceProvider, TextureUpdater, EventSource
- Texture interfaces (multi-touch, pressure, radius)
- W3C Pointer Events Level 3, Scroll events, Gesture events
- IME support for CJK input

## Released

### v0.12.0 (2026-04-09)
- CursorMode + PointerEvent DeltaX/DeltaY for mouse grab (gogpu#173)

### v0.11.0 (2026-03-20)
- WindowChrome interface for frameless window support

### v0.10.0 (2026-03-11)
- Typed Device/Queue interfaces, removed HalProvider

### v0.9.0 (2026-02-27)
- WindowProvider DPI/HiDPI support

### v0.8.0 (2026-02-15)
- WindowProvider, PlatformProvider, CursorShape, NullProviders

### v0.7.0 (2026-02-05)
- TextureUpdater interface for dynamic texture content

### v0.6.0 (2026-01-31)
- Gesture events (GestureEvent, GestureEventSource)

### v0.5.0 (2026-01-31)
- W3C Pointer Events Level 3, scroll events, CI/CD

### v0.4.0 (2026-01-30)
- Texture interfaces, touch input support

### v0.3.1 (2026-01-29)
- Update gputypes to v0.2.0

### v0.3.0 (2026-01-29)
- Import gputypes for unified WebGPU types

### v0.2.0 (2026-01-27)
- IME support for CJK input

### v0.1.1 (2026-01-27)
- Initial release with DeviceProvider, EventSource, Registry

## Future Considerations

### v1.0.0 — API Freeze
- Stable API guarantee
- Full WebGPU spec coverage
- Comprehensive documentation

## Non-Goals

- This package will **never** contain implementations
- This package will **never** have external dependencies (beyond gputypes)
- This package focuses on **interfaces**, not concrete types
