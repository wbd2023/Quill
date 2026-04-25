package bashstyle

import (
	"strings"
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

	output, err := CheckStructure(
		repoRoot,
		profiles.RepositoryConfig(t),
		contract.ScopeTools,
	)
	if err == nil {
		t.Fatal("expected bash structure failure")
	}

	if !strings.Contains(output, "missing set -euo pipefail") {
		t.Fatalf("expected strict-mode message, got:\n%s", output)
	}
}
