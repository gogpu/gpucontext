// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

package gpucontext

import "testing"

func TestNullWindowChrome_Defaults(t *testing.T) {
	var wc WindowChrome = NullWindowChrome{}

	if wc.IsFrameless() {
		t.Error("IsFrameless() should return false")
	}
	if wc.IsMaximized() {
		t.Error("IsMaximized() should return false")
	}
}

func TestNullWindowChrome_Actions(t *testing.T) {
	var wc WindowChrome = NullWindowChrome{}

	// All actions should succeed without panic
	wc.SetFrameless(true)
	wc.SetFrameless(false)
	wc.SetHitTestCallback(func(x, y float64) HitTestResult { return HitTestClient })
	wc.SetHitTestCallback(nil)
	wc.Minimize()
	wc.Maximize()
	wc.Close()
}

func TestHitTestResult_String(t *testing.T) {
	tests := []struct {
		result HitTestResult
		want   string
	}{
		{HitTestClient, "Client"},
		{HitTestCaption, "Caption"},
		{HitTestClose, "Close"},
		{HitTestMaximize, "Maximize"},
		{HitTestMinimize, "Minimize"},
		{HitTestResizeN, "ResizeN"},
		{HitTestResizeS, "ResizeS"},
		{HitTestResizeW, "ResizeW"},
		{HitTestResizeE, "ResizeE"},
		{HitTestResizeNW, "ResizeNW"},
		{HitTestResizeNE, "ResizeNE"},
		{HitTestResizeSW, "ResizeSW"},
		{HitTestResizeSE, "ResizeSE"},
		{HitTestResult(99), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.result.String(); got != tt.want {
				t.Errorf("HitTestResult(%d).String() = %q, want %q", tt.result, got, tt.want)
			}
		})
	}
}

func TestHitTestResult_Values(t *testing.T) {
	// Verify hit test result constants are sequential starting from 0
	if HitTestClient != 0 {
		t.Errorf("HitTestClient = %d, want 0", HitTestClient)
	}
	if HitTestCaption != 1 {
		t.Errorf("HitTestCaption = %d, want 1", HitTestCaption)
	}
	if HitTestClose != 2 {
		t.Errorf("HitTestClose = %d, want 2", HitTestClose)
	}
	if HitTestMaximize != 3 {
		t.Errorf("HitTestMaximize = %d, want 3", HitTestMaximize)
	}
	if HitTestMinimize != 4 {
		t.Errorf("HitTestMinimize = %d, want 4", HitTestMinimize)
	}
	if HitTestResizeN != 5 {
		t.Errorf("HitTestResizeN = %d, want 5", HitTestResizeN)
	}
	if HitTestResizeS != 6 {
		t.Errorf("HitTestResizeS = %d, want 6", HitTestResizeS)
	}
	if HitTestResizeW != 7 {
		t.Errorf("HitTestResizeW = %d, want 7", HitTestResizeW)
	}
	if HitTestResizeE != 8 {
		t.Errorf("HitTestResizeE = %d, want 8", HitTestResizeE)
	}
	if HitTestResizeNW != 9 {
		t.Errorf("HitTestResizeNW = %d, want 9", HitTestResizeNW)
	}
	if HitTestResizeNE != 10 {
		t.Errorf("HitTestResizeNE = %d, want 10", HitTestResizeNE)
	}
	if HitTestResizeSW != 11 {
		t.Errorf("HitTestResizeSW = %d, want 11", HitTestResizeSW)
	}
	if HitTestResizeSE != 12 {
		t.Errorf("HitTestResizeSE = %d, want 12", HitTestResizeSE)
	}
}

// mockWindowChrome verifies the interface can be implemented by custom types.
type mockWindowChrome struct {
	frameless bool
	maximized bool
	closed    bool
	minimized bool
	callback  HitTestCallback
}

func (m *mockWindowChrome) SetFrameless(frameless bool)           { m.frameless = frameless }
func (m *mockWindowChrome) IsFrameless() bool                     { return m.frameless }
func (m *mockWindowChrome) SetHitTestCallback(cb HitTestCallback) { m.callback = cb }
func (m *mockWindowChrome) Minimize()                             { m.minimized = true }
func (m *mockWindowChrome) Maximize()                             { m.maximized = !m.maximized }
func (m *mockWindowChrome) IsMaximized() bool                     { return m.maximized }
func (m *mockWindowChrome) Close()                                { m.closed = true }

// Ensure mockWindowChrome implements WindowChrome.
var _ WindowChrome = &mockWindowChrome{}

func TestWindowChrome_CustomImplementation(t *testing.T) {
	mock := &mockWindowChrome{}
	var wc WindowChrome = mock

	// Test frameless mode
	wc.SetFrameless(true)
	if !wc.IsFrameless() {
		t.Error("IsFrameless() should return true after SetFrameless(true)")
	}
	wc.SetFrameless(false)
	if wc.IsFrameless() {
		t.Error("IsFrameless() should return false after SetFrameless(false)")
	}

	// Test hit test callback
	called := false
	wc.SetHitTestCallback(func(x, y float64) HitTestResult {
		called = true
		if x > 100 {
			return HitTestClient
		}
		return HitTestCaption
	})
	if mock.callback == nil {
		t.Fatal("callback should be set")
	}
	result := mock.callback(50, 10)
	if !called {
		t.Error("callback should have been called")
	}
	if result != HitTestCaption {
		t.Errorf("callback(50, 10) = %v, want Caption", result)
	}
	result = mock.callback(150, 10)
	if result != HitTestClient {
		t.Errorf("callback(150, 10) = %v, want Client", result)
	}

	// Test minimize
	wc.Minimize()
	if !mock.minimized {
		t.Error("Minimize() should set minimized to true")
	}

	// Test maximize toggle
	wc.Maximize()
	if !wc.IsMaximized() {
		t.Error("IsMaximized() should return true after first Maximize()")
	}
	wc.Maximize()
	if wc.IsMaximized() {
		t.Error("IsMaximized() should return false after second Maximize()")
	}

	// Test close
	wc.Close()
	if !mock.closed {
		t.Error("Close() should set closed to true")
	}
}
