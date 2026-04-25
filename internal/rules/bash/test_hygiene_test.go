package bashstyle

import (
	"strings"
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/fixtures/profiles"
)

func TestCheckTestHygieneRequiresTrapCleanup(t *testing.T) {
	repoRoot := t.TempDir()
	fixtures.WriteFile(
		t,
		repoRoot,
		"test/smoke/example_test.sh",
		"#!/bin/bash\nset -euo pipefail\n\ntmp_dir=\"$(mktemp -d)\"\n",
	)

	output, err := CheckTestHygiene(
		repoRoot,
		profiles.RepositoryConfig(t),
		contract.ScopeAll,
	)
	if err == nil {
		t.Fatal("expected bash test hygiene failure")
	}

	if !strings.Contains(output, "mktemp must install trap-based cleanup") {
		t.Fatalf("expected trap-cleanup violation, got:\n%s", output)
	}
}

func TestCheckTestHygieneAcceptsTrapCleanup(t *testing.T) {
	repoRoot := t.TempDir()
	fixtures.WriteFile(
		t,
		repoRoot,
		"test/smoke/example_test.sh",
		"#!/bin/bash\nset -euo pipefail\n\ntmp_dir=\"$(mktemp -d)\"\n"+
			"trap 'rm -rf \"$tmp_dir\"' EXIT\n",
	)

	output, err := CheckTestHygiene(
		repoRoot,
		profiles.RepositoryConfig(t),
		contract.ScopeAll,
	)
	if err != nil {
		t.Fatalf("expected bash test hygiene fixture to pass, output:\n%s", output)
	}
}
