// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

package gpucontext

import (
	"testing"
	"time"
)

func TestGestureEvent_ZeroValue(t *testing.T) {
	var ev GestureEvent

	// Zero value should represent no gesture
	if ev.NumPointers != 0 {
		t.Errorf("NumPointers: got %d, want 0", ev.NumPointers)
	}
	if ev.ZoomDelta != 0 {
		t.Errorf("ZoomDelta: got %f, want 0", ev.ZoomDelta)
	}
	if ev.RotationDelta != 0 {
		t.Errorf("RotationDelta: got %f, want 0", ev.RotationDelta)
	}
	if ev.TranslationDelta.X != 0 || ev.TranslationDelta.Y != 0 {
		t.Errorf("TranslationDelta: got (%f, %f), want (0, 0)",
			ev.TranslationDelta.X, ev.TranslationDelta.Y)
	}
	if ev.PinchType != PinchNone {
		t.Errorf("PinchType: got %v, want PinchNone", ev.PinchType)
	}
	if ev.Center.X != 0 || ev.Center.Y != 0 {
		t.Errorf("Center: got (%f, %f), want (0, 0)", ev.Center.X, ev.Center.Y)
	}
	if ev.Timestamp != 0 {
		t.Errorf("Timestamp: got %v, want 0", ev.Timestamp)
	}
}

func TestGestureEvent_Fields(t *testing.T) {
	ev := GestureEvent{
		NumPointers:      2,
		ZoomDelta:        1.5,
		ZoomDelta2D:      Point{X: 1.5, Y: 1.0},
		RotationDelta:    0.1,
		TranslationDelta: Point{X: 10, Y: 20},
		PinchType:        PinchProportional,
		Center:           Point{X: 100, Y: 200},
		Timestamp:        time.Second,
	}

	if ev.NumPointers != 2 {
		t.Errorf("NumPointers: got %d, want 2", ev.NumPointers)
	}
	if ev.ZoomDelta != 1.5 {
		t.Errorf("ZoomDelta: got %f, want 1.5", ev.ZoomDelta)
	}
	if ev.ZoomDelta2D.X != 1.5 || ev.ZoomDelta2D.Y != 1.0 {
		t.Errorf("ZoomDelta2D: got (%f, %f), want (1.5, 1.0)",
			ev.ZoomDelta2D.X, ev.ZoomDelta2D.Y)
	}
	if ev.RotationDelta != 0.1 {
		t.Errorf("RotationDelta: got %f, want 0.1", ev.RotationDelta)
	}
	if ev.TranslationDelta.X != 10 || ev.TranslationDelta.Y != 20 {
		t.Errorf("TranslationDelta: got (%f, %f), want (10, 20)",
			ev.TranslationDelta.X, ev.TranslationDelta.Y)
	}
	if ev.PinchType != PinchProportional {
		t.Errorf("PinchType: got %v, want PinchProportional", ev.PinchType)
	}
	if ev.Center.X != 100 || ev.Center.Y != 200 {
		t.Errorf("Center: got (%f, %f), want (100, 200)", ev.Center.X, ev.Center.Y)
	}
	if ev.Timestamp != time.Second {
		t.Errorf("Timestamp: got %v, want 1s", ev.Timestamp)
	}
}

func TestPinchType_String(t *testing.T) {
	tests := []struct {
		pinchType PinchType
		want      string
	}{
		{PinchNone, "None"},
		{PinchHorizontal, "Horizontal"},
		{PinchVertical, "Vertical"},
		{PinchProportional, "Proportional"},
		{PinchType(99), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.pinchType.String()
			if got != tt.want {
				t.Errorf("String(): got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestPinchType_Classification(t *testing.T) {
	// Test the classification logic that will be used in GestureRecognizer
	classifyPinch := func(dx, dy float64) PinchType {
		absDx := dx
		if absDx < 0 {
			absDx = -absDx
		}
		absDy := dy
		if absDy < 0 {
			absDy = -absDy
		}

		if absDx > absDy*3 {
			return PinchHorizontal
		}
		if absDy > absDx*3 {
			return PinchVertical
		}
		return PinchProportional
	}

	tests := []struct {
		name string
		dx   float64
		dy   float64
		want PinchType
	}{
		{"horizontal dominant", 100, 10, PinchHorizontal},
		{"vertical dominant", 10, 100, PinchVertical},
		{"proportional equal", 50, 50, PinchProportional},
		{"proportional similar", 60, 50, PinchProportional},
		{"horizontal exactly 3x", 30, 10, PinchProportional}, // Not > 3x
		{"horizontal over 3x", 31, 10, PinchHorizontal},
		{"vertical exactly 3x", 10, 30, PinchProportional}, // Not > 3x
		{"vertical over 3x", 10, 31, PinchVertical},
		{"negative horizontal", -100, 10, PinchHorizontal},
		{"negative vertical", 10, -100, PinchVertical},
		{"both negative", -100, -100, PinchProportional},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := classifyPinch(tt.dx, tt.dy)
			if got != tt.want {
				t.Errorf("classifyPinch(%f, %f): got %v, want %v",
					tt.dx, tt.dy, got, tt.want)
			}
		})
	}
}

func TestPoint_Operations(t *testing.T) {
	p1 := Point{X: 10, Y: 20}
	p2 := Point{X: 5, Y: 10}

	// Test Add
	sum := p1.Add(p2)
	if sum.X != 15 || sum.Y != 30 {
		t.Errorf("Add: got (%f, %f), want (15, 30)", sum.X, sum.Y)
	}

	// Test Sub
	diff := p1.Sub(p2)
	if diff.X != 5 || diff.Y != 10 {
		t.Errorf("Sub: got (%f, %f), want (5, 10)", diff.X, diff.Y)
	}

	// Test Scale
	scaled := p1.Scale(2)
	if scaled.X != 20 || scaled.Y != 40 {
		t.Errorf("Scale: got (%f, %f), want (20, 40)", scaled.X, scaled.Y)
	}

	// Test Scale with negative
	scaledNeg := p1.Scale(-1)
	if scaledNeg.X != -10 || scaledNeg.Y != -20 {
		t.Errorf("Scale(-1): got (%f, %f), want (-10, -20)", scaledNeg.X, scaledNeg.Y)
	}
}

func TestNullGestureEventSource(t *testing.T) {
	var source NullGestureEventSource

	// Should not panic
	called := false
	source.OnGesture(func(GestureEvent) {
		called = true
	})

	if called {
		t.Error("NullGestureEventSource should not call the callback")
	}

	// Verify interface compliance
	var _ GestureEventSource = source
	var _ GestureEventSource = NullGestureEventSource{}
}
