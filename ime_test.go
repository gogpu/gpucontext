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
		{"reversed", IMETextRange{Start: 4, End: 1}, false},
		{"past end", IMETextRange{Start: 0, End: 6}, false},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.textRange.IsValid(text); got != tt.valid {
				t.Fatalf("IMETextRange%+v.IsValid(%q) = %v, want %v", tt.textRange, text, got, tt.valid)
			}
		})
	}

	if !(IMETextRange{Start: 1, End: 1}).Empty() {
		t.Fatal("collapsed range should be empty")
	}
	if (IMETextRange{Start: 1, End: 4}).Empty() {
		t.Fatal("non-collapsed range should not be empty")
	}
}

func TestIMECompositionRanges(t *testing.T) {
	composition := IMEComposition{
		CompositionText: "にほ",
		CursorBegin:     3,
		CursorEnd:       6,
		SelectionStart:  0,
		SelectionEnd:    3,
	}
	if !composition.IsValid() {
		t.Fatal("composition with UTF-8 boundary ranges should be valid")
	}
	if got := composition.CursorRange(); got != (IMETextRange{Start: 3, End: 6}) {
		t.Fatalf("CursorRange() = %+v", got)
	}
	if got := composition.SelectionRange(); got != (IMETextRange{Start: 0, End: 3}) {
		t.Fatalf("SelectionRange() = %+v", got)
	}

	composition.CursorEnd = 2 // inside the first three-byte rune.
	if composition.IsValid() {
		t.Fatal("composition with a non-boundary cursor must be invalid")
	}
	composition.CursorBegin, composition.CursorEnd = -1, -1
	if composition.HasCursor() || !composition.IsValid() {
		t.Fatal("a -1,-1 cursor range should represent a valid hidden cursor")
	}
}

func TestLegacyIMEStateCompositionConversion(t *testing.T) {
	legacy := IMEState{CompositionText: "ni", CursorPos: 2}
	got := legacy.Composition()
	if got.CursorRange() != (IMETextRange{Start: 2, End: 2}) {
		t.Fatalf("legacy collapsed cursor conversion = %+v", got.CursorRange())
	}

	ranged := IMEState{CompositionText: "に", CursorBegin: 0, CursorEnd: 3}
	if got := ranged.Composition().CursorRange(); got != (IMETextRange{Start: 0, End: 3}) {
		t.Fatalf("versioned cursor conversion = %+v", got)
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
	if ContentHintAutoCapitalize != ContentHintAutoCapitalization {
		t.Fatal("browser spelling must remain an alias")
	}
	if got := ContentHintNone.String(); got != stringNone {
		t.Fatalf("ContentHintNone.String() = %q, want %q", got, stringNone)
	}
	if ContentHintNone.Has(ContentHintNone) {
		t.Fatal("ContentHintNone must not report a capability bit")
	}
	if got := ContentHint(1 << 15).String(); got != "Unknown" {
		t.Fatalf("unknown ContentHint.String() = %q, want Unknown", got)
	}
}

func TestIMEDeleteSurroundingValidation(t *testing.T) {
	if !(IMEDeleteSurroundingEvent{Before: 2, After: 3}).IsValid() {
		t.Fatal("non-negative deletion lengths should be valid")
	}
	if (IMEDeleteSurroundingEvent{Before: -1}).IsValid() {
		t.Fatal("negative deletion length should be invalid")
	}
}

// legacyIMEController intentionally implements only the original interface.
// Adding the versioned extension must not make this implementation invalid.
type legacyIMEController struct{}

func (legacyIMEController) SetIMEPosition(_, _ int) {}
func (legacyIMEController) SetIMEEnabled(bool)      {}

var _ IMEController = legacyIMEController{}

type richIMEProvider struct{}

func (richIMEProvider) SetIMEPosition(_, _ int)                                {}
func (richIMEProvider) SetIMEEnabled(bool)                                     {}
func (richIMEProvider) SetIMECursorArea(IMECursorArea)                         {}
func (richIMEProvider) SetIMEContentType(ContentPurpose, ContentHint)          {}
func (richIMEProvider) SetIMESurroundingText(IMESurroundingText)               {}
func (richIMEProvider) CancelIME()                                             {}
func (richIMEProvider) OnIMECompositionUpdateV2(func(IMEComposition))          {}
func (richIMEProvider) OnIMECanceled(func())                                   {}
func (richIMEProvider) OnIMEDisabled(func())                                   {}
func (richIMEProvider) OnIMEDeleteSurrounding(func(IMEDeleteSurroundingEvent)) {}
func (richIMEProvider) IMECapabilities() IMECapabilities {
	return IMECapabilities{
		Version: IMEContractVersion,
		Features: IMECapabilityComposition |
			IMECapabilityCommit |
			IMECapabilityCancel |
			IMECapabilityDisabled |
			IMECapabilityDeleteSurrounding |
			IMECapabilityCursorArea |
			IMECapabilitySurroundingText |
			IMECapabilityContentPurpose |
			IMECapabilityContentHints,
	}
}

var _ IMEControllerV2 = richIMEProvider{}
var _ IMEEventSourceV2 = richIMEProvider{}
var _ IMECapabilityProviderV2 = richIMEProvider{}
var _ IMEContractV2 = richIMEProvider{}

func TestDiscoverIMECapabilities(t *testing.T) {
	caps, ok := DiscoverIMECapabilities(richIMEProvider{})
	if !ok {
		t.Fatal("rich provider should advertise capabilities")
	}
	if caps.Version != IMEContractVersion {
		t.Fatalf("capability version = %d, want %d", caps.Version, IMEContractVersion)
	}
	if !caps.Supports(IMECapabilityCursorArea | IMECapabilityContentHints) {
		t.Fatal("provider should advertise cursor-area and content-hint support")
	}
	if _, ok := DiscoverIMECapabilities(legacyIMEController{}); ok {
		t.Fatal("legacy provider must not be reported as versioned")
	}
}

func TestZeroCapabilitiesAreNotAdvertised(t *testing.T) {
	provider := zeroCapabilitiesProvider{}
	if _, ok := DiscoverIMECapabilities(provider); ok {
		t.Fatal("version zero must not be advertised")
	}
}

type zeroCapabilitiesProvider struct{}

func (zeroCapabilitiesProvider) IMECapabilities() IMECapabilities { return IMECapabilities{} }
