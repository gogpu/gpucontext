// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

package gpucontext

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

// IME contract (ADR-045). One IMEController, extended in place before v1.0,
// and one optional IMEEventSource found by type assertion. All text positions
// are UTF-8 byte offsets, the representation used by the Wayland text-input
// protocol, GTK4 and winit; native backends convert their platform's units
// (UTF-16 on Windows and macOS) at the boundary.

// IMETextRange is a half-open range [Start, End) in a UTF-8 string.
//
// Start and End are byte offsets, not rune or UTF-16 indexes, and must land on
// UTF-8 rune boundaries.
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
// IME payloads use UTF-8 offsets. An out-of-range or mid-rune offset is
// invalid; it is never clamped, so a host drops such a payload instead of
// guessing.
func (r IMETextRange) IsValid(text string) bool {
	if !utf8.ValidString(text) || r.Start < 0 || r.End < r.Start || r.End > len(text) {
		return false
	}
	if r.Start < len(text) && !utf8.RuneStart(text[r.Start]) {
		return false
	}
	return r.End == len(text) || utf8.RuneStart(text[r.End])
}

// imeHiddenCursor is the Cursor value of an IMEState whose IME asked for no
// visible cursor (Wayland text-input-v3 sends -1 for both ends).
var imeHiddenCursor = IMETextRange{Start: -1, End: -1}

// IMEState is the current Input Method Editor composition, delivered through
// EventSource.OnIMECompositionUpdate. It is used for CJK and other complex
// text input: the user types phonetic characters that the IME converts, and
// the consumer renders CompositionText inline with composition styling.
//
// Cursor and Selection are UTF-8 byte ranges in CompositionText. Most IMEs
// report a collapsed Cursor (Start == End). Selection is the target or
// clause segment the IME is converting, when the platform reports one. A
// Cursor of {-1, -1} means the IME asked for no visible cursor.
type IMEState struct {
	// Composing indicates whether IME is currently in composition mode.
	Composing bool

	// CompositionText is the current preedit text.
	CompositionText string

	// Cursor is the cursor range in CompositionText, or {-1, -1} when the IME
	// hides the cursor.
	Cursor IMETextRange

	// Selection is the target or clause range in CompositionText; empty when
	// the platform reports none.
	Selection IMETextRange
}

// HasCursor reports whether the IME asked for a visible cursor.
func (s IMEState) HasCursor() bool {
	return s.Cursor != imeHiddenCursor
}

// IsValid reports whether Cursor (unless hidden) and Selection are valid UTF-8
// ranges in CompositionText.
func (s IMEState) IsValid() bool {
	cursorValid := !s.HasCursor() || s.Cursor.IsValid(s.CompositionText)
	return cursorValid && s.Selection.IsValid(s.CompositionText)
}

// IMESurroundingText is the text around the insertion point supplied to the
// platform IME. Cursor and Anchor are UTF-8 byte offsets into Text. They differ
// when the surrounding text contains a selection; their order is preserved
// because some native IMEs use the direction of the selection.
//
// A platform may pass the IME only part of Text around the cursor: Wayland
// text-input-v3 caps a request at 4000 bytes. Consumers send a window around
// the insertion point rather than a whole document.
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
// WindowProvider and the EventSource pointer callbacks; the host scales them
// to the platform's units. Width and Height may be zero when only a caret
// point is available.
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
// follow the cross-platform purpose set of winit and the Wayland text-input
// protocol, plus Chat. A host treats a value its platform cannot express
// (Chat on most platforms) and any unknown value as Normal.
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
	// ContentPurposeChat accepts chat or conversational text. Platforms
	// without a chat purpose treat it as Normal.
	ContentPurposeChat

	// contentPurposeEnd stays last: the name-coverage test walks every value
	// below it.
	contentPurposeEnd
)

// contentPurposeNames is keyed by the constants, so a new purpose without a
// name shows up in TestContentPurposeStringCoversEveryValue.
var contentPurposeNames = map[ContentPurpose]string{
	ContentPurposeNormal:   "Normal",
	ContentPurposeAlpha:    "Alpha",
	ContentPurposeDigits:   "Digits",
	ContentPurposeNumber:   "Number",
	ContentPurposePhone:    "Phone",
	ContentPurposeURL:      "URL",
	ContentPurposeEmail:    "Email",
	ContentPurposeName:     "Name",
	ContentPurposePassword: "Password",
	ContentPurposePin:      "Pin",
	ContentPurposeDate:     "Date",
	ContentPurposeTime:     "Time",
	ContentPurposeDateTime: "DateTime",
	ContentPurposeTerminal: "Terminal",
	ContentPurposeChat:     "Chat",
}

// String returns a stable name for the content purpose.
func (p ContentPurpose) String() string {
	if name, ok := contentPurposeNames[p]; ok {
		return name
	}
	return stringUnknown
}

// ContentHint is a bitmask of input-method hints. Hints are advisory: a host
// ignores a hint its platform cannot express. ContentHintNone is the zero
// value; the other values are independent flags.
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

	// contentHintEnd stays last: the name-coverage test checks that every
	// flag below it has a name.
	contentHintEnd
)

// contentHintNames lists every hint flag in bit order, keyed by the constants.
var contentHintNames = []struct {
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
}

// Has reports whether all bits in hint are present.
func (h ContentHint) Has(hint ContentHint) bool {
	return hint != ContentHintNone && h&hint == hint
}

// String returns a stable, pipe-separated list of set hint names.
func (h ContentHint) String() string {
	if h == ContentHintNone {
		return stringNone
	}
	parts := make([]string, 0, len(contentHintNames))
	known := ContentHintNone
	for _, item := range contentHintNames {
		if h.Has(item.hint) {
			parts = append(parts, item.name)
			known |= item.hint
		}
	}
	if unknown := h &^ known; unknown != 0 {
		parts = append(parts, "0x"+strconv.FormatUint(uint64(unknown), 16))
	}
	return strings.Join(parts, "|")
}

// IMECapabilities tells a consumer which optional IMEController operations
// the host's platform really performs. Composition, commit, cancel and
// enable/disable are the baseline of every IMEController and have no bit.
// The groups follow winit, where each one is a real platform boundary.
type IMECapabilities uint8

const (
	// IMECapabilityCursorArea means SetIMECursorArea positions the candidate
	// window.
	IMECapabilityCursorArea IMECapabilities = 1 << iota
	// IMECapabilityContentType means SetIMEContentType reaches the platform
	// IME (purpose and hints).
	IMECapabilityContentType
	// IMECapabilitySurroundingText means SetIMESurroundingText reaches the
	// platform IME and delete-surrounding requests are delivered through
	// IMEEventSource.
	IMECapabilitySurroundingText
)

// Has reports whether all capabilities in c2 are present.
func (c IMECapabilities) Has(c2 IMECapabilities) bool {
	return c2 != 0 && c&c2 == c2
}

// IMEController lets widgets control the platform IME. It is implemented by
// the host window system. IME is disabled until SetIMEEnabled(true); a host
// stops exposing surrounding text once it is disabled. Methods for an
// operation the platform does not support (see IMECapabilities) are no-ops.
//
// Threading: every method must be safe to call from any goroutine, including
// from a draw callback. Platform IME APIs are bound to one thread (TSF and
// IMM32 to the window's thread, Cocoa to the main thread, Wayland to its event
// queue), so an implementation marshals the call there and may apply it
// asynchronously, no later than the next frame. It never blocks the caller
// waiting for that thread.
type IMEController interface {
	// SetIMEEnabled enables or disables IME for the focused input. When
	// disabled, key presses are delivered directly without IME processing,
	// as for password or non-text inputs.
	SetIMEEnabled(enabled bool)

	// SetIMECursorArea updates the caret rectangle used to place the
	// candidate window. Coordinates are logical DIP in the content area.
	SetIMECursorArea(area IMECursorArea)

	// SetIMEContentType supplies an advisory purpose and set of hints.
	SetIMEContentType(purpose ContentPurpose, hints ContentHint)

	// SetIMESurroundingText supplies the text around the insertion point for
	// context-sensitive candidates and delete-surrounding requests. A host
	// ignores an invalid value (see IMESurroundingText.IsValid).
	SetIMESurroundingText(text IMESurroundingText)

	// CancelIME cancels an active composition without committing it.
	CancelIME()

	// IMECapabilities reports which optional operations the platform
	// performs.
	IMECapabilities() IMECapabilities
}

// IMEDeleteSurroundingEvent asks the consumer to delete text around its
// current insertion point. Before and After are non-negative UTF-8 byte
// counts in the surrounding text, not rune counts. The consumer applies the
// deletion atomically and then publishes the updated surrounding text.
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

// IMEEventSource is the optional IME event extension of an EventSource,
// found by type assertion. EventSource's OnIMECompositionStart, Update and
// End callbacks carry composition; these carry the rest of the IME lifecycle.
type IMEEventSource interface {
	// OnIMECanceled registers a callback delivered when an active composition
	// is canceled. Cancellation never commits the preedit text.
	OnIMECanceled(fn func())

	// OnIMEDisabled registers a callback delivered when the platform disables
	// text input, for example after focus is lost. It carries no text and
	// must not commit anything.
	OnIMEDisabled(fn func())

	// OnIMEDeleteSurrounding registers a callback for an IME request to
	// delete text around the current insertion point.
	OnIMEDeleteSurrounding(fn func(IMEDeleteSurroundingEvent))
}
