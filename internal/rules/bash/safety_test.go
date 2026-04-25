package bashstyle

import (
	"strings"
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/fixtures/profiles"
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

	output, err := CheckSafety(repoRoot, profiles.RepositoryConfig(t), contract.ScopeTools)
	if err == nil {
		t.Fatal("expected bash safety failure")
	}

	required := []string{
		"Bash function names should use lower-case with underscores",
		"Bash constants and exported variables should use upper-case with underscores",
		"detect dependencies with command -v, not which",
		"temporary resources must be created with mktemp",
		"avoid cmd | while read loops when loop state must survive",
		"shellcheck suppressions must include rule IDs and a short reason",
		"non-trivial Bash scripts must keep main() as the bottom-most function",
		"must end with main \"$@\"",
	}
	for _, snippet := range required {
		if strings.Contains(output, snippet) {
			continue
		}

		t.Fatalf("expected %q in output, got:\n%s", snippet, output)
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

	output, err := CheckSafety(repoRoot, profiles.RepositoryConfig(t), contract.ScopeTools)
	if err != nil {
		t.Fatalf("expected bash safety check to pass, output:\n%s", output)
	}
}
