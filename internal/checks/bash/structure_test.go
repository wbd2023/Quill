package bash

import (
	"testing"

	"ciphera/tools/internal/style"
	"ciphera/tools/internal/testutil"
	"ciphera/tools/internal/testutil/profiles"
)

func TestCheckStructureFindsMissingStrictMode(t *testing.T) {
	repoRoot := t.TempDir()
	testutil.WriteFile(
		t,
		repoRoot,
		"tools/test.sh",
		"#!/bin/bash\nprintf 'hello\\n'\n",
	)

	result, err := CheckStructure(
		repoRoot,
		profiles.RepositoryConfig(t),
		style.Scope("tools"),
	)
	if err == nil {
		t.Fatal("expected bash structure failure")
	}

	if !hasDiagnostic(result, "bash/structure/invalid", "", 0, "missing set -euo pipefail") {
		t.Fatalf("expected strict-mode diagnostic, got: %#v", result.Diagnostics)
	}
}
