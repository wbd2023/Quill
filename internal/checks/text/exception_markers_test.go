package text

import (
	"testing"

	"github.com/wbd2023/Quill/internal/style"
	"github.com/wbd2023/Quill/internal/testutil"
	"github.com/wbd2023/Quill/internal/testutil/profiles"
)

func TestCheckExceptionMarkersRejectsMalformedMarkers(t *testing.T) {
	repoRoot := t.TempDir()
	testutil.WriteFile(
		t,
		repoRoot,
		"tools/test.sh",
		"#!/bin/bash\nset -euo pipefail\necho test # style: allow long-line\n",
	)

	result, err := CheckExceptionMarkers(
		repoRoot,
		profiles.RepositoryConfig(t),
		style.Scope("all"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Diagnostics) == 0 {
		t.Fatal("expected malformed exception marker to fail")
	}

	if !hasDiagnostic(result, "text/exception-markers/invalid", "tools/test.sh", 3, "") {
		t.Fatalf("expected malformed-marker diagnostic, got: %#v", result.Diagnostics)
	}
}

func TestCheckExceptionMarkersAcceptsValidMarkers(t *testing.T) {
	repoRoot := t.TempDir()
	testutil.WriteFile(
		t,
		repoRoot,
		"internal/example/example.go",
		"package example\n\n// style: allow-non-ascii because: test vector\nconst value = \"ok\"\n",
	)

	result, err := CheckExceptionMarkers(
		repoRoot,
		profiles.RepositoryConfig(t),
		style.Scope("all"),
	)
	if err != nil {
		t.Fatalf("expected valid marker to pass, diagnostics: %#v", result.Diagnostics)
	}
}
