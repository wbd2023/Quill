package scripts_test

import (
	"os"
	"strings"
	"testing"
)

/* -------------------------------------------- Tests ------------------------------------------- */

func TestCheckBashLineLengthScript(t *testing.T) {
	harness := newScriptHarness(t, checkPathBashLineLength)
	targetScript := harness.scriptPath(checkPathBashLineLength)

	harness.writeProjectFile(
		t,
		"tools/short.sh",
		"#!/bin/bash\nset -euo pipefail\necho \"short\"\n",
		0o600,
	)

	if output, err := runBashScript(targetScript, "--scope", "tools"); err != nil {
		t.Fatalf("expected short script to pass, output:\n%s", output)
	}

	harness.writeProjectFile(
		t,
		"tools/tab-width.sh",
		"#!/bin/bash\nset -euo pipefail\n\t"+strings.Repeat("a", 97)+"\n",
		0o600,
	)

	output, err := runBashScript(targetScript, "--scope", "tools")
	if err == nil {
		t.Fatalf("expected tab-width script to fail, output:\n%s", output)
	}

	if !strings.Contains(output, "[1.1] Bash line exceeds 100 columns") {
		t.Fatalf("expected tab-width line length finding, got:\n%s", output)
	}

	harness.writeProjectFile(
		t,
		"tools/tab-width.sh",
		"#!/bin/bash\nset -euo pipefail\n\t"+strings.Repeat("a", 97)+
			" # style: allow-long-line\n",
		0o600,
	)

	if output, err := runBashScript(targetScript, "--scope", "tools"); err != nil {
		t.Fatalf("expected marked tab-width line to pass, output:\n%s", output)
	}

	longPayload := strings.Repeat("a", 120)
	harness.writeProjectFile(
		t,
		"tools/long.sh",
		"#!/bin/bash\nset -euo pipefail\necho \""+longPayload+"\"\n",
		0o600,
	)

	output, err = runBashScript(targetScript, "--scope", "tools")
	if err == nil {
		t.Fatalf("expected long script to fail, output:\n%s", output)
	}

	if !strings.Contains(output, "[1.1] Bash line exceeds 100 columns") {
		t.Fatalf("expected line length finding, got:\n%s", output)
	}

	harness.writeProjectFile(
		t,
		"tools/long.sh",
		"#!/bin/bash\nset -euo pipefail\necho \""+longPayload+"\" # style: allow-long-line\n",
		0o600,
	)

	if output, err := runBashScript(targetScript, "--scope", "tools"); err != nil {
		t.Fatalf("expected marked long line to pass, output:\n%s", output)
	}
}

func TestCheckBashLineLengthScriptUsage(t *testing.T) {
	harness := newScriptHarness(t, checkPathBashLineLength)
	targetScript := harness.scriptPath(checkPathBashLineLength)

	output, err := runBashScript(targetScript, "--bad-flag")
	if err == nil {
		t.Fatalf("expected invalid arguments to fail, output:\n%s", output)
	}

	if !strings.Contains(output, "usage:") {
		t.Fatalf("expected usage output for invalid arguments, got:\n%s", output)
	}

	output, err = runBashScript(targetScript, "--scope")
	if err == nil {
		t.Fatalf("expected missing scope value to fail, output:\n%s", output)
	}

	if !strings.Contains(output, "usage:") {
		t.Fatalf("expected usage output for missing scope value, got:\n%s", output)
	}
}

func TestCheckGoLineLengthScriptTabWidth(t *testing.T) {
	harness := newScriptHarness(t, checkPathGoLineLength)
	targetScript := harness.scriptPath(checkPathGoLineLength)

	fakeLinterContent := "#!/bin/bash\nset -euo pipefail\n"
	fakeLinterContent += "if [ \"${1:-}\" = \"run\" ] && [ \"${2:-}\" = \"--help\" ]; then\n"
	fakeLinterContent += "\techo \"--enable-only\"\n\texit 0\nfi\n"
	fakeLinterContent += "if [ \"${1:-}\" = \"run\" ]; then\n\texit 0\nfi\nexit 0\n"
	harness.writeFakeCommand(t, "golangci-lint", fakeLinterContent)

	harness.writeProjectFile(
		t,
		"tools/tab_width.go",
		"package tools\n\nfunc tabWidth() {\n\t// "+strings.Repeat("a", 97)+"\n}\n",
		0o600,
	)

	harness.writeProjectFile(
		t,
		"tools/tab_width_test.go",
		"package tools\n\nfunc TestTabWidth() {\n\t// "+strings.Repeat("a", 97)+"\n}\n",
		0o600,
	)

	environment := harness.env([]string{harness.fakeBinDirectory, os.Getenv("PATH")})
	output, err := runBashScriptWithEnv(targetScript, environment, "--scope", "tools")
	if err == nil {
		t.Fatalf("expected tab-width go file to fail, output:\n%s", output)
	}

	if !strings.Contains(output, "[1.1] Go line exceeds 100 columns") {
		t.Fatalf("expected tab-width finding, got:\n%s", output)
	}

	harness.writeProjectFile(
		t,
		"tools/tab_width.go",
		"package tools\n\nfunc tabWidth() {\n\t// short\n}\n",
		0o600,
	)

	harness.writeProjectFile(
		t,
		"tools/tab_width_test.go",
		"package tools\n\nfunc TestTabWidth() {\n\t// short\n}\n",
		0o600,
	)

	if output, err := runBashScriptWithEnv(
		targetScript,
		environment,
		"--scope",
		"tools",
	); err != nil {
		t.Fatalf("expected fixed go file to pass, output:\n%s", output)
	}
}
