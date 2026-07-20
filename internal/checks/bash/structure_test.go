package bash

import (
	"testing"

	"github.com/wbd2023/Quill/internal/style"
	"github.com/wbd2023/Quill/internal/testutil"
	"github.com/wbd2023/Quill/internal/testutil/profiles"
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
		style.Scope("all"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Diagnostics) == 0 {
		t.Fatal("expected bash structure failure")
	}

	if !hasDiagnostic(result, "bash/structure/invalid", "", 0, "missing set -euo pipefail") {
		t.Fatalf("expected strict-mode diagnostic, got: %#v", result.Diagnostics)
	}
}
