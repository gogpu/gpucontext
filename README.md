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
- **EventSource** — Interface for input events (keyboard, mouse, window, IME)
- **PointerEventSource** — W3C Pointer Events Level 3 (unified mouse/touch/pen)
- **ScrollEventSource** — Scroll/wheel events with pixel/line/page modes
- **TouchEventSource** — Interface for multi-touch input (mobile, tablets, touchscreens)
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

### Touch Input (Multi-touch Support)

`TouchEventSource` enables multi-touch handling for mobile and tablet applications:

```go
// Touch phases follow platform conventions (iOS, Android, W3C)
const (
    TouchBegan     // First contact
    TouchMoved     // Touch moved
    TouchEnded     // Touch lifted
    TouchCanceled // System interrupted
)

// TouchPoint represents a single touch contact
type TouchPoint struct {
    ID       TouchID   // Unique within session
    X, Y     float64   // Position in logical pixels
    Pressure *float32  // Optional: 0.0-1.0
    Radius   *float32  // Optional: contact radius
}

// TouchEvent contains all touch information
type TouchEvent struct {
    Phase     TouchPhase    // Lifecycle stage
    Changed   []TouchPoint  // Touches that triggered this event
    All       []TouchPoint  // All active touches
    Modifiers Modifiers     // Keyboard modifiers (Ctrl+drag, etc.)
    Timestamp time.Duration // For velocity calculations
}
```

Usage for gesture handling:

```go
// Implement pinch-to-zoom
func (app *App) AttachTouchEvents(source gpucontext.EventSource) {
    // Check if touch is supported
    if tes, ok := source.(gpucontext.TouchEventSource); ok {
        tes.OnTouch(func(ev gpucontext.TouchEvent) {
            switch ev.Phase {
            case gpucontext.TouchBegan:
                app.startGesture(ev.Changed)
            case gpucontext.TouchMoved:
                if len(ev.All) == 2 {
                    // Pinch gesture
                    app.handlePinch(ev.All[0], ev.All[1])
                } else if len(ev.All) == 1 {
                    // Pan gesture
                    app.handlePan(ev.All[0])
                }
            case gpucontext.TouchEnded, gpucontext.TouchCanceled:
                app.endGesture()
            }
        })
    }
}

// Calculate pinch distance
func (app *App) handlePinch(t1, t2 gpucontext.TouchPoint) {
    dx := t1.X - t2.X
    dy := t1.Y - t2.Y
    distance := math.Sqrt(dx*dx + dy*dy)
    app.zoom = distance / app.initialPinchDistance
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
          DeviceProvider, EventSource, Texture
              TouchEventSource, Registry
                          │
          ┌───────────────┼───────────────┐
          │               │               │
          ▼               ▼               ▼
        gogpu            gg          born-ml/born
     (implements)      (uses)      (implements & uses)
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
