// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

package gpucontext

import "testing"

func TestNullWindowProvider_Size(t *testing.T) {
	tests := []struct {
		name       string
		provider   NullWindowProvider
		wantWidth  int
		wantHeight int
	}{
		{"zero value", NullWindowProvider{}, 0, 0},
		{"standard HD", NullWindowProvider{W: 1920, H: 1080}, 1920, 1080},
		{"small window", NullWindowProvider{W: 320, H: 240}, 320, 240},
		{"4K", NullWindowProvider{W: 3840, H: 2160}, 3840, 2160},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, h := tt.provider.Size()
			if w != tt.wantWidth {
				t.Errorf("width = %d, want %d", w, tt.wantWidth)
			}
			if h != tt.wantHeight {
				t.Errorf("height = %d, want %d", h, tt.wantHeight)
			}
		})
	}
}

func TestNullWindowProvider_ScaleFactor(t *testing.T) {
	tests := []struct {
		name string
		sf   float64
		want float64
	}{
		{"zero defaults to 1.0", 0, 1.0},
		{"standard", 1.0, 1.0},
		{"retina", 2.0, 2.0},
		{"high DPI", 1.5, 1.5},
		{"3x scale", 3.0, 3.0},
		{"fractional", 1.25, 1.25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := NullWindowProvider{SF: tt.sf}
			got := provider.ScaleFactor()
			if got != tt.want {
				t.Errorf("ScaleFactor() = %f, want %f", got, tt.want)
			}
		})
	}
}

func TestNullWindowProvider_RequestRedraw(t *testing.T) {
	// RequestRedraw should not panic
	provider := NullWindowProvider{}
	provider.RequestRedraw()
}

func TestNullWindowProvider_Interface(t *testing.T) {
	// Verify NullWindowProvider can be used through the WindowProvider interface
	var wp WindowProvider = NullWindowProvider{W: 800, H: 600, SF: 2.0}

	w, h := wp.Size()
	if w != 800 {
		t.Errorf("width = %d, want 800", w)
	}
	if h != 600 {
		t.Errorf("height = %d, want 600", h)
	}

	sf := wp.ScaleFactor()
	if sf != 2.0 {
		t.Errorf("ScaleFactor() = %f, want 2.0", sf)
	}

	// Should not panic
	wp.RequestRedraw()
}

// mockWindowProvider verifies the interface can be implemented by custom types.
type mockWindowProvider struct {
	w, h    int
	sf      float64
	redraws int
}

func (m *mockWindowProvider) Size() (int, int)     { return m.w, m.h }
func (m *mockWindowProvider) ScaleFactor() float64 { return m.sf }
func (m *mockWindowProvider) RequestRedraw()       { m.redraws++ }

// Ensure mockWindowProvider implements WindowProvider.
var _ WindowProvider = &mockWindowProvider{}

func TestWindowProvider_CustomImplementation(t *testing.T) {
	mock := &mockWindowProvider{w: 1024, h: 768, sf: 1.5}
	var wp WindowProvider = mock

	w, h := wp.Size()
	if w != 1024 || h != 768 {
		t.Errorf("Size() = (%d, %d), want (1024, 768)", w, h)
	}

	if sf := wp.ScaleFactor(); sf != 1.5 {
		t.Errorf("ScaleFactor() = %f, want 1.5", sf)
	}

	wp.RequestRedraw()
	wp.RequestRedraw()
	if mock.redraws != 2 {
		t.Errorf("redraws = %d, want 2", mock.redraws)
	}
}
