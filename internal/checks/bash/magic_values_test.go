package bash

import (
	"strings"
	"testing"

	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/fixtures/profiles"
	"ciphera/tools/internal/style"
)

func TestCheckMagicValuesFindsNonTrivialLiterals(t *testing.T) {
	repoRoot := t.TempDir()
	source := strings.Join([]string{
		"#!/bin/bash",
		"set -euo pipefail",
		"exit 2",
		"if [ \"$value\" -eq 9 ]; then",
		"\texit 0",
		"fi",
		"",
	}, "\n")
	fixtures.WriteFile(
		t,
		repoRoot,
		"tools/example.sh",
		source,
	)

	result, err := CheckMagicValues(
		repoRoot,
		profiles.RepositoryConfig(t),
		style.Scope("tools"),
	)
	if err == nil {
		t.Fatal("expected bash magic-value failure")
	}

	if !hasDiagnostic(result, "bash/magic-values/non-trivial", "tools/example.sh", 3, "") ||
		!hasDiagnostic(result, "bash/magic-values/non-trivial", "tools/example.sh", 4, "") {
		t.Fatalf("expected diagnostics to include offending lines, got: %#v", result.Diagnostics)
	}
}

func TestCheckMagicValuesAllowsTrivialLiterals(t *testing.T) {
	repoRoot := t.TempDir()
	source := strings.Join([]string{
		"#!/bin/bash",
		"set -euo pipefail",
		"if [ \"$#\" -eq 1 ]; then",
		"\texit 1",
		"fi",
		"head -1 file.txt",
		"",
	}, "\n")
	fixtures.WriteFile(
		t,
		repoRoot,
		"tools/example.sh",
		source,
	)

	result, err := CheckMagicValues(
		repoRoot,
		profiles.RepositoryConfig(t),
		style.Scope("tools"),
	)
	if err != nil {
		t.Fatalf("expected trivial literals to pass, diagnostics: %#v", result.Diagnostics)
	}
}
