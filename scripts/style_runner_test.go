package scripts_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

/* -------------------------------------------- Tests ------------------------------------------- */

func TestCheckStyleScriptWarnsOnRecommendationFindings(t *testing.T) {
	harness := newScriptHarness(t, "check-style.sh")
	harness.writeProxyCommand(t, "bash")
	harness.writeProxyCommand(t, "dirname")
	harness.writeProxyCommand(t, "mkdir")
	writeFakeGoCommand(t, harness, false)

	logPath := filepath.Join(harness.projectRoot, "style.log")
	harness.writeScript(
		t,
		"check-required.sh",
		"#!/bin/bash\nset -euo pipefail\necho \"required:$*\" >> \"$STYLE_TEST_LOG\"\n",
	)
	harness.writeScript(
		t,
		"check-app-only.sh",
		"#!/bin/bash\nset -euo pipefail\necho \"app:$*\" >> \"$STYLE_TEST_LOG\"\nexit 1\n",
	)
	harness.writeScript(
		t,
		"check-recommend.sh",
		"#!/bin/bash\nset -euo pipefail\necho \"recommend:$*\" >> \"$STYLE_TEST_LOG\"\n"+
			"echo \"fake recommendation\"\nexit 1\n",
	)

	registryTable := writeRegistryTable(
		t,
		[]string{
			"# tier|level|rule|name|scope|runner|target",
			"tier2|required|9.1|Required tools check|tools|script_scope|check-required.sh",
			"tier2|required|9.2|App only check|app|script_scope|check-app-only.sh",
			"tier2|recommendation|R9.1|Optional tools check|all|script_scope|check-recommend.sh",
		},
	)

	environment := harness.env(
		[]string{harness.fakeBinDirectory},
		"STYLE_REGISTRY_TABLE_FILE="+registryTable,
		"STYLE_TEST_LOG="+logPath,
	)
	output, err := runBashScriptWithEnv(
		harness.scriptPath("check-style.sh"),
		environment,
		"--scope",
		"tools",
		"--profile",
		"all",
	)
	if err != nil {
		t.Fatalf("expected warning-only style run to pass, output:\n%s", output)
	}

	if !strings.Contains(output, "WARN") {
		t.Fatalf("expected warning result in output, got:\n%s", output)
	}

	if !strings.Contains(output, "Required checks passed with recommendations.") {
		t.Fatalf("expected warning summary, got:\n%s", output)
	}

	loggedRuns, readErr := os.ReadFile(logPath)
	if readErr != nil {
		t.Fatalf("read style log: %v", readErr)
	}

	logOutput := string(loggedRuns)
	if !strings.Contains(logOutput, "required:--scope tools") {
		t.Fatalf("expected tools check execution, got log:\n%s", logOutput)
	}

	if !strings.Contains(logOutput, "recommend:--scope tools") {
		t.Fatalf("expected recommendation execution, got log:\n%s", logOutput)
	}

	if strings.Contains(logOutput, "app:") {
		t.Fatalf("app-scoped check should not run for tools scope, got log:\n%s", logOutput)
	}
}

func TestCheckStyleScriptFailsStrictRecommendations(t *testing.T) {
	harness := newScriptHarness(t, "check-style.sh")
	harness.writeProxyCommand(t, "bash")
	harness.writeProxyCommand(t, "dirname")
	harness.writeProxyCommand(t, "mkdir")
	writeFakeGoCommand(t, harness, false)

	harness.writeScript(
		t,
		"check-required.sh",
		"#!/bin/bash\nset -euo pipefail\nexit 0\n",
	)
	harness.writeScript(
		t,
		"check-recommend.sh",
		"#!/bin/bash\nset -euo pipefail\necho \"fake recommendation\"\nexit 1\n",
	)

	registryTable := writeRegistryTable(
		t,
		[]string{
			"# tier|level|rule|name|scope|runner|target",
			"tier2|required|9.1|Required tools check|all|script_scope|check-required.sh",
			"tier2|recommendation|R9.1|Optional tools check|all|script_scope|check-recommend.sh",
		},
	)

	environment := harness.env(
		[]string{harness.fakeBinDirectory},
		"STYLE_REGISTRY_TABLE_FILE="+registryTable,
	)
	output, err := runBashScriptWithEnv(
		harness.scriptPath("check-style.sh"),
		environment,
		"--profile",
		"all",
		"--strict-recommendations",
	)
	if err == nil {
		t.Fatalf("expected strict recommendation failure, output:\n%s", output)
	}

	if !strings.Contains(output, "FAIL") {
		t.Fatalf("expected failure result in output, got:\n%s", output)
	}

	if !strings.Contains(output, "Some checks failed") {
		t.Fatalf("expected failure summary, got:\n%s", output)
	}
}

func TestCheckStyleScriptRunsExecutorTargets(t *testing.T) {
	harness := newScriptHarness(t, "check-style.sh")
	harness.writeProxyCommand(t, "dirname")
	harness.writeProxyCommand(t, "mkdir")
	writeFakeGoCommand(t, harness, false)

	if err := os.MkdirAll(
		filepath.Join(harness.projectRoot, "tools", "stylecheck"),
		0o700,
	); err != nil {
		t.Fatalf("mkdir stylecheck module: %v", err)
	}

	logPath := filepath.Join(harness.projectRoot, "executor.log")
	harness.writeFakeCommand(
		t,
		"golangci-lint",
		"#!/bin/bash\nset -euo pipefail\necho \"executor:$*|pwd:$PWD\" >> \"$STYLE_TEST_LOG\"\n",
	)

	registryTable := writeRegistryTable(
		t,
		[]string{
			"# tier|level|rule|name|scope|runner|target",
			"tier1|required|1-2|Executor check|tools|runner|golangci_tools",
		},
	)

	environment := harness.env(
		[]string{harness.fakeBinDirectory},
		"STYLE_REGISTRY_TABLE_FILE="+registryTable,
		"STYLE_TEST_LOG="+logPath,
	)
	output, err := runBashScriptWithEnv(
		harness.scriptPath("check-style.sh"),
		environment,
		"--scope",
		"tools",
	)
	if err != nil {
		t.Fatalf("expected executor-backed style run to pass, output:\n%s", output)
	}

	loggedRuns, readErr := os.ReadFile(logPath)
	if readErr != nil {
		t.Fatalf("read executor log: %v", readErr)
	}

	logOutput := string(loggedRuns)
	if !strings.Contains(logOutput, "executor:run ./...") {
		t.Fatalf("expected golangci-lint execution, got log:\n%s", logOutput)
	}

	if !strings.Contains(logOutput, "/tools/stylecheck") {
		t.Fatalf("expected executor to run inside tools/stylecheck, got log:\n%s", logOutput)
	}
}

/* ------------------------------------------- Helpers ------------------------------------------ */

func writeFakeGoCommand(t *testing.T, harness scriptHarness, allowInstall bool) {
	t.Helper()

	goContents := "#!/bin/bash\nset -euo pipefail\n"
	goContents += "if [ \"${1:-}\" = \"env\" ] && [ \"${2:-}\" = \"GOPATH\" ]; then\n"
	goContents += "\techo \"$FAKE_GOPATH\"\n\texit 0\nfi\n"
	if allowInstall {
		goContents += "if [ \"${1:-}\" = \"install\" ]; then\n\texit 0\nfi\n"
	}
	goContents += "echo \"unexpected go command: $*\" >&2\nexit 1\n"

	harness.writeFakeCommand(t, "go", goContents)
}
