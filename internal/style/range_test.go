package style

import (
	"errors"
	"strings"
	"testing"
)

/* ----------------------------------------- Acceptance ----------------------------------------- */

func TestVerifyRangeAcceptsValidLocations(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		file string
		rng  Range
	}{
		{"empty file unknown range", "", Range{}},
		{"relative file unknown range", "a/b.go", Range{}},
		{"relative file start line only", "a/b.go", Range{Start: Position{Line: 3}}},
		{"relative file start line and byte column", "a/b.go",
			Range{Start: Position{Line: 3, Column: 7}}},

		// One-line ordering.
		{"one-line span end after start", "a/b.go",
			Range{Start: Position{Line: 3, Column: 4}, End: Position{Line: 3, Column: 6}}},
		{"one-line point range start equals end", "a/b.go",
			Range{Start: Position{Line: 3, Column: 4}, End: Position{Line: 3, Column: 4}}},
		{"one-line span unknown columns", "a/b.go",
			Range{Start: Position{Line: 3}, End: Position{Line: 3}}},

		// Multi-line ordering.
		{"multi-line span", "a/b.go",
			Range{Start: Position{Line: 3}, End: Position{Line: 5}}},
		{"multi-line span end column on later line", "a/b.go",
			Range{Start: Position{Line: 3, Column: 4}, End: Position{Line: 5, Column: 2}}},
	}

	for _, test := range cases {
		if err := VerifyRange(test.file, test.rng); err != nil {
			t.Errorf("%s: expected acceptance, got %v", test.name, err)
		}
	}
}

// TestVerifyRangeTreatsColumnsAsUTF8Bytes exercises the byte-column contract: columns count UTF-8
// bytes, not runes. In a line beginning with a two-byte rune, the next rune starts at byte column 3
// (not rune column 2), so a span computed from byte offsets is accepted and ordering compares those
// byte offsets.
func TestVerifyRangeTreatsColumnsAsUTF8Bytes(t *testing.T) {
	t.Parallel()

	// "\u00e9xample" - the escaped rune is two bytes, so 'x' starts at byte column 3.
	line := "\u00e9xample"
	startColumn := strings.IndexByte(line, 'x') + 1 // 1-based byte offset
	endColumn := startColumn + 1                    // exclusive end covering the single byte

	if err := VerifyRange("a/b.go", Range{
		Start: Position{Line: 1, Column: startColumn},
		End:   Position{Line: 1, Column: endColumn},
	}); err != nil {
		t.Fatalf("expected byte-column span accepted: %v", err)
	}

	// An end column before the start column is rejected even on a single line.
	if err := VerifyRange("a/b.go", Range{
		Start: Position{Line: 1, Column: endColumn},
		End:   Position{Line: 1, Column: startColumn},
	}); !errors.Is(err, errRangeEndBeforeStart) {
		t.Fatalf("expected end-before-start rejection, got %v", err)
	}
}

/* ----------------------------------------- Rejections ----------------------------------------- */

func TestVerifyRangeRejectsBadPaths(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		file    string
		wantErr error
	}{
		{"absolute path", "/a/b.go", errPathAbsolute},
		{"windows drive path", "C:/x.go", errPathAbsolute},
		{"parent escape prefix", "../a.go", errPathEscape},
		{"mid parent escape", "a/../b.go", errPathEscape},
		{"dot segment", "./a.go", errPathEscape},
		{"trailing slash empty segment", "a/", errPathEscape},
		{"backslash separator", "a\\b.go", errPathUnclean},
		{"nul byte", "a\x00b", errPathUnclean},
	}

	for _, test := range cases {
		err := VerifyRange(test.file, Range{})
		if !errors.Is(err, test.wantErr) {
			t.Errorf("%s: expected %v, got %v", test.name, test.wantErr, err)
		}
	}
}

func TestVerifyRangeRejectsBadPositions(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		rng     Range
		wantErr error
	}{
		{"column without line", Range{Start: Position{Column: 3}}, errColumnWithoutLine},
		{"end column without line",
			Range{Start: Position{Line: 1}, End: Position{Column: 3}}, errColumnWithoutLine},
		{"negative line", Range{Start: Position{Line: -1}}, errPositionNegative},
		{"negative column", Range{Start: Position{Line: 1, Column: -2}}, errPositionNegative},
		{"one-line end before start",
			Range{Start: Position{Line: 3, Column: 6}, End: Position{Line: 3, Column: 4}},
			errRangeEndBeforeStart},
		{"multi-line end before start",
			Range{Start: Position{Line: 5}, End: Position{Line: 3}}, errRangeEndBeforeStart},
		{"end without start", Range{End: Position{Line: 3}}, errRangeEndWithoutStart},
	}

	for _, test := range cases {
		err := VerifyRange("a/b.go", test.rng)
		if !errors.Is(err, test.wantErr) {
			t.Errorf("%s: expected %v, got %v", test.name, test.wantErr, err)
		}
	}
}

func TestVerifyRangeRejectsRangeWithoutFile(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		rng  Range
	}{
		{"known start without file", Range{Start: Position{Line: 1}}},
		{"known start and column without file", Range{Start: Position{Line: 1, Column: 1}}},
		{"known end without file", Range{End: Position{Line: 1}}},
	}

	for _, test := range cases {
		err := VerifyRange("", test.rng)
		if !errors.Is(err, errRangeWithoutFile) {
			t.Errorf("%s: expected %v, got %v", test.name, errRangeWithoutFile, err)
		}
	}
}

/* -------------------------------------- Knownness Helpers ------------------------------------- */

func TestPositionAndRangeKnownness(t *testing.T) {
	t.Parallel()

	if (Position{}).IsKnown() {
		t.Error("zero position should be unknown")
	}
	if !(Position{Line: 1}.IsKnown()) {
		t.Error("position with a line should be known")
	}
	if !(Range{}).IsUnknown() {
		t.Error("zero range should be unknown")
	}
	if (Range{Start: Position{Line: 1}}).IsUnknown() {
		t.Error("range with a known start should not be unknown")
	}
	if (Range{End: Position{Line: 1}}).IsUnknown() {
		t.Error("range with a known end should not be unknown")
	}
}
