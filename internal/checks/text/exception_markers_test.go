package text

import (
	"testing"

	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/fixtures/profiles"
	"ciphera/tools/internal/style"
)

func TestCheckExceptionMarkersRejectsMalformedMarkers(t *testing.T) {
	repoRoot := t.TempDir()
	fixtures.WriteFile(
		t,
		repoRoot,
		"tools/test.sh",
		"#!/bin/bash\nset -euo pipefail\necho test # style: allow long-line\n",
	)

	result, err := CheckExceptionMarkers(
		repoRoot,
		profiles.RepositoryConfig(t),
		style.Scope("tools"),
	)
	if err == nil {
		t.Fatal("expected malformed exception marker to fail")
	}

	if !hasDiagnostic(result, "text/exception-markers/invalid", "tools/test.sh", 3, "") {
		t.Fatalf("expected malformed-marker diagnostic, got: %#v", result.Diagnostics)
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

	result, err := CheckExceptionMarkers(
		repoRoot,
		profiles.RepositoryConfig(t),
		style.Scope("app"),
	)
	if err != nil {
		t.Fatalf("expected valid marker to pass, diagnostics: %#v", result.Diagnostics)
	}
}
