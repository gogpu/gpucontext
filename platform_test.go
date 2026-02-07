// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

package gpucontext

import "testing"

func TestNullPlatformProvider_ClipboardRead(t *testing.T) {
	var pp PlatformProvider = NullPlatformProvider{}

	text, err := pp.ClipboardRead()
	if text != "" {
		t.Errorf("ClipboardRead() text = %q, want empty", text)
	}
	if err != nil {
		t.Errorf("ClipboardRead() err = %v, want nil", err)
	}
}

func TestNullPlatformProvider_ClipboardWrite(t *testing.T) {
	var pp PlatformProvider = NullPlatformProvider{}

	err := pp.ClipboardWrite("hello")
	if err != nil {
		t.Errorf("ClipboardWrite() err = %v, want nil", err)
	}
}

func TestNullPlatformProvider_SetCursor(t *testing.T) {
	var pp PlatformProvider = NullPlatformProvider{}

	// All cursor shapes should be accepted without panic
	cursors := []CursorShape{
		CursorDefault, CursorPointer, CursorText, CursorCrosshair,
		CursorMove, CursorResizeNS, CursorResizeEW, CursorResizeNWSE,
		CursorResizeNESW, CursorNotAllowed, CursorWait, CursorNone,
	}
	for _, c := range cursors {
		pp.SetCursor(c)
	}
}

func TestNullPlatformProvider_Defaults(t *testing.T) {
	var pp PlatformProvider = NullPlatformProvider{}

	if pp.DarkMode() {
		t.Error("DarkMode() should return false")
	}
	if pp.ReduceMotion() {
		t.Error("ReduceMotion() should return false")
	}
	if pp.HighContrast() {
		t.Error("HighContrast() should return false")
	}
	if got := pp.FontScale(); got != 1.0 {
		t.Errorf("FontScale() = %f, want 1.0", got)
	}
}

func TestCursorShape_String(t *testing.T) {
	tests := []struct {
		cursor CursorShape
		want   string
	}{
		{CursorDefault, "Default"},
		{CursorPointer, "Pointer"},
		{CursorText, "Text"},
		{CursorCrosshair, "Crosshair"},
		{CursorMove, "Move"},
		{CursorResizeNS, "ResizeNS"},
		{CursorResizeEW, "ResizeEW"},
		{CursorResizeNWSE, "ResizeNWSE"},
		{CursorResizeNESW, "ResizeNESW"},
		{CursorNotAllowed, "NotAllowed"},
		{CursorWait, "Wait"},
		{CursorNone, "None"},
		{CursorShape(99), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.cursor.String(); got != tt.want {
				t.Errorf("CursorShape(%d).String() = %q, want %q", tt.cursor, got, tt.want)
			}
		})
	}
}

func TestCursorShape_Values(t *testing.T) {
	// Verify cursor shape constants are sequential starting from 0
	if CursorDefault != 0 {
		t.Errorf("CursorDefault = %d, want 0", CursorDefault)
	}
	if CursorPointer != 1 {
		t.Errorf("CursorPointer = %d, want 1", CursorPointer)
	}
	if CursorText != 2 {
		t.Errorf("CursorText = %d, want 2", CursorText)
	}
	if CursorCrosshair != 3 {
		t.Errorf("CursorCrosshair = %d, want 3", CursorCrosshair)
	}
	if CursorMove != 4 {
		t.Errorf("CursorMove = %d, want 4", CursorMove)
	}
	if CursorResizeNS != 5 {
		t.Errorf("CursorResizeNS = %d, want 5", CursorResizeNS)
	}
	if CursorResizeEW != 6 {
		t.Errorf("CursorResizeEW = %d, want 6", CursorResizeEW)
	}
	if CursorResizeNWSE != 7 {
		t.Errorf("CursorResizeNWSE = %d, want 7", CursorResizeNWSE)
	}
	if CursorResizeNESW != 8 {
		t.Errorf("CursorResizeNESW = %d, want 8", CursorResizeNESW)
	}
	if CursorNotAllowed != 9 {
		t.Errorf("CursorNotAllowed = %d, want 9", CursorNotAllowed)
	}
	if CursorWait != 10 {
		t.Errorf("CursorWait = %d, want 10", CursorWait)
	}
	if CursorNone != 11 {
		t.Errorf("CursorNone = %d, want 11", CursorNone)
	}
}

// mockPlatformProvider verifies the interface can be implemented by custom types.
type mockPlatformProvider struct {
	clipboard    string
	cursor       CursorShape
	darkMode     bool
	reduceMotion bool
	highContrast bool
	fontScale    float32
}

func (m *mockPlatformProvider) ClipboardRead() (string, error) { return m.clipboard, nil }
func (m *mockPlatformProvider) ClipboardWrite(text string) error {
	m.clipboard = text
	return nil
}
func (m *mockPlatformProvider) SetCursor(cursor CursorShape) { m.cursor = cursor }
func (m *mockPlatformProvider) DarkMode() bool               { return m.darkMode }
func (m *mockPlatformProvider) ReduceMotion() bool           { return m.reduceMotion }
func (m *mockPlatformProvider) HighContrast() bool           { return m.highContrast }
func (m *mockPlatformProvider) FontScale() float32           { return m.fontScale }

// Ensure mockPlatformProvider implements PlatformProvider.
var _ PlatformProvider = &mockPlatformProvider{}

func TestPlatformProvider_CustomImplementation(t *testing.T) {
	mock := &mockPlatformProvider{
		darkMode:     true,
		reduceMotion: true,
		highContrast: true,
		fontScale:    1.5,
	}
	var pp PlatformProvider = mock

	// Test clipboard round-trip
	err := pp.ClipboardWrite("test clipboard")
	if err != nil {
		t.Fatalf("ClipboardWrite() err = %v", err)
	}
	text, err := pp.ClipboardRead()
	if err != nil {
		t.Fatalf("ClipboardRead() err = %v", err)
	}
	if text != "test clipboard" {
		t.Errorf("ClipboardRead() = %q, want \"test clipboard\"", text)
	}

	// Test cursor
	pp.SetCursor(CursorPointer)
	if mock.cursor != CursorPointer {
		t.Errorf("cursor = %v, want CursorPointer", mock.cursor)
	}

	// Test system preferences
	if !pp.DarkMode() {
		t.Error("DarkMode() should return true")
	}
	if !pp.ReduceMotion() {
		t.Error("ReduceMotion() should return true")
	}
	if !pp.HighContrast() {
		t.Error("HighContrast() should return true")
	}
	if got := pp.FontScale(); got != 1.5 {
		t.Errorf("FontScale() = %f, want 1.5", got)
	}
}
