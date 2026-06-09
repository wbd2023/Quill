package bash

import (
	"testing"

	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/fixtures/profiles"
	"ciphera/tools/internal/style"
)

func TestCheckSafetyFindsConventionAndSafetyViolations(t *testing.T) {
	repoRoot := t.TempDir()
	fixtures.WriteFile(
		t,
		repoRoot,
		"tools/sample.sh",
		"#!/bin/bash\n"+
			"set -euo pipefail\n\n"+
			"readonly colour_reset='\\033[0m'\n"+
			"tmp_dir=/tmp/example\n"+
			"BadFunction() {\n"+
			"\twhich git >/dev/null 2>&1\n"+
			"\tprintf '%s\\n' \"$colour_reset\"\n"+
			"}\n\n"+
			"worker() {\n"+
			"\tprintf 'work\\n'\n"+
			"}\n\n"+
			"printf 'value' | while read -r line; do\n"+
			"\tprintf '%s\\n' \"$line\"\n"+
			"done\n\n"+
			"# shellcheck disable=SC2086\n"+
			"worker\n",
	)

	result, err := CheckSafety(repoRoot, profiles.RepositoryConfig(t), style.Scope("tools"))
	if err == nil {
		t.Fatal("expected bash safety failure")
	}

	required := []struct {
		code    string
		message string
	}{
		{"bash/safety/naming", "Bash function names should use lower-case with underscores"},
		{"bash/safety/naming", "Bash constants and exported variables should use upper-case"},
		{"bash/safety/script-shape", "detect dependencies with command -v, not which"},
		{"bash/safety/temp-path", "temporary resources must be created with mktemp"},
		{"bash/safety/script-shape", "avoid cmd | while read loops when loop state must survive"},
		{"bash/safety/suppression", "shellcheck suppressions must include rule IDs"},
		{"bash/safety/script-shape", "non-trivial Bash scripts must keep main()"},
		{"bash/safety/script-shape", "must end with main \"$@\""},
	}
	for _, expected := range required {
		if hasDiagnostic(result, expected.code, "", 0, expected.message) {
			continue
		}

		t.Fatalf(
			"expected diagnostic containing %q, got: %#v",
			expected.message,
			result.Diagnostics,
		)
	}
}

func TestCheckSafetyPassesCleanScript(t *testing.T) {
	repoRoot := t.TempDir()
	fixtures.WriteFile(
		t,
		repoRoot,
		"tools/sample.sh",
		"#!/bin/bash\n"+
			"set -euo pipefail\n\n"+
			"readonly COLOUR_RESET='\\033[0m'\n\n"+
			"main() {\n"+
			"\tlocal tmp_dir\n"+
			"\ttmp_dir=\"$(mktemp -d)\"\n"+
			"\ttrap 'rm -rf \"$tmp_dir\"' EXIT\n"+
			"\tif ! command -v git >/dev/null 2>&1; then\n"+
			"\t\tprintf 'git is required\\n'\n"+
			"\t\treturn 1\n"+
			"\tfi\n"+
			"}\n\n"+
			"main \"$@\"\n",
	)

	result, err := CheckSafety(repoRoot, profiles.RepositoryConfig(t), style.Scope("tools"))
	if err != nil {
		t.Fatalf("expected bash safety check to pass, diagnostics: %#v", result.Diagnostics)
	}
}
