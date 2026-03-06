package style_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var requiredStyleToolNames = []string{
	"misspell",
	"golangci-lint",
	"shfmt",
	"rg",
	"shellcheck",
}

/* -------------------------------------------- Tests ------------------------------------------- */

func TestInstallStyleToolsScriptUsesExistingTools(t *testing.T) {
	harness := newScriptHarness(t, "entrypoints/install-tools.sh")
	harness.writeProxyCommand(t, "awk")
	harness.writeProxyCommand(t, "sort")
	harness.writeProxyCommand(t, "head")
	harness.writeProxyCommand(t, "mkdir")
	writeFakeGoCommand(t, harness, false)

	goBinDirectory := filepath.Join(harness.fakeGoPath, goBinRelativePath)
	localBinDirectory := filepath.Join(harness.fakeHomeDirectory, localBinRelativePath)
	for _, toolName := range requiredStyleToolNames {
		writeStubExecutable(t, filepath.Join(goBinDirectory, toolName))
	}
	writeStubExecutable(t, filepath.Join(localBinDirectory, "markdownlint"))

	environment := harness.env([]string{harness.fakeBinDirectory})
	output, err := runBashScriptWithEnv(
		harness.scriptPath("entrypoints/install-tools.sh"),
		environment,
	)
	if err != nil {
		t.Fatalf("expected installer to succeed with existing tools, output:\n%s", output)
	}

	if !strings.Contains(output, "Style tools installed.") {
		t.Fatalf("expected success message, got:\n%s", output)
	}
}

func TestInstallStyleToolsScriptReportsMissingToolsAfterInstallAttempt(t *testing.T) {
	harness := newScriptHarness(t, "entrypoints/install-tools.sh")
	harness.writeProxyCommand(t, "awk")
	harness.writeProxyCommand(t, "sort")
	harness.writeProxyCommand(t, "head")
	harness.writeProxyCommand(t, "mkdir")
	writeFakeGoCommand(t, harness, true)

	environment := harness.env([]string{harness.fakeBinDirectory})
	output, err := runBashScriptWithEnv(
		harness.scriptPath("entrypoints/install-tools.sh"),
		environment,
	)
	if err == nil {
		t.Fatalf("expected installer to fail when tools remain missing, output:\n%s", output)
	}

	if !strings.Contains(output, "missing required style tools:") {
		t.Fatalf("expected missing-tools message, got:\n%s", output)
	}

	if !strings.Contains(output, "markdownlint") {
		t.Fatalf("expected missing markdownlint in output, got:\n%s", output)
	}
}

/* ------------------------------------------- Helpers ------------------------------------------ */

func writeStubExecutable(t *testing.T, path string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatalf("mkdir stub executable dir: %v", err)
	}

	contents := "#!/bin/bash\nset -euo pipefail\n"
	contents += "tool_name=\"${0##*/}\"\n"
	contents += "case \"$tool_name\" in\n"
	contents += "golangci-lint)\n"
	contents += "\tif [ \"${1:-}\" = \"version\" ]; then\n"
	contents += "\t\techo \"golangci-lint has version 2.6.2\"\n"
	contents += "\tfi\n"
	contents += "\t;;\n"
	contents += "shfmt)\n"
	contents += "\tif [ \"${1:-}\" = \"--version\" ]; then\n"
	contents += "\t\techo \"v3.12.0\"\n"
	contents += "\tfi\n"
	contents += "\t;;\n"
	contents += "shellcheck)\n"
	contents += "\tif [ \"${1:-}\" = \"--version\" ]; then\n"
	contents += "\t\techo \"version: 0.10.0\"\n"
	contents += "\tfi\n"
	contents += "\t;;\n"
	contents += "rg)\n"
	contents += "\tif [ \"${1:-}\" = \"--version\" ]; then\n"
	contents += "\t\techo \"ripgrep 14.1.0\"\n"
	contents += "\tfi\n"
	contents += "\t;;\n"
	contents += "markdownlint)\n"
	contents += "\tif [ \"${1:-}\" = \"--version\" ]; then\n"
	contents += "\t\techo \"0.45.0\"\n"
	contents += "\tfi\n"
	contents += "\t;;\n"
	contents += "esac\n"
	contents += "exit 0\n"
	if err := os.WriteFile(path, []byte(contents), 0o700); err != nil {
		t.Fatalf("write stub executable %s: %v", path, err)
	}
}
