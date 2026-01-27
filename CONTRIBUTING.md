# Contributing to gpucontext

Thank you for your interest in contributing to gpucontext!

## Guidelines

### Zero Dependencies Rule

**This package must have zero external dependencies.** This is a hard requirement, not a preference.

Why:
- Foundation packages cannot create circular dependencies
- Maximum compatibility across Go versions
- Minimal binary size impact

### Interface Design

When proposing new interfaces:

1. **Minimal** — Include only essential methods
2. **Orthogonal** — Interfaces should be composable
3. **Stable** — Changes break downstream packages

Example of good interface:
```go
type Device interface {
    Poll(wait bool)
    Destroy()
}
```

Example of bad interface (too specific):
```go
type Device interface {
    CreateBuffer(size int, usage BufferUsage) Buffer
    CreateTexture(desc TextureDescriptor) Texture
    CreateShaderModule(code string) ShaderModule
    // ... 50 more methods
}
```

### Code Style

- Follow standard Go conventions
- Run `go fmt` before committing
- Run `golangci-lint run` for linting
- Add tests for new functionality

### Pull Request Process

1. Fork the repository
2. Create a feature branch: `git checkout -b feat/my-feature`
3. Make changes with tests
4. Run: `go test ./... && golangci-lint run`
5. Submit PR with clear description

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add Capabilities interface
fix: correct TextureFormat values
docs: update README examples
test: add Registry edge cases
```

## Questions?

Open a [Discussion](https://github.com/gogpu/gpucontext/discussions) on GitHub.
