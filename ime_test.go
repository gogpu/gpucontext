// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

package gpucontext

import "testing"

func TestIMETextRangeUsesUTF8ByteOffsets(t *testing.T) {
	text := "A你B"
	// A = 1 byte, 你 = 3 bytes, B = 1 byte.
	cases := []struct {
		name      string
		textRange IMETextRange
		valid     bool
	}{
		{"empty at start", IMETextRange{Start: 0, End: 0}, true},
		{"multibyte rune", IMETextRange{Start: 1, End: 4}, true},
		{"end of string", IMETextRange{Start: 5, End: 5}, true},
		{"inside rune", IMETextRange{Start: 2, End: 4}, false},
		{"end inside rune", IMETextRange{Start: 1, End: 3}, false},
		{"reversed", IMETextRange{Start: 4, End: 1}, false},
		{"past end", IMETextRange{Start: 0, End: 6}, false},
		{"negative", IMETextRange{Start: -1, End: 0}, false},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.textRange.IsValid(text); got != tt.valid {
				t.Fatalf("IMETextRange%+v.IsValid(%q) = %v, want %v", tt.textRange, text, got, tt.valid)
			}
		})
	}
	if (IMETextRange{Start: 0, End: 1}).IsValid("\xff") {
		t.Fatal("a range in invalid UTF-8 must be invalid")
	}

	if !(IMETextRange{Start: 1, End: 1}).Empty() {
		t.Fatal("collapsed range should be empty")
	}
	if (IMETextRange{Start: 1, End: 4}).Empty() {
		t.Fatal("non-collapsed range should not be empty")
	}
}

func TestIMEStateRanges(t *testing.T) {
	state := IMEState{
		Composing:       true,
		CompositionText: "にほ",
		Cursor:          IMETextRange{Start: 3, End: 6},
		Selection:       IMETextRange{Start: 0, End: 3},
	}
	if !state.HasCursor() || !state.IsValid() {
		t.Fatal("state with UTF-8 boundary ranges should be valid with a visible cursor")
	}

	state.Cursor.End = 2 // inside the first three-byte rune, and before Start.
	if state.IsValid() {
		t.Fatal("state with a non-boundary cursor must be invalid")
	}

	state.Cursor = IMETextRange{Start: -1, End: -1}
	if state.HasCursor() || !state.IsValid() {
		t.Fatal("a {-1, -1} cursor should be a valid hidden cursor")
	}

	state.Selection = IMETextRange{Start: 1, End: 3}
	if state.IsValid() {
		t.Fatal("state with a non-boundary selection must be invalid")
	}

	if !(IMEState{}).IsValid() {
		t.Fatal("the zero IMEState (no composition) must be valid")
	}
}

func TestIMESurroundingTextSelection(t *testing.T) {
	surrounding := IMESurroundingText{Text: "ab你cd", Cursor: 6, Anchor: 1}
	if !surrounding.IsValid() {
		t.Fatal("surrounding text endpoints should be valid UTF-8 offsets")
	}
	if got := surrounding.SelectionRange(); got != (IMETextRange{Start: 1, End: 6}) {
		t.Fatalf("SelectionRange() = %+v", got)
	}
	surrounding.Cursor, surrounding.Anchor = 1, 6
	if got := surrounding.SelectionRange(); got != (IMETextRange{Start: 1, End: 6}) {
		t.Fatalf("forward SelectionRange() = %+v", got)
	}

	surrounding.Cursor = 3 // inside 你.
	if surrounding.IsValid() {
		t.Fatal("surrounding text with a non-boundary cursor must be invalid")
	}
}

func TestContentPurposeString(t *testing.T) {
	tests := []struct {
		purpose ContentPurpose
		want    string
	}{
		{ContentPurposeNormal, "Normal"},
		{ContentPurposeAlpha, "Alpha"},
		{ContentPurposeDigits, "Digits"},
		{ContentPurposeNumber, "Number"},
		{ContentPurposePhone, "Phone"},
		{ContentPurposeURL, "URL"},
		{ContentPurposeEmail, "Email"},
		{ContentPurposeName, "Name"},
		{ContentPurposePassword, "Password"},
		{ContentPurposePin, "Pin"},
		{ContentPurposeDate, "Date"},
		{ContentPurposeTime, "Time"},
		{ContentPurposeDateTime, "DateTime"},
		{ContentPurposeTerminal, "Terminal"},
		{ContentPurposeChat, "Chat"},
		{ContentPurpose(255), "Unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.purpose.String(); got != tt.want {
				t.Fatalf("ContentPurpose(%d).String() = %q, want %q", tt.purpose, got, tt.want)
			}
		})
	}
}

// A purpose added to the iota block without a name in contentPurposeNames
// fails here.
func TestContentPurposeStringCoversEveryValue(t *testing.T) {
	for p := ContentPurposeNormal; p < contentPurposeEnd; p++ {
		if p.String() == stringUnknown {
			t.Errorf("ContentPurpose(%d) has no name", p)
		}
	}
	if got, want := len(contentPurposeNames), int(contentPurposeEnd); got != want {
		t.Fatalf("contentPurposeNames has %d entries, want %d", got, want)
	}
}

func TestContentHintFlags(t *testing.T) {
	hints := ContentHintCompletion | ContentHintSensitiveData | ContentHintMultiline
	for _, hint := range []ContentHint{ContentHintCompletion, ContentHintSensitiveData, ContentHintMultiline} {
		if !hints.Has(hint) {
			t.Errorf("%v should be present in %v", hint, hints)
		}
	}
	if hints.Has(ContentHintLatin) {
		t.Error("Latin should not be present")
	}
	if got := hints.String(); got != "Completion|SensitiveData|Multiline" {
		t.Fatalf("ContentHint.String() = %q", got)
	}
	if got := ContentHintNone.String(); got != stringNone {
		t.Fatalf("ContentHintNone.String() = %q, want %q", got, stringNone)
	}
	if ContentHintNone.Has(ContentHintNone) {
		t.Fatal("ContentHintNone must not report a hint bit")
	}
	if got := ContentHint(1 << 15).String(); got != "0x8000" {
		t.Fatalf("unknown ContentHint.String() = %q, want 0x8000", got)
	}
	if got := (ContentHintLatin | ContentHint(1<<14)).String(); got != "Latin|0x4000" {
		t.Fatalf("mixed ContentHint.String() = %q, want Latin|0x4000", got)
	}
}

// Every hint flag in the iota block has a distinct single bit and a name.
func TestContentHintNamesCoverEveryFlag(t *testing.T) {
	var all ContentHint
	for _, item := range contentHintNames {
		if item.hint == 0 || item.hint&(item.hint-1) != 0 {
			t.Errorf("%s is not a single bit: %#x", item.name, item.hint)
		}
		if all&item.hint != 0 {
			t.Errorf("%s reuses a bit", item.name)
		}
		all |= item.hint
	}
	if want := contentHintEnd - 1; all != want {
		t.Fatalf("named hints cover %#x, want every flag in the block (%#x)", uint64(all), uint64(want))
	}
}

func TestIMEDeleteSurroundingValidation(t *testing.T) {
	if !(IMEDeleteSurroundingEvent{Before: 2, After: 3}).IsValid() {
		t.Fatal("non-negative deletion lengths should be valid")
	}
	if (IMEDeleteSurroundingEvent{Before: -1}).IsValid() {
		t.Fatal("negative deletion length should be invalid")
	}
	if (IMEDeleteSurroundingEvent{After: -1}).IsValid() {
		t.Fatal("negative deletion length should be invalid")
	}
}

func TestIMECapabilitiesHas(t *testing.T) {
	caps := IMECapabilityCursorArea | IMECapabilitySurroundingText
	if !caps.Has(IMECapabilityCursorArea) || !caps.Has(IMECapabilitySurroundingText) {
		t.Fatal("set capabilities must be reported")
	}
	if caps.Has(IMECapabilityContentType) {
		t.Fatal("ContentType is not set")
	}
	if caps.Has(IMECapabilityCursorArea | IMECapabilityContentType) {
		t.Fatal("Has must require every requested bit")
	}
	if caps.Has(0) {
		t.Fatal("Has(0) must be false")
	}
}

// fullIMEHost implements the whole contract: IMEController plus the optional
// IMEEventSource.
type fullIMEHost struct{ caps IMECapabilities }

func (fullIMEHost) SetIMEEnabled(bool)                                     {}
func (fullIMEHost) SetIMECursorArea(IMECursorArea)                         {}
func (fullIMEHost) SetIMEContentType(ContentPurpose, ContentHint)          {}
func (fullIMEHost) SetIMESurroundingText(IMESurroundingText)               {}
func (fullIMEHost) CancelIME()                                             {}
func (h fullIMEHost) IMECapabilities() IMECapabilities                     { return h.caps }
func (fullIMEHost) OnIMECanceled(func())                                   {}
func (fullIMEHost) OnIMEDisabled(func())                                   {}
func (fullIMEHost) OnIMEDeleteSurrounding(func(IMEDeleteSurroundingEvent)) {}

var (
	_ IMEController  = fullIMEHost{}
	_ IMEEventSource = fullIMEHost{}
)

func TestIMEEventSourceIsFoundByTypeAssertion(t *testing.T) {
	var controller IMEController = fullIMEHost{caps: IMECapabilitySurroundingText}
	if !controller.IMECapabilities().Has(IMECapabilitySurroundingText) {
		t.Fatal("capabilities come from the controller itself")
	}
	if _, ok := controller.(IMEEventSource); !ok {
		t.Fatal("a host with IME events must be found by a plain type assertion")
	}
	var plain any = struct{}{}
	if _, ok := plain.(IMEEventSource); ok {
		t.Fatal("a value without the methods must not satisfy IMEEventSource")
	}
}
