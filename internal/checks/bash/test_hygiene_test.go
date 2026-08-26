package bash

import (
	"testing"

	"github.com/wbd2023/quill/internal/style"
	"github.com/wbd2023/quill/internal/testutil"
	"github.com/wbd2023/quill/internal/testutil/profiles"
)

func TestCheckTestHygieneRequiresTrapCleanup(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(
		t,
		root,
		"test/smoke/example_test.sh",
		"#!/bin/bash\nset -euo pipefail\n\ntmp_dir=\"$(mktemp -d)\"\n",
	)

	result, err := CheckTestHygiene(
		root,
		profiles.RepositoryConfig(),
		style.Scope("all"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Diagnostics) == 0 {
		t.Fatal("expected bash test hygiene failure")
	}

	if !hasDiagnostic(
		result,
		"bash/test-hygiene/missing-cleanup",
		"",
		0,
		"mktemp must install trap-based cleanup",
	) {
		t.Fatalf("expected trap-cleanup diagnostic, got: %#v", result.Diagnostics)
	}
}

func TestCheckTestHygieneAcceptsTrapCleanup(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(
		t,
		root,
		"test/smoke/example_test.sh",
		"#!/bin/bash\nset -euo pipefail\n\ntmp_dir=\"$(mktemp -d)\"\n"+
			"trap 'rm -rf \"$tmp_dir\"' EXIT\n",
	)

	result, err := CheckTestHygiene(
		root,
		profiles.RepositoryConfig(),
		style.Scope("all"),
	)
	if err != nil {
		t.Fatalf("expected bash test hygiene fixture to pass, diagnostics: %#v", result.Diagnostics)
	}
}
