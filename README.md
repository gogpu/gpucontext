# gpucontext

Shared GPU infrastructure for the [gogpu](https://github.com/gogpu) ecosystem.

## Overview

`gpucontext` provides interfaces and utilities for sharing GPU resources across multiple packages without circular dependencies. It follows the pattern used by [wgpu-types](https://docs.rs/wgpu-types) in the Rust ecosystem.

## Installation

```bash
go get github.com/gogpu/gpucontext
```

**Requires:** Go 1.25+

## Features

- **DeviceProvider** — Interface for injecting GPU device and queue
- **EventSource** — Interface for input events (keyboard, mouse, window, IME)
- **IME Support** — Input Method Editor for CJK languages (Chinese, Japanese, Korean)
- **Registry[T]** — Generic registry with priority-based backend selection
- **WebGPU Types** — Device, Queue, Adapter, Surface interfaces (zero dependencies)

## Usage

### DeviceProvider Pattern

The `DeviceProvider` interface enables dependency injection of GPU capabilities:

```go
// In gogpu/gogpu - implements DeviceProvider
type App struct {
    device gpucontext.Device
    queue  gpucontext.Queue
}

func (app *App) Device() gpucontext.Device       { return app.device }
func (app *App) Queue() gpucontext.Queue         { return app.queue }
func (app *App) SurfaceFormat() gpucontext.TextureFormat { return app.format }
func (app *App) Adapter() gpucontext.Adapter     { return app.adapter }

// In gogpu/gg - uses DeviceProvider
func NewGPUCanvas(provider gpucontext.DeviceProvider) *Canvas {
    return &Canvas{
        device: provider.Device(),
        queue:  provider.Queue(),
    }
}
```

### EventSource (for UI frameworks)

`EventSource` enables UI frameworks to receive input events from host applications:

```go
// In gogpu/ui - uses EventSource
func (ui *UI) AttachEvents(source gpucontext.EventSource) {
    source.OnKeyPress(func(key gpucontext.Key, mods gpucontext.Modifiers) {
        ui.focused.HandleKeyDown(key, mods)
    })

    source.OnMousePress(func(button gpucontext.MouseButton, x, y float64) {
        widget := ui.hitTest(x, y)
        widget.HandleMouseDown(button, x, y)
    })
}

// In gogpu/gogpu - implements EventSource
type App struct {
    keyHandlers []func(gpucontext.Key, gpucontext.Modifiers)
}

func (app *App) OnKeyPress(fn func(gpucontext.Key, gpucontext.Modifiers)) {
    app.keyHandlers = append(app.keyHandlers, fn)
}
```

### IME Support (CJK Input)

`IMEState` and related interfaces enable Input Method Editor support for Chinese, Japanese, and Korean input:

```go
// In gogpu/ui - handle IME composition
func (input *TextInput) AttachIME(source gpucontext.EventSource) {
    source.OnIMECompositionStart(func() {
        input.showCompositionWindow()
    })

    source.OnIMECompositionUpdate(func(state gpucontext.IMEState) {
        // Show composition text with cursor
        input.setCompositionText(state.CompositionText, state.CursorPos)
    })

    source.OnIMECompositionEnd(func(committed string) {
        // Insert final text
        input.insertText(committed)
        input.hideCompositionWindow()
    })
}

// Control IME position (for composition window placement)
func (input *TextInput) Focus(controller gpucontext.IMEController) {
    controller.SetIMEEnabled(true)
    controller.SetIMEPosition(input.cursorX, input.cursorY)
}
```

### Backend Registry

The `Registry[T]` provides thread-safe registration with priority-based selection:

```go
import "github.com/gogpu/gpucontext"

// Create registry with priority order
var backends = gpucontext.NewRegistry[Backend](
    gpucontext.WithPriority("vulkan", "dx12", "metal", "gles", "software"),
)

// Register backends (typically in init())
func init() {
    backends.Register("vulkan", NewVulkanBackend)
    backends.Register("software", NewSoftwareBackend)
}

// Get best available backend
backend := backends.Best()

// Or get specific backend
vulkan := backends.Get("vulkan")

// Check availability
if backends.Has("vulkan") {
    // Vulkan is available
}

// List all available
names := backends.Available() // ["vulkan", "software"]
```

## Dependency Graph

```
                   gpucontext
                  (DeviceProvider, Registry, types)
                       ^
          +------------+------------+
          |            |            |
        gogpu         gg         born-ml/born
        (implements)  (uses)     (implements & uses)
          |
          v
        wgpu/hal
```

## Ecosystem

| Package | Description |
|---------|-------------|
| [gogpu/gogpu](https://github.com/gogpu/gogpu) | Graphics framework, implements DeviceProvider |
| [gogpu/gg](https://github.com/gogpu/gg) | 2D graphics, uses DeviceProvider |
| [gogpu/wgpu](https://github.com/gogpu/wgpu) | Pure Go WebGPU implementation |
| [born-ml/born](https://github.com/born-ml/born) | ML framework, implements & uses |

## License

MIT License — see [LICENSE](LICENSE) for details.
