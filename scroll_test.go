// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

package gpucontext

import (
	"testing"
	"time"
)

func TestScrollDeltaMode_String(t *testing.T) {
	tests := []struct {
		mode ScrollDeltaMode
		want string
	}{
		{ScrollDeltaPixel, "Pixel"},
		{ScrollDeltaLine, "Line"},
		{ScrollDeltaPage, "Page"},
		{ScrollDeltaMode(99), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.mode.String(); got != tt.want {
				t.Errorf("ScrollDeltaMode(%d).String() = %q, want %q", tt.mode, got, tt.want)
			}
		})
	}
}

func TestScrollDeltaMode_Values(t *testing.T) {
	// Verify delta mode constants are sequential
	if ScrollDeltaPixel != 0 {
		t.Errorf("ScrollDeltaPixel = %d, want 0", ScrollDeltaPixel)
	}
	if ScrollDeltaLine != 1 {
		t.Errorf("ScrollDeltaLine = %d, want 1", ScrollDeltaLine)
	}
	if ScrollDeltaPage != 2 {
		t.Errorf("ScrollDeltaPage = %d, want 2", ScrollDeltaPage)
	}
}

func TestScrollEvent_ZeroValue(t *testing.T) {
	var ev ScrollEvent

	if ev.X != 0 {
		t.Errorf("Zero value X = %f, want 0", ev.X)
	}
	if ev.Y != 0 {
		t.Errorf("Zero value Y = %f, want 0", ev.Y)
	}
	if ev.DeltaX != 0 {
		t.Errorf("Zero value DeltaX = %f, want 0", ev.DeltaX)
	}
	if ev.DeltaY != 0 {
		t.Errorf("Zero value DeltaY = %f, want 0", ev.DeltaY)
	}
	if ev.DeltaMode != ScrollDeltaPixel {
		t.Errorf("Zero value DeltaMode = %v, want ScrollDeltaPixel", ev.DeltaMode)
	}
	if ev.Modifiers != 0 {
		t.Errorf("Zero value Modifiers = %d, want 0", ev.Modifiers)
	}
	if ev.Timestamp != 0 {
		t.Errorf("Zero value Timestamp = %v, want 0", ev.Timestamp)
	}
}

func TestScrollEvent_FullConstruction(t *testing.T) {
	ev := ScrollEvent{
		X:         100.5,
		Y:         200.5,
		DeltaX:    10.0,
		DeltaY:    -20.0,
		DeltaMode: ScrollDeltaLine,
		Modifiers: ModControl,
		Timestamp: time.Millisecond * 5000,
	}

	if ev.X != 100.5 {
		t.Errorf("X = %f, want 100.5", ev.X)
	}
	if ev.Y != 200.5 {
		t.Errorf("Y = %f, want 200.5", ev.Y)
	}
	if ev.DeltaX != 10.0 {
		t.Errorf("DeltaX = %f, want 10.0", ev.DeltaX)
	}
	if ev.DeltaY != -20.0 {
		t.Errorf("DeltaY = %f, want -20.0", ev.DeltaY)
	}
	if ev.DeltaMode != ScrollDeltaLine {
		t.Errorf("DeltaMode = %v, want ScrollDeltaLine", ev.DeltaMode)
	}
	if !ev.Modifiers.HasControl() {
		t.Error("Modifiers should have control")
	}
	if ev.Timestamp != time.Millisecond*5000 {
		t.Errorf("Timestamp = %v, want %v", ev.Timestamp, time.Millisecond*5000)
	}
}

func TestScrollEvent_VerticalScroll(t *testing.T) {
	// Typical vertical scroll event from mouse wheel
	ev := ScrollEvent{
		X:         400,
		Y:         300,
		DeltaX:    0,
		DeltaY:    3, // Scroll down 3 lines
		DeltaMode: ScrollDeltaLine,
	}

	if ev.DeltaX != 0 {
		t.Errorf("DeltaX = %f, want 0", ev.DeltaX)
	}
	if ev.DeltaY != 3 {
		t.Errorf("DeltaY = %f, want 3", ev.DeltaY)
	}
	if ev.DeltaMode != ScrollDeltaLine {
		t.Errorf("DeltaMode = %v, want ScrollDeltaLine", ev.DeltaMode)
	}
}

func TestScrollEvent_HorizontalScroll(t *testing.T) {
	// Horizontal scroll from trackpad
	ev := ScrollEvent{
		X:         400,
		Y:         300,
		DeltaX:    50, // Scroll right 50 pixels
		DeltaY:    0,
		DeltaMode: ScrollDeltaPixel,
	}

	if ev.DeltaX != 50 {
		t.Errorf("DeltaX = %f, want 50", ev.DeltaX)
	}
	if ev.DeltaY != 0 {
		t.Errorf("DeltaY = %f, want 0", ev.DeltaY)
	}
	if ev.DeltaMode != ScrollDeltaPixel {
		t.Errorf("DeltaMode = %v, want ScrollDeltaPixel", ev.DeltaMode)
	}
}

func TestScrollEvent_PageScroll(t *testing.T) {
	// Page scroll (rare, but possible)
	ev := ScrollEvent{
		X:         400,
		Y:         300,
		DeltaX:    0,
		DeltaY:    1, // Scroll down 1 page
		DeltaMode: ScrollDeltaPage,
	}

	if ev.DeltaY != 1 {
		t.Errorf("DeltaY = %f, want 1", ev.DeltaY)
	}
	if ev.DeltaMode != ScrollDeltaPage {
		t.Errorf("DeltaMode = %v, want ScrollDeltaPage", ev.DeltaMode)
	}
}

func TestScrollEvent_CtrlScroll(t *testing.T) {
	// Ctrl+scroll typically means zoom
	ev := ScrollEvent{
		X:         400,
		Y:         300,
		DeltaX:    0,
		DeltaY:    -1,
		DeltaMode: ScrollDeltaLine,
		Modifiers: ModControl,
	}

	if !ev.Modifiers.HasControl() {
		t.Error("Should detect Ctrl modifier for zoom detection")
	}
}

func TestNullScrollEventSource(t *testing.T) {
	// NullScrollEventSource should implement ScrollEventSource
	var ses ScrollEventSource = NullScrollEventSource{}

	// Method should be callable without panic
	ses.OnScrollEvent(func(ScrollEvent) {})
}

// mockScrollEventSource is used to verify ScrollEventSource interface.
type mockScrollEventSource struct {
	handler func(ScrollEvent)
}

func (m *mockScrollEventSource) OnScrollEvent(fn func(ScrollEvent)) {
	m.handler = fn
}

// Ensure mockScrollEventSource implements ScrollEventSource.
var _ ScrollEventSource = &mockScrollEventSource{}

func TestScrollEventSource_Interface(t *testing.T) {
	// Verify ScrollEventSource can be used through the interface
	mock := &mockScrollEventSource{}
	var source ScrollEventSource = mock

	var received *ScrollEvent
	source.OnScrollEvent(func(ev ScrollEvent) {
		received = &ev
	})

	// Simulate event dispatch
	testEvent := ScrollEvent{
		X:         100,
		Y:         200,
		DeltaX:    0,
		DeltaY:    -120, // Scroll up 120 pixels
		DeltaMode: ScrollDeltaPixel,
	}

	mock.handler(testEvent)

	if received == nil {
		t.Fatal("Handler was not called")
	}
	if received.X != 100 {
		t.Errorf("received.X = %f, want 100", received.X)
	}
	if received.Y != 200 {
		t.Errorf("received.Y = %f, want 200", received.Y)
	}
	if received.DeltaY != -120 {
		t.Errorf("received.DeltaY = %f, want -120", received.DeltaY)
	}
}
