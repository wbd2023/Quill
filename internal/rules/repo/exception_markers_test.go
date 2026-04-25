package repostyle

import (
	"strings"
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/fixtures/profiles"
)

func TestCheckExceptionMarkersRejectsMalformedMarkers(t *testing.T) {
	repoRoot := t.TempDir()
	fixtures.WriteFile(
		t,
		repoRoot,
		"tools/test.sh",
		"#!/bin/bash\nset -euo pipefail\necho test # style: allow long-line\n",
	)

	output, err := CheckExceptionMarkers(
		repoRoot,
		profiles.RepositoryConfig(t),
		contract.ScopeTools,
	)
	if err == nil {
		t.Fatal("expected malformed exception marker to fail")
	}

	if !strings.Contains(output, "tools/test.sh:3") {
		t.Fatalf("expected output to include malformed marker location, got:\n%s", output)
	}
}

func TestCheckExceptionMarkersAcceptsValidMarkers(t *testing.T) {
	repoRoot := t.TempDir()
	fixtures.WriteFile(
		t,
		repoRoot,
		"internal/example/example.go",
		"package example\n\n// style: allow-non-ascii because: test vector\nconst value = \"ok\"\n",
	)

	output, err := CheckExceptionMarkers(
		repoRoot,
		profiles.RepositoryConfig(t),
		contract.ScopeApp,
	)
	if err != nil {
		t.Fatalf("expected valid marker to pass, output:\n%s", output)
	}
}
