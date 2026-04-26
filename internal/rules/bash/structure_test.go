package bash

import (
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/fixtures/profiles"
)

func TestCheckStructureFindsMissingStrictMode(t *testing.T) {
	repoRoot := t.TempDir()
	fixtures.WriteFile(
		t,
		repoRoot,
		"tools/test.sh",
		"#!/bin/bash\nprintf 'hello\\n'\n",
	)

	result, err := CheckStructure(
		repoRoot,
		profiles.RepositoryConfig(t),
		contract.Scope("tools"),
	)
	if err == nil {
		t.Fatal("expected bash structure failure")
	}

	if !hasDiagnostic(result, "bash/structure/invalid", "", 0, "missing set -euo pipefail") {
		t.Fatalf("expected strict-mode diagnostic, got: %#v", result.Diagnostics)
	}
}
