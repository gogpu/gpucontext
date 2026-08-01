// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

package gpucontext

import "testing"

func TestInputEventInterface(t *testing.T) {
	events := []InputEvent{
		KeyEvent{Key: KeyA, Mods: ModShift, Pressed: true},
		KeyEvent{Key: KeyEscape, Pressed: false},
		CharEvent{Char: 'A'},
		FocusEvent{Focused: true},
		FocusEvent{Focused: false},
		ResizeEvent{Width: 800, Height: 600},
		ScaleChangedEvent{ScaleFactor: 2.0, Width: 800, Height: 600},
		PointerEvent{Type: PointerDown, X: 100, Y: 200},
		PointerEvent{Type: PointerMove, X: 150, Y: 250},
		ScrollEvent{DeltaY: -3.0, DeltaMode: ScrollDeltaLine},
	}

	if len(events) != 10 {
		t.Fatalf("expected 10 events, got %d", len(events))
	}

	for i, ev := range events {
		if ev == nil {
			t.Errorf("event %d is nil", i)
		}
	}
}

func TestInputEventTypeSwitch(t *testing.T) {
	events := []InputEvent{
		KeyEvent{Key: KeyA, Pressed: true},
		PointerEvent{Type: PointerDown, X: 10, Y: 20},
		ScrollEvent{DeltaY: 1.0},
		CharEvent{Char: 'x'},
		FocusEvent{Focused: true},
		ResizeEvent{Width: 1024, Height: 768},
		ScaleChangedEvent{ScaleFactor: 2.0, Width: 800, Height: 600},
	}

	counts := map[string]int{}
	for _, ev := range events {
		switch ev.(type) {
		case KeyEvent:
			counts["key"]++
		case PointerEvent:
			counts["pointer"]++
		case ScrollEvent:
			counts["scroll"]++
		case CharEvent:
			counts["char"]++
		case FocusEvent:
			counts["focus"]++
		case ResizeEvent:
			counts["resize"]++
		case ScaleChangedEvent:
			counts["scale"]++
		default:
			t.Errorf("unexpected event type: %T", ev)
		}
	}

	for _, name := range []string{"key", "pointer", "scroll", "char", "focus", "resize", "scale"} {
		if counts[name] != 1 {
			t.Errorf("expected 1 %s event, got %d", name, counts[name])
		}
	}
}

func TestKeyEventFields(t *testing.T) {
	ev := KeyEvent{Key: KeyA, Mods: ModShift | ModControl, Pressed: true}
	if ev.Key != KeyA {
		t.Errorf("Key = %v, want KeyA", ev.Key)
	}
	if !ev.Mods.HasShift() {
		t.Error("expected HasShift")
	}
	if !ev.Mods.HasControl() {
		t.Error("expected HasControl")
	}
	if !ev.Pressed {
		t.Error("expected Pressed = true")
	}
}

func TestCharEventFields(t *testing.T) {
	ev := CharEvent{Char: '日'}
	if ev.Char != '日' {
		t.Errorf("Char = %q, want '日'", ev.Char)
	}
}

func TestFocusEventFields(t *testing.T) {
	gained := FocusEvent{Focused: true}
	lost := FocusEvent{Focused: false}
	if !gained.Focused {
		t.Error("expected Focused = true")
	}
	if lost.Focused {
		t.Error("expected Focused = false")
	}
}

func TestResizeEventFields(t *testing.T) {
	ev := ResizeEvent{Width: 1920, Height: 1080}
	if ev.Width != 1920 || ev.Height != 1080 {
		t.Errorf("got %dx%d, want 1920x1080", ev.Width, ev.Height)
	}
}

func TestScaleChangedEventFields(t *testing.T) {
	ev := ScaleChangedEvent{ScaleFactor: 2.0, Width: 800, Height: 600}
	if ev.ScaleFactor != 2.0 {
		t.Errorf("ScaleFactor = %v, want 2.0", ev.ScaleFactor)
	}
	if ev.Width != 800 || ev.Height != 600 {
		t.Errorf("got %dx%d, want 800x600", ev.Width, ev.Height)
	}
}

func TestScaleChangedEventImplementsInputEvent(t *testing.T) {
	var ev InputEvent = ScaleChangedEvent{ScaleFactor: 1.5, Width: 1024, Height: 768}
	se, ok := ev.(ScaleChangedEvent)
	if !ok {
		t.Fatal("ScaleChangedEvent does not implement InputEvent")
	}
	if se.ScaleFactor != 1.5 {
		t.Errorf("ScaleFactor = %v, want 1.5", se.ScaleFactor)
	}
}

func TestPointerEventImplementsInputEvent(t *testing.T) {
	var ev InputEvent = PointerEvent{Type: PointerDown, X: 50, Y: 100}
	pe, ok := ev.(PointerEvent)
	if !ok {
		t.Fatal("PointerEvent does not implement InputEvent")
	}
	if pe.Type != PointerDown {
		t.Errorf("Type = %v, want PointerDown", pe.Type)
	}
	if pe.X != 50 || pe.Y != 100 {
		t.Errorf("coordinates = (%v, %v), want (50, 100)", pe.X, pe.Y)
	}
}

func TestScrollEventImplementsInputEvent(t *testing.T) {
	var ev InputEvent = ScrollEvent{DeltaY: -1.0, Phase: ScrollPhaseBegan}
	se, ok := ev.(ScrollEvent)
	if !ok {
		t.Fatal("ScrollEvent does not implement InputEvent")
	}
	if se.DeltaY != -1.0 {
		t.Errorf("DeltaY = %v, want -1.0", se.DeltaY)
	}
	if se.Phase != ScrollPhaseBegan {
		t.Errorf("Phase = %v, want ScrollPhaseBegan", se.Phase)
	}
}
