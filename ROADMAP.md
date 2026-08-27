# Roadmap

## Vision

`gpucontext` is the shared foundation for the [gogpu](https://github.com/gogpu) ecosystem, providing interfaces and utilities for GPU resource sharing without circular dependencies.

## Current: v0.29.0

- `DeviceProvider.Features()` — query device capabilities (FeatureRayQuery, etc.)
- SurfaceCompositor interface (ADR-067)
- MouseButton/Modifiers stringers
- Key enum redesign (grouped explicit bases, binary-stable)
- DamageSource registration (ADR-065)
- SubpixelLayout on PlatformProvider (ADR-024)
- AdapterInfo on DeviceProvider (ADR-020)
- WindowChrome.SetFullscreen / IsFullscreen (ADR-018)
- CursorMode (Locked/Confined/Normal) for mouse grab
- WindowProvider, PlatformProvider, DeviceProvider, TextureUpdater, EventSource

## Released

### v0.29.0 (2026-08-27)
- `DeviceProvider.Features()` — gputypes.Features bitfield for capability queries

### v0.28.0 (2026-08-13)
- MouseButton.String(), Modifiers.String()

### v0.27.0 (2026-08-10)
- SurfaceCompositor interface (ADR-067)

### v0.26.0 (2026-08-10)
- Key enum redesign (grouped explicit bases)

### v0.24.0 (2026-07-21)
- DamageSource registration (ADR-065)

### v0.18.0 (2026-05-09)
- SubpixelLayout on PlatformProvider

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
