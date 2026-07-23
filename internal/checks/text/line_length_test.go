package text

import (
	"fmt"
	"strings"
	"testing"

	"github.com/wbd2023/Quill/internal/markers"
	"github.com/wbd2023/Quill/internal/style"
	"github.com/wbd2023/Quill/internal/testutil"
)

/* ---------------------------------------- Source files ---------------------------------------- */

func TestCheckLineLengthsFindsLongGoLines(t *testing.T) {
	repoRoot := t.TempDir()
	longLine := strings.Repeat("a", 101)
	path := testutil.WriteFile(
		t,
		repoRoot,
		"internal/example/example.go",
		"package example\n\nconst value = \""+longLine+"\"\n",
	)

	result, err := CheckLineLengths(
		repoRoot,
		[]string{path},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Diagnostics) == 0 {
		t.Fatal("expected long-line failure")
	}

	if !hasDiagnostic(result, "text/line-length/too-long", "internal/example/example.go", 3, "") {
		t.Fatalf("expected diagnostic to include offending file, got: %#v", result.Diagnostics)
	}
}

func TestCheckLineLengthsHonoursShellAllowMarker(t *testing.T) {
	repoRoot := t.TempDir()
	longLine := strings.Repeat("b", 101)
	source := strings.Join([]string{
		"#!/bin/bash",
		"set -euo pipefail",
		"echo \"" + longLine + "\" # " + markers.Text(longLineMarker),
		"",
	}, "\n")
	path := testutil.WriteFile(
		t,
		repoRoot,
		"tools/test.sh",
		source,
	)

	result, err := CheckLineLengths(
		repoRoot,
		[]string{path},
	)
	if err != nil {
		t.Fatalf("expected allow-marker line to pass, diagnostics: %#v", result.Diagnostics)
	}
}

/* ------------------------------------ Reference definitions ----------------------------------- */

// The Markdown tests below pin the approved narrow exception: a source line is still limited to
// 100 columns, and an overlong Markdown line is tolerated only when it is a Goldmark link
// reference definition whose lowercase http(s) destination alone exceeds 100 columns while the
// non-destination remainder stays within 100 columns. Every other overlong construct must fail.

func TestCheckLineLengthsAcceptsOverlongHTTPSReferenceDefinition(t *testing.T) {
	repoRoot := t.TempDir()
	destination := linkDestination("https://", 105) // 105 columns, strictly over the limit
	source := "[ref]: " + destination + "\n"        // remainder 7 columns, total 112

	path := testutil.WriteFile(t, repoRoot, "docs/links.md", source)
	result := checkLengths(t, repoRoot, path)
	requireNoDiagnostics(t, result)
}

func TestCheckLineLengthsRejectsReferenceDefinitionWithShortDestination(t *testing.T) {
	repoRoot := t.TempDir()
	destination := linkDestination("https://", 88) // 88 columns, not over the limit
	label := strings.Repeat("a", 14)               // "[aaaa...]: " remainder is 18 columns
	line := "[" + label + "]: " + destination      // 18 + 88 = 106; destination not eligible

	path := testutil.WriteFile(t, repoRoot, "docs/links.md", line+"\n")
	result := checkLengths(t, repoRoot, path)
	requireTooLong(t, result, "docs/links.md", 1, len(line))
}

func TestCheckLineLengthsAcceptsTwoLineReferenceDefinition(t *testing.T) {
	repoRoot := t.TempDir()
	destination := linkDestination("https://", 105) // eligible destination on its own line
	// Goldmark permits the optional title on the following line; the overlong line still carries
	// the destination, and the short title line stays within the limit.
	source := "[ref]: " + destination + "\n  \"Reference Title\"\n"

	path := testutil.WriteFile(t, repoRoot, "docs/links.md", source)
	result := checkLengths(t, repoRoot, path)
	requireNoDiagnostics(t, result)
}

func TestCheckLineLengthsRejectsLongBareURL(t *testing.T) {
	repoRoot := t.TempDir()
	line := linkDestination("https://", 105) // bare URL, not a reference definition

	path := testutil.WriteFile(t, repoRoot, "docs/links.md", line+"\n")
	result := checkLengths(t, repoRoot, path)
	requireTooLong(t, result, "docs/links.md", 1, len(line))
}

func TestCheckLineLengthsRejectsLongInlineLink(t *testing.T) {
	repoRoot := t.TempDir()
	destination := linkDestination("https://", 91)
	line := "[click here](" + destination + ")" // inline link, not a reference definition

	path := testutil.WriteFile(t, repoRoot, "docs/links.md", line+"\n")
	result := checkLengths(t, repoRoot, path)
	requireTooLong(t, result, "docs/links.md", 1, len(line))
}

func TestCheckLineLengthsRejectsNonHTTPSReferenceDestination(t *testing.T) {
	repoRoot := t.TempDir()
	destination := linkDestination("ftp://", 105) // non-http(s) scheme, 105 columns
	line := "[ref]: " + destination               // 7 + 105 = 112 columns

	path := testutil.WriteFile(t, repoRoot, "docs/links.md", line+"\n")
	result := checkLengths(t, repoRoot, path)
	requireTooLong(t, result, "docs/links.md", 1, len(line))
}

func TestCheckLineLengthsRejectsExactHundredReferenceDestination(t *testing.T) {
	repoRoot := t.TempDir()
	destination := linkDestination("https://", 100) // exactly 100, not strictly over the limit
	line := "[ref]: " + destination                 // 7 + 100 = 107 columns

	path := testutil.WriteFile(t, repoRoot, "docs/links.md", line+"\n")
	result := checkLengths(t, repoRoot, path)
	requireTooLong(t, result, "docs/links.md", 1, len(line))
}

func TestCheckLineLengthsRejectsEligibleDestinationWithExcessLabel(t *testing.T) {
	repoRoot := t.TempDir()
	destination := linkDestination("https://", 105) // eligible destination
	label := strings.Repeat("a", 100)               // "[aaaa...]: " remainder is 104 columns
	line := "[" + label + "]: " + destination       // 104 + 105 = 209 columns

	path := testutil.WriteFile(t, repoRoot, "docs/links.md", line+"\n")
	result := checkLengths(t, repoRoot, path)
	requireTooLong(t, result, "docs/links.md", 1, len(line))
}

func TestCheckLineLengthsRejectsMalformedDefinitionLikeLine(t *testing.T) {
	repoRoot := t.TempDir()
	destination := linkDestination("https://", 105)
	line := "[ref] " + destination // missing colon, not a valid link reference definition

	path := testutil.WriteFile(t, repoRoot, "docs/links.md", line+"\n")
	result := checkLengths(t, repoRoot, path)
	requireTooLong(t, result, "docs/links.md", 1, len(line))
}

// TestCheckLineLengthsAcceptsContinuationLineDestination pins the exact-span contract:
// CommonMark permits a single line ending between the label colon and the destination, so
// Goldmark stores a valid reference-definition destination in a later source segment than
// the label line. Only the real destination span earns the overlong allowance: the label
// line stays within the limit and the overlong continuation line must pass. The first-
// segment-only search rejects this qualifying destination, so this test fails until the
// destination is classified from its true segment.
func TestCheckLineLengthsAcceptsContinuationLineDestination(t *testing.T) {
	repoRoot := t.TempDir()
	destination := linkDestination("https://", 105) // eligible destination, strictly over the limit
	// The post-colon line ending places the destination on the continuation line; Goldmark
	// records it in a segment other than the first, while the short label line stays legal.
	source := "[ref]:\n" + destination + "\n"

	path := testutil.WriteFile(t, repoRoot, "docs/links.md", source)
	result := checkLengths(t, repoRoot, path)
	requireNoDiagnostics(t, result)
}

// TestCheckLineLengthsAcceptsAngleBracketedContinuationDestination pins the
// angle-delimiter contract alongside the exact-span contract: CommonMark permits a
// destination wrapped in optional angle brackets, which Goldmark strips from its parsed
// value, and the post-colon line ending still places the real span on the continuation
// line. Only the destination's own bytes (delimiters excluded) earn the overlong
// allowance, so the short label line stays legal and the overlong angle-delimited line
// must pass. A span search that compares at the raw opening '<' rejects this qualifying
// destination, so this test fails until the delimiters are skipped before matching.
func TestCheckLineLengthsAcceptsAngleBracketedContinuationDestination(t *testing.T) {
	repoRoot := t.TempDir()
	destination := linkDestination("https://", 105) // eligible destination, strictly over the limit
	// The destination carries the optional angle delimiters Goldmark strips from its
	// parsed value, and the post-colon line ending places the real span on the
	// continuation line; only the destination bytes (delimiters excluded) earn the
	// allowance, so the overlong line passes while the short label line stays legal.
	source := "[ref]:\n<" + destination + ">\n"

	path := testutil.WriteFile(t, repoRoot, "docs/links.md", source)
	result := checkLengths(t, repoRoot, path)
	requireNoDiagnostics(t, result)
}

// TestCheckLineLengthsRejectsLongURLInLabel proves the allowance pins the real destination
// span and does not leak to identical bytes that merely sit in a link label: the same long
// URL carried by the label does not qualify the line when the actual destination stays
// short, so the overlong line still fails.
func TestCheckLineLengthsRejectsLongURLInLabel(t *testing.T) {
	repoRoot := t.TempDir()
	longURL := linkDestination("https://", 105) // identical bytes, placed in the label
	shortDestination := "https://x.example"     // real destination, well within the limit
	line := "[" + longURL + "]: " + shortDestination

	path := testutil.WriteFile(t, repoRoot, "docs/links.md", line+"\n")
	result := checkLengths(t, repoRoot, path)
	requireTooLong(t, result, "docs/links.md", 1, len(line))
}

/* ------------------------------------------- Helpers ------------------------------------------ */

// checkLengths runs CheckLineLengths and fails the test on any returned error.
func checkLengths(t *testing.T, repoRoot string, path string) (result style.ExecutionResult) {
	t.Helper()
	result, err := CheckLineLengths(repoRoot, []string{path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return result
}

// requireNoDiagnostics fails the test when the result carries any diagnostics.
func requireNoDiagnostics(t *testing.T, result style.ExecutionResult) {
	t.Helper()
	if len(result.Diagnostics) != 0 {
		t.Fatalf("expected no diagnostics, got: %#v", result.Diagnostics)
	}
}

// requireTooLong fails the test unless there is exactly one too-long diagnostic for the given
// file and line, carrying the expected column count in its message.
func requireTooLong(
	t *testing.T,
	result style.ExecutionResult,
	file string,
	line int,
	columns int,
) {
	t.Helper()
	if len(result.Diagnostics) != 1 {
		t.Fatalf("expected exactly one diagnostic for %s:%d, got: %#v",
			file, line, result.Diagnostics)
	}
	if !hasDiagnostic(result, "text/line-length/too-long", file, line,
		fmt.Sprintf("%d columns", columns)) {
		t.Fatalf("expected too-long diagnostic (%d columns) for %s:%d, got: %#v",
			columns, file, line, result.Diagnostics)
	}
}

// linkDestination builds an exact-length lowercase link destination with the given scheme.
// scheme includes the "://" suffix, e.g. "https://" or "ftp://".
func linkDestination(scheme string, columns int) (destination string) {
	prefix := scheme + "example.com/"
	destination = prefix + strings.Repeat("a", columns-len(prefix))
	return destination
}
