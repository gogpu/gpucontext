// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

package gpucontext

import (
	"testing"
	"time"
)

func TestPointerEventType_String(t *testing.T) {
	tests := []struct {
		eventType PointerEventType
		want      string
	}{
		{PointerDown, "PointerDown"},
		{PointerUp, "PointerUp"},
		{PointerMove, "PointerMove"},
		{PointerEnter, "PointerEnter"},
		{PointerLeave, "PointerLeave"},
		{PointerCancel, "PointerCancel"},
		{PointerEventType(99), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.eventType.String(); got != tt.want {
				t.Errorf("PointerEventType(%d).String() = %q, want %q", tt.eventType, got, tt.want)
			}
		})
	}
}

func TestPointerType_String(t *testing.T) {
	tests := []struct {
		pointerType PointerType
		want        string
	}{
		{PointerTypeMouse, "Mouse"},
		{PointerTypeTouch, "Touch"},
		{PointerTypePen, "Pen"},
		{PointerType(99), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.pointerType.String(); got != tt.want {
				t.Errorf("PointerType(%d).String() = %q, want %q", tt.pointerType, got, tt.want)
			}
		})
	}
}

func TestButton_String(t *testing.T) {
	tests := []struct {
		button Button
		want   string
	}{
		{ButtonNone, "None"},
		{ButtonLeft, "Left"},
		{ButtonMiddle, "Middle"},
		{ButtonRight, "Right"},
		{ButtonX1, "X1"},
		{ButtonX2, "X2"},
		{ButtonEraser, "Eraser"},
		{Button(99), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.button.String(); got != tt.want {
				t.Errorf("Button(%d).String() = %q, want %q", tt.button, got, tt.want)
			}
		})
	}
}

func TestButtonConstants(t *testing.T) {
	// Verify W3C-compliant button values
	if ButtonNone != -1 {
		t.Errorf("ButtonNone = %d, want -1", ButtonNone)
	}
	if ButtonLeft != 0 {
		t.Errorf("ButtonLeft = %d, want 0", ButtonLeft)
	}
	if ButtonMiddle != 1 {
		t.Errorf("ButtonMiddle = %d, want 1", ButtonMiddle)
	}
	if ButtonRight != 2 {
		t.Errorf("ButtonRight = %d, want 2", ButtonRight)
	}
	if ButtonX1 != 3 {
		t.Errorf("ButtonX1 = %d, want 3", ButtonX1)
	}
	if ButtonX2 != 4 {
		t.Errorf("ButtonX2 = %d, want 4", ButtonX2)
	}
}

func TestButtons_HasMethods(t *testing.T) {
	tests := []struct {
		name    string
		buttons Buttons
		left    bool
		right   bool
		middle  bool
		x1      bool
		x2      bool
		eraser  bool
	}{
		{"none", ButtonsNone, false, false, false, false, false, false},
		{"left only", ButtonsLeft, true, false, false, false, false, false},
		{"right only", ButtonsRight, false, true, false, false, false, false},
		{"middle only", ButtonsMiddle, false, false, true, false, false, false},
		{"x1 only", ButtonsX1, false, false, false, true, false, false},
		{"x2 only", ButtonsX2, false, false, false, false, true, false},
		{"eraser only", ButtonsEraser, false, false, false, false, false, true},
		{"left and right", ButtonsLeft | ButtonsRight, true, true, false, false, false, false},
		{"all buttons", ButtonsLeft | ButtonsRight | ButtonsMiddle | ButtonsX1 | ButtonsX2 | ButtonsEraser,
			true, true, true, true, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.buttons.HasLeft(); got != tt.left {
				t.Errorf("Buttons.HasLeft() = %v, want %v", got, tt.left)
			}
			if got := tt.buttons.HasRight(); got != tt.right {
				t.Errorf("Buttons.HasRight() = %v, want %v", got, tt.right)
			}
			if got := tt.buttons.HasMiddle(); got != tt.middle {
				t.Errorf("Buttons.HasMiddle() = %v, want %v", got, tt.middle)
			}
			if got := tt.buttons.HasX1(); got != tt.x1 {
				t.Errorf("Buttons.HasX1() = %v, want %v", got, tt.x1)
			}
			if got := tt.buttons.HasX2(); got != tt.x2 {
				t.Errorf("Buttons.HasX2() = %v, want %v", got, tt.x2)
			}
			if got := tt.buttons.HasEraser(); got != tt.eraser {
				t.Errorf("Buttons.HasEraser() = %v, want %v", got, tt.eraser)
			}
		})
	}
}

func TestButtons_Count(t *testing.T) {
	tests := []struct {
		buttons Buttons
		count   int
	}{
		{ButtonsNone, 0},
		{ButtonsLeft, 1},
		{ButtonsLeft | ButtonsRight, 2},
		{ButtonsLeft | ButtonsRight | ButtonsMiddle, 3},
		{ButtonsLeft | ButtonsRight | ButtonsMiddle | ButtonsX1 | ButtonsX2 | ButtonsEraser, 6},
	}

	for _, tt := range tests {
		if got := tt.buttons.Count(); got != tt.count {
			t.Errorf("Buttons(%d).Count() = %d, want %d", tt.buttons, got, tt.count)
		}
	}
}

func TestButtonsConstants(t *testing.T) {
	// Verify button bitmask values
	if ButtonsNone != 0 {
		t.Errorf("ButtonsNone = %d, want 0", ButtonsNone)
	}
	if ButtonsLeft != 1 {
		t.Errorf("ButtonsLeft = %d, want 1", ButtonsLeft)
	}
	if ButtonsRight != 2 {
		t.Errorf("ButtonsRight = %d, want 2", ButtonsRight)
	}
	if ButtonsMiddle != 4 {
		t.Errorf("ButtonsMiddle = %d, want 4", ButtonsMiddle)
	}
	if ButtonsX1 != 8 {
		t.Errorf("ButtonsX1 = %d, want 8", ButtonsX1)
	}
	if ButtonsX2 != 16 {
		t.Errorf("ButtonsX2 = %d, want 16", ButtonsX2)
	}
	if ButtonsEraser != 32 {
		t.Errorf("ButtonsEraser = %d, want 32", ButtonsEraser)
	}
}

func TestPointerEvent_ZeroValue(t *testing.T) {
	var ev PointerEvent

	if ev.Type != PointerDown {
		t.Errorf("Zero value Type = %v, want PointerDown", ev.Type)
	}
	if ev.PointerID != 0 {
		t.Errorf("Zero value PointerID = %d, want 0", ev.PointerID)
	}
	if ev.X != 0 {
		t.Errorf("Zero value X = %f, want 0", ev.X)
	}
	if ev.Y != 0 {
		t.Errorf("Zero value Y = %f, want 0", ev.Y)
	}
	if ev.Pressure != 0 {
		t.Errorf("Zero value Pressure = %f, want 0", ev.Pressure)
	}
	if ev.PointerType != PointerTypeMouse {
		t.Errorf("Zero value PointerType = %v, want PointerTypeMouse", ev.PointerType)
	}
	if ev.IsPrimary {
		t.Error("Zero value IsPrimary should be false")
	}
}

func TestPointerEvent_FullConstruction(t *testing.T) {
	ev := PointerEvent{
		Type:        PointerMove,
		PointerID:   42,
		X:           100.5,
		Y:           200.5,
		Pressure:    0.75,
		TiltX:       15.0,
		TiltY:       -10.0,
		Twist:       45.0,
		Width:       20.0,
		Height:      25.0,
		PointerType: PointerTypePen,
		IsPrimary:   true,
		Button:      ButtonNone,
		Buttons:     ButtonsLeft | ButtonsMiddle,
		Modifiers:   ModShift | ModControl,
		Timestamp:   time.Millisecond * 12345,
	}

	if ev.Type != PointerMove {
		t.Errorf("Type = %v, want PointerMove", ev.Type)
	}
	if ev.PointerID != 42 {
		t.Errorf("PointerID = %d, want 42", ev.PointerID)
	}
	if ev.X != 100.5 {
		t.Errorf("X = %f, want 100.5", ev.X)
	}
	if ev.Y != 200.5 {
		t.Errorf("Y = %f, want 200.5", ev.Y)
	}
	if ev.Pressure != 0.75 {
		t.Errorf("Pressure = %f, want 0.75", ev.Pressure)
	}
	if ev.TiltX != 15.0 {
		t.Errorf("TiltX = %f, want 15.0", ev.TiltX)
	}
	if ev.TiltY != -10.0 {
		t.Errorf("TiltY = %f, want -10.0", ev.TiltY)
	}
	if ev.Twist != 45.0 {
		t.Errorf("Twist = %f, want 45.0", ev.Twist)
	}
	if ev.Width != 20.0 {
		t.Errorf("Width = %f, want 20.0", ev.Width)
	}
	if ev.Height != 25.0 {
		t.Errorf("Height = %f, want 25.0", ev.Height)
	}
	if ev.PointerType != PointerTypePen {
		t.Errorf("PointerType = %v, want PointerTypePen", ev.PointerType)
	}
	if !ev.IsPrimary {
		t.Error("IsPrimary should be true")
	}
	if ev.Button != ButtonNone {
		t.Errorf("Button = %v, want ButtonNone", ev.Button)
	}
	if !ev.Buttons.HasLeft() {
		t.Error("Buttons should have left")
	}
	if !ev.Buttons.HasMiddle() {
		t.Error("Buttons should have middle")
	}
	if !ev.Modifiers.HasShift() {
		t.Error("Modifiers should have shift")
	}
	if !ev.Modifiers.HasControl() {
		t.Error("Modifiers should have control")
	}
	if ev.Timestamp != time.Millisecond*12345 {
		t.Errorf("Timestamp = %v, want %v", ev.Timestamp, time.Millisecond*12345)
	}
}

func TestPointerEventType_Values(t *testing.T) {
	// Verify event type constants are sequential
	if PointerDown != 0 {
		t.Errorf("PointerDown = %d, want 0", PointerDown)
	}
	if PointerUp != 1 {
		t.Errorf("PointerUp = %d, want 1", PointerUp)
	}
	if PointerMove != 2 {
		t.Errorf("PointerMove = %d, want 2", PointerMove)
	}
	if PointerEnter != 3 {
		t.Errorf("PointerEnter = %d, want 3", PointerEnter)
	}
	if PointerLeave != 4 {
		t.Errorf("PointerLeave = %d, want 4", PointerLeave)
	}
	if PointerCancel != 5 {
		t.Errorf("PointerCancel = %d, want 5", PointerCancel)
	}
}

func TestPointerType_Values(t *testing.T) {
	// Verify pointer type constants are sequential
	if PointerTypeMouse != 0 {
		t.Errorf("PointerTypeMouse = %d, want 0", PointerTypeMouse)
	}
	if PointerTypeTouch != 1 {
		t.Errorf("PointerTypeTouch = %d, want 1", PointerTypeTouch)
	}
	if PointerTypePen != 2 {
		t.Errorf("PointerTypePen = %d, want 2", PointerTypePen)
	}
}

func TestNullPointerEventSource(t *testing.T) {
	// NullPointerEventSource should implement PointerEventSource
	var pes PointerEventSource = NullPointerEventSource{}

	// Method should be callable without panic
	pes.OnPointer(func(PointerEvent) {})
}

// mockPointerEventSource is used to verify PointerEventSource interface.
type mockPointerEventSource struct {
	handler func(PointerEvent)
}

func (m *mockPointerEventSource) OnPointer(fn func(PointerEvent)) {
	m.handler = fn
}

// Ensure mockPointerEventSource implements PointerEventSource.
var _ PointerEventSource = &mockPointerEventSource{}

func TestPointerEventSource_Interface(t *testing.T) {
	// Verify PointerEventSource can be used through the interface
	mock := &mockPointerEventSource{}
	var source PointerEventSource = mock

	var received *PointerEvent
	source.OnPointer(func(ev PointerEvent) {
		received = &ev
	})

	// Simulate event dispatch
	testEvent := PointerEvent{
		Type:        PointerDown,
		PointerID:   1,
		X:           100,
		Y:           200,
		PointerType: PointerTypeMouse,
		IsPrimary:   true,
		Button:      ButtonLeft,
		Buttons:     ButtonsLeft,
	}

	mock.handler(testEvent)

	if received == nil {
		t.Fatal("Handler was not called")
	}
	if received.Type != PointerDown {
		t.Errorf("received.Type = %v, want PointerDown", received.Type)
	}
	if received.PointerID != 1 {
		t.Errorf("received.PointerID = %d, want 1", received.PointerID)
	}
	if received.X != 100 {
		t.Errorf("received.X = %f, want 100", received.X)
	}
	if received.Y != 200 {
		t.Errorf("received.Y = %f, want 200", received.Y)
	}
}
