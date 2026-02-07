# gpucontext

Shared GPU infrastructure for the [gogpu](https://github.com/gogpu) ecosystem.

## Overview

`gpucontext` provides interfaces and utilities for sharing GPU resources across multiple packages without circular dependencies.

## Relationship to gputypes

| Package | Purpose | Dependencies |
|---------|---------|--------------|
| [gputypes](https://github.com/gogpu/gputypes) | WebGPU types (enums, structs, constants) | **ZERO** |
| **gpucontext** | Interfaces (DeviceProvider, EventSource, Texture) | imports gputypes |

gpucontext imports gputypes to use shared types in interface signatures, ensuring type compatibility across the ecosystem.

## Installation

```bash
go get github.com/gogpu/gpucontext
```

**Requires:** Go 1.25+

## Features

- **DeviceProvider** — Interface for injecting GPU device and queue
- **WindowProvider** — Window geometry, DPI scale factor, and redraw requests
- **PlatformProvider** — Clipboard, cursor, dark mode, and accessibility preferences
- **CursorShape** — 12 standard cursor shapes (arrow, pointer, text, resize, etc.)
- **EventSource** — Interface for input events (keyboard, mouse, window, IME)
- **PointerEventSource** — W3C Pointer Events Level 3 (unified mouse/touch/pen)
- **ScrollEventSource** — Scroll/wheel events with pixel/line/page modes
- **Texture** — Minimal interface for GPU textures with TextureUpdater/TextureDrawer/TextureCreator
- **IME Support** — Input Method Editor for CJK languages (Chinese, Japanese, Korean)
- **Registry[T]** — Generic registry with priority-based backend selection
- **WebGPU Interfaces** — Device, Queue, Adapter, Surface interfaces
- **WebGPU Types** — Re-exports from [gputypes](https://github.com/gogpu/gputypes) (TextureFormat, etc.)

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

### WindowProvider (for UI frameworks)

The `WindowProvider` interface enables UI frameworks to query window dimensions and DPI:

```go
// In gogpu/ui - uses WindowProvider for layout
func (ui *UI) Layout(wp gpucontext.WindowProvider) {
    w, h := wp.Size()
    scale := wp.ScaleFactor()
    dpW := float64(w) / scale
    dpH := float64(h) / scale
    ui.root.Layout(dpW, dpH)
}
```

### PlatformProvider (optional OS integration)

`PlatformProvider` exposes clipboard, cursor, and system preferences.
Not all hosts support it — use type assertion to check:

```go
// In gogpu/ui - cursor management
func (ui *UI) UpdateCursor(provider gpucontext.WindowProvider) {
    if pp, ok := provider.(gpucontext.PlatformProvider); ok {
        pp.SetCursor(gpucontext.CursorPointer) // hand cursor
    }
}

// In gogpu/ui - clipboard
func (ui *UI) Paste(provider gpucontext.WindowProvider) {
    if pp, ok := provider.(gpucontext.PlatformProvider); ok {
        text, err := pp.ClipboardRead()
        if err == nil {
            ui.focused.InsertText(text)
        }
    }
}

// In gogpu/ui - theme detection
func (ui *UI) DetectTheme(provider gpucontext.WindowProvider) {
    if pp, ok := provider.(gpucontext.PlatformProvider); ok {
        if pp.DarkMode() {
            ui.SetTheme(DarkTheme)
        }
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

### Texture Interface

`Texture` provides a minimal interface for GPU textures, enabling sharing between packages:

```go
// Texture is a minimal interface for GPU textures
type Texture interface {
    Width() int
    Height() int
}

// TextureDrawer can draw textures (implemented by renderers)
type TextureDrawer interface {
    DrawTexture(tex Texture, x, y float32) error
    DrawTextureEx(tex Texture, opts TextureDrawOptions) error
}

// TextureCreator can create textures from pixel data
type TextureCreator interface {
    CreateTexture(width, height int, pixels []byte) (Texture, error)
}
```

### TextureUpdater (Dynamic Content)

`TextureUpdater` enables efficient texture updates without recreating textures:

```go
// TextureUpdater updates existing texture pixel data
type TextureUpdater interface {
    UpdateData(data []byte) error
}
```

Usage in integration packages:

```go
// In gg/integration/ggcanvas - creates textures from CPU canvas
func (c *Canvas) Flush() (gpucontext.Texture, error) {
    pixels := c.pixmap.Pix()
    return c.creator.CreateTexture(c.width, c.height, pixels)
}

// In gogpu - implements TextureDrawer
func (ctx *Context) DrawTexture(tex gpucontext.Texture, x, y float32) error {
    return ctx.renderer.DrawTexture(tex, x, y)
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
                   gputypes (ZERO deps)
                 All WebGPU types (100+)
                          │
                          ▼
                   gpucontext
                  (imports gputypes)
          DeviceProvider, WindowProvider,
          PlatformProvider, EventSource,
          Texture, PointerEventSource, Registry
                          │
          ┌───────────────┼───────────────┐
          │               │               │
          ▼               ▼               ▼
        gogpu            gg              ui
     (implements)      (uses)         (uses)
          │
          ▼
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
