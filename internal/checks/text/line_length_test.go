package text

import (
	"strings"
	"testing"

	"ciphera/tools/internal/markers"
	"ciphera/tools/internal/testutil"
)

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
