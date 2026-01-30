// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

package gpucontext

// Texture is the minimal interface for GPU textures.
// This interface enables type-safe cross-package texture handling
// without circular dependencies.
//
// Implementations:
//   - gogpu.Texture implements Texture
//
// Design note: This interface intentionally contains only read-only
// methods that are universally needed for texture operations.
// Implementation-specific methods (like UpdateData) remain on
// concrete types.
type Texture interface {
	// Width returns the texture width in pixels.
	Width() int

	// Height returns the texture height in pixels.
	Height() int
}

// TextureDrawer provides texture drawing capabilities for 2D rendering.
// This interface enables packages like ggcanvas to draw textures without
// depending directly on gogpu, following the Dependency Inversion Principle.
//
// Implementations:
//   - gogpu.Context implements TextureDrawer (via adapter)
//
// Example usage in ggcanvas:
//
//	func (c *Canvas) RenderTo(drawer gpucontext.TextureDrawer) error {
//	    tex, _ := c.Flush()
//	    return drawer.DrawTexture(tex, 0, 0)
//	}
type TextureDrawer interface {
	// DrawTexture draws a texture at the specified position.
	//
	// Coordinate system:
	//   - (0, 0) is the top-left corner
	//   - Positive X goes right
	//   - Positive Y goes down
	//   - Coordinates are in pixels
	//
	// The texture must have been created by TextureCreator from this drawer.
	DrawTexture(tex Texture, x, y float32) error

	// TextureCreator returns the texture creator associated with this drawer.
	// Use this to create textures that can be drawn by this drawer.
	TextureCreator() TextureCreator
}

// TextureCreator provides texture creation from raw pixel data.
// This interface enables packages to create GPU textures without
// depending directly on specific GPU implementations.
//
// Implementations:
//   - gogpu.Renderer implements TextureCreator (via adapter)
//
// Example usage:
//
//	creator := drawer.TextureCreator()
//	tex, err := creator.NewTextureFromRGBA(800, 600, pixelData)
//	if err != nil {
//	    return err
//	}
//	drawer.DrawTexture(tex, 0, 0)
type TextureCreator interface {
	// NewTextureFromRGBA creates a texture from RGBA pixel data.
	// The data must be width * height * 4 bytes (RGBA, 8 bits per channel).
	//
	// The returned Texture can be drawn using TextureDrawer.DrawTexture.
	//
	// Returns error if:
	//   - Data size doesn't match width * height * 4
	//   - GPU texture creation fails
	NewTextureFromRGBA(width, height int, data []byte) (Texture, error)
}
