# Roadmap

## Vision

`gpucontext` is the shared foundation for the [gogpu](https://github.com/gogpu) ecosystem, providing interfaces and utilities for GPU resource sharing without circular dependencies.

## Current: v0.2.0

**Status:** Released

- DeviceProvider interface
- EventSource interface (with IME support)
- Registry[T] generic
- WebGPU type definitions
- IME support for CJK input

## Previous: v0.1.1

- Initial release with DeviceProvider, EventSource, Registry

## Planned: v0.3.0

**Focus:** Extended capabilities

- [ ] `Capabilities` interface for feature queries
- [ ] `ResourceLimits` struct for GPU limits
- [ ] `ShaderFormat` enum (WGSL, SPIR-V, GLSL)
- [ ] `BufferUsage`, `TextureUsage` flags

## Future Considerations

### v0.4.0 — Compute Support
- ComputeProvider interface
- WorkgroupLimits
- Storage buffer types

### v1.0.0 — API Freeze
- Stable API guarantee
- Full WebGPU spec coverage
- Comprehensive documentation

## Non-Goals

- This package will **never** contain implementations
- This package will **never** have external dependencies
- This package focuses on **interfaces**, not concrete types

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.
