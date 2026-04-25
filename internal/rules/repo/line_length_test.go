package repostyle

import (
	"strings"
	"testing"

	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/styleguide"
)

func TestCheckLineLengthsFindsLongGoLines(t *testing.T) {
	repoRoot := t.TempDir()
	longLine := strings.Repeat("a", 101)
	path := fixtures.WriteFile(
		t,
		repoRoot,
		"internal/example/example.go",
		"package example\n\nconst value = \""+longLine+"\"\n",
	)

	output, err := CheckLineLengths(
		repoRoot,
		[]string{path},
	)
	if err == nil {
		t.Fatal("expected long-line failure")
	}

	if !strings.Contains(output, "internal/example/example.go:3") {
		t.Fatalf("expected output to include offending file, got:\n%s", output)
	}
}

func TestCheckLineLengthsHonoursShellAllowMarker(t *testing.T) {
	repoRoot := t.TempDir()
	longLine := strings.Repeat("b", 101)
	source := strings.Join([]string{
		"#!/bin/bash",
		"set -euo pipefail",
		"echo \"" + longLine + "\" # " + styleguide.ExceptionMarker(styleguide.ExceptionLongLine),
		"",
	}, "\n")
	path := fixtures.WriteFile(
		t,
		repoRoot,
		"tools/test.sh",
		source,
	)

	output, err := CheckLineLengths(
		repoRoot,
		[]string{path},
	)
	if err != nil {
		t.Fatalf("expected allow-marker line to pass, output:\n%s", output)
	}
}
