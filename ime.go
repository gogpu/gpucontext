// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

package gpucontext

import (
	"strings"
	"unicode/utf8"
)

// IMEContractVersion is the version of the optional, richer IME contract.
//
// The original EventSource IME callbacks and IMEController interface remain
// unchanged. A host advertises this version by implementing
// IMECapabilityProviderV2; consumers must still check the individual feature
// bits before using an optional operation.
const IMEContractVersion uint = 2

// IMETextRange is a half-open range [Start, End) in a UTF-8 string.
//
// Start and End are byte offsets, not rune or UTF-16 indexes. This is the
// representation used by the Wayland text-input protocol and lets a host
// forward a Go string without changing its encoding. Values must land on
// UTF-8 rune boundaries. Native backends are responsible for converting their
// platform's range units at the boundary.
type IMETextRange struct {
	// Start is the first byte included in the range.
	Start int

	// End is the first byte excluded from the range.
	End int
}

// Empty reports whether the range contains no bytes.
func (r IMETextRange) Empty() bool {
	return r.Start == r.End
}

// IsValid reports whether the range is within text and starts and ends on
// UTF-8 rune boundaries. Invalid UTF-8 text is rejected as well because all
// IME payloads use UTF-8 offsets.
func (r IMETextRange) IsValid(text string) bool {
	if !utf8.ValidString(text) || r.Start < 0 || r.End < r.Start || r.End > len(text) {
		return false
	}
	if r.Start < len(text) && !utf8.RuneStart(text[r.Start]) {
		return false
	}
	return r.End == len(text) || utf8.RuneStart(text[r.End])
}

// IMEComposition is the full-fidelity composition payload delivered through
// IMEEventSourceV2.OnIMECompositionUpdateV2. All positions refer to
// CompositionText and use UTF-8 byte offsets. The original
// EventSource.OnIMECompositionUpdate callback remains available for legacy
// consumers.
//
// CursorBegin and CursorEnd describe the active cursor range. Most IMEs report
// a collapsed range, but a range is required for platforms that expose a
// selected segment inside the preedit. SelectionStart and SelectionEnd are the
// target/marked range supplied by the IME, when one is available. A pair of
// -1 values in CursorBegin/CursorEnd means that the IME requests a hidden
// cursor; all other positions must be valid UTF-8 boundaries.
type IMEComposition struct {
	// CompositionText is the current preedit text.
	CompositionText string

	// CursorBegin is the first byte of the active cursor range.
	CursorBegin int

	// CursorEnd is the first byte after the active cursor range.
	CursorEnd int

	// SelectionStart is the first byte of the marked/selected range.
	SelectionStart int

	// SelectionEnd is the first byte after the marked/selected range.
	SelectionEnd int
}

// Composition converts the legacy IMEState payload to the versioned
// full-fidelity shape. Providers that only populate the legacy collapsed
// CursorPos field are represented as a collapsed cursor range; providers that
// populate CursorBegin/End retain the richer range.
func (s IMEState) Composition() IMEComposition {
	begin, end := s.CursorBegin, s.CursorEnd
	if begin == 0 && end == 0 && s.CursorPos != 0 {
		begin, end = s.CursorPos, s.CursorPos
	}
	return IMEComposition{
		CompositionText: s.CompositionText,
		CursorBegin:     begin,
		CursorEnd:       end,
		SelectionStart:  s.SelectionStart,
		SelectionEnd:    s.SelectionEnd,
	}
}

// CursorRange returns the active cursor range in CompositionText.
func (s IMEComposition) CursorRange() IMETextRange {
	return IMETextRange{Start: s.CursorBegin, End: s.CursorEnd}
}

// HasCursor reports whether the IME supplied a visible cursor range.
func (s IMEComposition) HasCursor() bool {
	return s.CursorBegin >= 0 && s.CursorEnd >= 0
}

// SelectionRange returns the marked/selected range in CompositionText.
func (s IMEComposition) SelectionRange() IMETextRange {
	return IMETextRange{Start: s.SelectionStart, End: s.SelectionEnd}
}

// IsValid reports whether both ranges are valid UTF-8 ranges in the preedit.
func (s IMEComposition) IsValid() bool {
	cursorValid := (!s.HasCursor() && s.CursorBegin == -1 && s.CursorEnd == -1) ||
		s.CursorRange().IsValid(s.CompositionText)
	return cursorValid &&
		s.SelectionRange().IsValid(s.CompositionText)
}

// IMESurroundingText is the text around the insertion point supplied to the
// platform IME. Cursor and Anchor are UTF-8 byte offsets into Text. They may
// differ when the surrounding text contains a selection; their order is
// preserved because some native IMEs use the direction of the selection.
type IMESurroundingText struct {
	// Text is the UTF-8 text visible to the IME around the insertion point.
	Text string

	// Cursor is the insertion point byte offset in Text.
	Cursor int

	// Anchor is the other endpoint of the selection byte offset in Text.
	Anchor int
}

// SelectionRange returns the normalized selection range represented by Cursor
// and Anchor.
func (s IMESurroundingText) SelectionRange() IMETextRange {
	if s.Cursor <= s.Anchor {
		return IMETextRange{Start: s.Cursor, End: s.Anchor}
	}
	return IMETextRange{Start: s.Anchor, End: s.Cursor}
}

// IsValid reports whether Text is valid UTF-8 and both endpoints are on rune
// boundaries within it.
func (s IMESurroundingText) IsValid() bool {
	return IMETextRange{Start: s.Cursor, End: s.Cursor}.IsValid(s.Text) &&
		IMETextRange{Start: s.Anchor, End: s.Anchor}.IsValid(s.Text)
}

// IMECursorArea is the caret rectangle used to place an IME candidate window.
// Coordinates are logical DIP relative to the window content area, matching
// WindowProvider and the EventSource pointer callbacks. Width and Height may
// be zero when only a caret point is available.
type IMECursorArea struct {
	// X is the left edge in logical DIP.
	X float64

	// Y is the top edge in logical DIP.
	Y float64

	// Width is the rectangle width in logical DIP.
	Width float64

	// Height is the rectangle height in logical DIP.
	Height float64
}

// ContentPurpose describes the kind of text expected by an IME. The values
// intentionally match the cross-platform purpose set used by winit and the
// Wayland text-input protocol. Unknown values must be treated as Normal by a
// host.
type ContentPurpose uint8

const (
	// ContentPurposeNormal is unrestricted text.
	ContentPurposeNormal ContentPurpose = iota
	// ContentPurposeAlpha accepts alphabetic characters.
	ContentPurposeAlpha
	// ContentPurposeDigits accepts decimal digits.
	ContentPurposeDigits
	// ContentPurposeNumber accepts a numeric value (including punctuation).
	ContentPurposeNumber
	// ContentPurposePhone accepts a telephone number.
	ContentPurposePhone
	// ContentPurposeURL accepts a URL.
	ContentPurposeURL
	// ContentPurposeEmail accepts an email address.
	ContentPurposeEmail
	// ContentPurposeName accepts a person's name.
	ContentPurposeName
	// ContentPurposePassword accepts a password.
	ContentPurposePassword
	// ContentPurposePin accepts a short numeric PIN.
	ContentPurposePin
	// ContentPurposeDate accepts a calendar date.
	ContentPurposeDate
	// ContentPurposeTime accepts a time of day.
	ContentPurposeTime
	// ContentPurposeDateTime accepts a date and time.
	ContentPurposeDateTime
	// ContentPurposeTerminal accepts shell/terminal input.
	ContentPurposeTerminal
	// ContentPurposeChat accepts chat or conversational text.
	ContentPurposeChat
)

// String returns a stable name for the content purpose.
func (p ContentPurpose) String() string {
	switch p {
	case ContentPurposeNormal:
		return "Normal"
	case ContentPurposeAlpha:
		return "Alpha"
	case ContentPurposeDigits:
		return "Digits"
	case ContentPurposeNumber:
		return "Number"
	case ContentPurposePhone:
		return "Phone"
	case ContentPurposeURL:
		return "URL"
	case ContentPurposeEmail:
		return "Email"
	case ContentPurposeName:
		return "Name"
	case ContentPurposePassword:
		return "Password"
	case ContentPurposePin:
		return "Pin"
	case ContentPurposeDate:
		return "Date"
	case ContentPurposeTime:
		return "Time"
	case ContentPurposeDateTime:
		return "DateTime"
	case ContentPurposeTerminal:
		return "Terminal"
	case ContentPurposeChat:
		return "Chat"
	default:
		return "Unknown"
	}
}

// ContentHint is a bitmask of input-method hints. Hints are advisory: a host
// may ignore a hint when its platform cannot express it. ContentHintNone is
// the zero value; the other ten values are independent flags.
type ContentHint uint16

const (
	// ContentHintNone requests no additional hint.
	ContentHintNone ContentHint = 0
	// ContentHintCompletion permits completion suggestions.
	ContentHintCompletion ContentHint = 1 << (iota - 1)
	// ContentHintSpellcheck permits spell checking.
	ContentHintSpellcheck
	// ContentHintAutoCapitalization permits automatic capitalization.
	ContentHintAutoCapitalization
	// ContentHintLowercase requests lowercase text.
	ContentHintLowercase
	// ContentHintUppercase requests uppercase text.
	ContentHintUppercase
	// ContentHintTitlecase requests title-case text.
	ContentHintTitlecase
	// ContentHintHiddenText marks text that must not be displayed.
	ContentHintHiddenText
	// ContentHintSensitiveData marks text that must not be learned or stored.
	ContentHintSensitiveData
	// ContentHintLatin requests Latin-script input where possible.
	ContentHintLatin
	// ContentHintMultiline permits line breaks.
	ContentHintMultiline
)

// ContentHintAutoCapitalize is the common spelling used by browser APIs. It
// is an alias of ContentHintAutoCapitalization.
const ContentHintAutoCapitalize = ContentHintAutoCapitalization

// Has reports whether all bits in hint are present.
func (h ContentHint) Has(hint ContentHint) bool {
	return hint != ContentHintNone && h&hint == hint
}

// String returns a stable, pipe-separated list of set hint names.
func (h ContentHint) String() string {
	if h == ContentHintNone {
		return stringNone
	}
	parts := make([]string, 0, 10)
	known := ContentHintNone
	for _, item := range []struct {
		hint ContentHint
		name string
	}{
		{ContentHintCompletion, "Completion"},
		{ContentHintSpellcheck, "Spellcheck"},
		{ContentHintAutoCapitalization, "AutoCapitalization"},
		{ContentHintLowercase, "Lowercase"},
		{ContentHintUppercase, "Uppercase"},
		{ContentHintTitlecase, "Titlecase"},
		{ContentHintHiddenText, "HiddenText"},
		{ContentHintSensitiveData, "SensitiveData"},
		{ContentHintLatin, "Latin"},
		{ContentHintMultiline, "Multiline"},
	} {
		if h.Has(item.hint) {
			parts = append(parts, item.name)
			known |= item.hint
		}
	}
	if unknown := h &^ known; unknown != 0 {
		parts = append(parts, "Unknown")
	}
	return strings.Join(parts, "|")
}

// IMEDeleteSurroundingEvent asks the consumer to delete text around its
// current insertion point. Before and After are non-negative UTF-8 byte
// counts in the surrounding text, not rune counts. The consumer should apply
// the deletion atomically and then publish the updated surrounding text.
type IMEDeleteSurroundingEvent struct {
	// Before is the number of bytes to delete immediately before the cursor.
	Before int

	// After is the number of bytes to delete immediately after the cursor.
	After int
}

// IsValid reports whether both deletion lengths are non-negative.
func (e IMEDeleteSurroundingEvent) IsValid() bool {
	return e.Before >= 0 && e.After >= 0
}

// IMEEventSourceV2 is the optional event extension for the richer IME
// contract. It intentionally does not modify EventSource, so every existing
// EventSource implementation remains source-compatible. Hosts should continue
// to expose EventSource's start/update/end callbacks and implement this
// interface when the platform can report these additional events.
type IMEEventSourceV2 interface {
	// OnIMECompositionUpdateV2 registers a callback for a full-fidelity
	// composition update. It is delivered in the same lifecycle as the
	// original EventSource callback, but carries a cursor range rather than a
	// single cursor position.
	OnIMECompositionUpdateV2(func(IMEComposition))

	// OnIMECanceled registers a callback delivered when an active composition
	// is canceled. Cancellation never commits the preedit text.
	OnIMECanceled(func())

	// OnIMEDisabled registers a callback delivered when the platform disables
	// text input (for example, after focus is lost or a password field is
	// selected). The callback carries no text and must not commit anything.
	OnIMEDisabled(func())

	// OnIMEDeleteSurrounding registers a callback for an IME request to delete
	// text around the current insertion point.
	OnIMEDeleteSurrounding(func(IMEDeleteSurroundingEvent))
}

// IMEEventSource is the unversioned spelling of IMEEventSourceV2. It is an
// optional capability and is not embedded in EventSource for compatibility.
type IMEEventSource = IMEEventSourceV2

// IMEControllerV2 extends the original IMEController without changing it.
// Consumers should use a type assertion before calling these optional
// methods. Implementations must treat IME as disabled until SetIMEEnabled(true)
// is called, and must stop exposing surrounding text after it is disabled.
type IMEControllerV2 interface {
	IMEController

	// SetIMECursorArea updates the caret rectangle used for candidate-window
	// placement. Coordinates are logical DIP in the window content area.
	SetIMECursorArea(area IMECursorArea)

	// SetIMEContentType supplies an advisory purpose and set of content hints.
	SetIMEContentType(purpose ContentPurpose, hints ContentHint)

	// SetIMESurroundingText supplies UTF-8 text and cursor/anchor byte offsets
	// for context-sensitive candidate generation and delete-surrounding
	// requests. Implementations should ignore invalid ranges.
	SetIMESurroundingText(text IMESurroundingText)

	// CancelIME cancels an active composition without committing its preedit.
	CancelIME()
}

// IMECapability identifies one optional operation in IMEContractVersion.
type IMECapability uint16

const (
	// IMECapabilityComposition indicates composition start/update/end events.
	IMECapabilityComposition IMECapability = 1 << iota
	// IMECapabilityCommit indicates committed text is reported exactly once.
	IMECapabilityCommit
	// IMECapabilityCancel indicates cancellation is supported.
	IMECapabilityCancel
	// IMECapabilityDisabled indicates OnIMEDisabled is supported.
	IMECapabilityDisabled
	// IMECapabilityDeleteSurrounding indicates delete-surrounding requests.
	IMECapabilityDeleteSurrounding
	// IMECapabilityCursorArea indicates candidate cursor-area updates.
	IMECapabilityCursorArea
	// IMECapabilitySurroundingText indicates surrounding-text updates.
	IMECapabilitySurroundingText
	// IMECapabilityContentPurpose indicates content-purpose support.
	IMECapabilityContentPurpose
	// IMECapabilityContentHints indicates content-hint support.
	IMECapabilityContentHints
)

// IMECapabilities describes the optional IME operations provided by a host.
// Version is the highest contract version understood by the provider. A zero
// Version means that no versioned contract is available.
type IMECapabilities struct {
	// Version is the supported IME contract version.
	Version uint

	// Features contains the operations supported at Version.
	Features IMECapability
}

// Supports reports whether all requested capabilities are available.
func (c IMECapabilities) Supports(capability IMECapability) bool {
	return capability != 0 && c.Features&capability == capability
}

// IMECapabilityProviderV2 lets a host advertise the optional IME contract.
// The method is deliberately separate from EventSource and IMEController so a
// host can expose capability discovery on whichever provider owns the window.
type IMECapabilityProviderV2 interface {
	IMECapabilities() IMECapabilities
}

// IMEContractV2 is the complete optional contract. A host may implement the
// smaller interfaces independently; consumers should prefer capability
// assertions so partial platform support remains useful.
type IMEContractV2 interface {
	IMEControllerV2
	IMEEventSourceV2
	IMECapabilityProviderV2
}

// DiscoverIMECapabilities returns the capabilities advertised by value. It
// returns false for legacy providers and for providers that return Version 0.
// A future Version is returned unchanged so consumers can safely use the bits
// they understand while ignoring newer operations.
func DiscoverIMECapabilities(value interface{}) (IMECapabilities, bool) {
	provider, ok := value.(IMECapabilityProviderV2)
	if !ok {
		return IMECapabilities{}, false
	}
	caps := provider.IMECapabilities()
	return caps, caps.Version != 0
}
