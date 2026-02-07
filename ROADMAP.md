# Roadmap

## Vision

`gpucontext` is the shared foundation for the [gogpu](https://github.com/gogpu) ecosystem, providing interfaces and utilities for GPU resource sharing without circular dependencies.

## Current: v0.8.0

**Status:** In Development

- WindowProvider interface (window geometry, DPI, redraw requests)
- PlatformProvider interface (clipboard, cursor, dark mode, accessibility)
- CursorShape enum (12 standard cursor shapes)
- NullWindowProvider / NullPlatformProvider for testing

## Released

### v0.7.0 (2026-02-05)
- TextureUpdater interface for dynamic texture content

### v0.6.0 (2026-01-31)
- Gesture events (GestureEvent, GestureEventSource)

### v0.5.0 (2026-01-31)
- W3C Pointer Events Level 3
- Scroll events with delta modes
- CI/CD infrastructure

### v0.4.0 (2026-01-30)
- Texture interfaces (Texture, TextureDrawer, TextureCreator)
- Touch input support (multi-touch, pressure, radius)

### v0.3.1 (2026-01-29)
- Update gputypes to v0.2.0

### v0.3.0 (2026-01-29)
- Import gputypes for unified WebGPU types

### v0.2.0 (2026-01-27)
- IME support for CJK input

### v0.1.1 (2026-01-27)
- Initial release with DeviceProvider, EventSource, Registry

## Future Considerations

### v0.9.0 — Extended Capabilities
- Capabilities interface for feature queries
- ResourceLimits for GPU limits
- ShaderFormat enum (WGSL, SPIR-V, GLSL)

### v1.0.0 — API Freeze
- Stable API guarantee
- Full WebGPU spec coverage
- Comprehensive documentation

## Non-Goals

- This package will **never** contain implementations
- This package will **never** have external dependencies (beyond gputypes)
- This package focuses on **interfaces**, not concrete types

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.
