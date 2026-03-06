package style_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

/* -------------------------------------------- Tests ------------------------------------------- */

func TestCheckStyleScriptWarnsOnRecommendationFindings(t *testing.T) {
	harness := newScriptHarness(t, "entrypoints/check-style.sh")
	harness.writeProxyCommand(t, "bash")
	harness.writeProxyCommand(t, "dirname")
	harness.writeProxyCommand(t, "mkdir")
	writeFakeGoCommand(t, harness, false)

	logPath := filepath.Join(harness.projectRoot, "style.log")
	harness.writeScript(
		t,
		requiredCheckScriptPath,
		"#!/bin/bash\nset -euo pipefail\necho \"required:$*\" >> \"$STYLE_TEST_LOG\"\n",
	)
	harness.writeScript(
		t,
		appOnlyCheckScriptPath,
		"#!/bin/bash\nset -euo pipefail\necho \"app:$*\" >> \"$STYLE_TEST_LOG\"\nexit 1\n",
	)
	harness.writeScript(
		t,
		recommendationCheckScriptPath,
		"#!/bin/bash\nset -euo pipefail\necho \"recommend:$*\" >> \"$STYLE_TEST_LOG\"\n"+
			"echo \"fake recommendation\"\nexit 1\n",
	)

	registryTable := writeRegistryTable(
		t,
		registryTableLines(
			requiredRegistryRow(
				registryTierTwo,
				"9.1",
				"Required tools check",
				registryScopeTools,
				registryRunnerScriptScope,
				requiredCheckScriptPath,
			),
			requiredRegistryRow(
				registryTierTwo,
				"9.2",
				"App only check",
				registryScopeApp,
				registryRunnerScriptScope,
				appOnlyCheckScriptPath,
			),
			recommendationRegistryRow(
				registryTierTwo,
				"R9.1",
				"Optional tools check",
				registryScopeAll,
				registryRunnerScriptScope,
				recommendationCheckScriptPath,
			),
		),
	)

	environment := harness.env(
		[]string{harness.fakeBinDirectory},
		styleRegistryTableEnvName+"="+registryTable,
		styleTestLogEnvName+"="+logPath,
	)
	output, err := runBashScriptWithEnv(
		harness.scriptPath("entrypoints/check-style.sh"),
		environment,
		"--scope",
		registryScopeTools,
		"--profile",
		styleProfileAll,
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
	harness := newScriptHarness(t, "entrypoints/check-style.sh")
	harness.writeProxyCommand(t, "bash")
	harness.writeProxyCommand(t, "dirname")
	harness.writeProxyCommand(t, "mkdir")
	writeFakeGoCommand(t, harness, false)

	harness.writeScript(
		t,
		requiredCheckScriptPath,
		"#!/bin/bash\nset -euo pipefail\nexit 0\n",
	)
	harness.writeScript(
		t,
		recommendationCheckScriptPath,
		"#!/bin/bash\nset -euo pipefail\necho \"fake recommendation\"\nexit 1\n",
	)

	registryTable := writeRegistryTable(
		t,
		registryTableLines(
			requiredRegistryRow(
				registryTierTwo,
				"9.1",
				"Required tools check",
				registryScopeAll,
				registryRunnerScriptScope,
				requiredCheckScriptPath,
			),
			recommendationRegistryRow(
				registryTierTwo,
				"R9.1",
				"Optional tools check",
				registryScopeAll,
				registryRunnerScriptScope,
				recommendationCheckScriptPath,
			),
		),
	)

	environment := harness.env(
		[]string{harness.fakeBinDirectory},
		styleRegistryTableEnvName+"="+registryTable,
	)
	output, err := runBashScriptWithEnv(
		harness.scriptPath("entrypoints/check-style.sh"),
		environment,
		"--profile",
		styleProfileAll,
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
	harness := newScriptHarness(t, "entrypoints/check-style.sh")
	harness.writeProxyCommand(t, "dirname")
	harness.writeProxyCommand(t, "mkdir")
	writeFakeGoCommand(t, harness, false)

	if err := os.MkdirAll(
		filepath.Join(harness.projectRoot, "tools", "style", "ast"),
		0o700,
	); err != nil {
		t.Fatalf("mkdir style ast module: %v", err)
	}

	logPath := filepath.Join(harness.projectRoot, "executor.log")
	harness.writeFakeCommand(
		t,
		"golangci-lint",
		"#!/bin/bash\nset -euo pipefail\necho \"executor:$*|pwd:$PWD\" >> \"$STYLE_TEST_LOG\"\n",
	)

	registryTable := writeRegistryTable(
		t,
		registryTableLines(
			requiredRegistryRow(
				registryTierOne,
				"1-2",
				"Executor check",
				registryScopeTools,
				registryRunnerExecutor,
				registryTargetGolangciTools,
			),
		),
	)

	environment := harness.env(
		[]string{harness.fakeBinDirectory},
		styleRegistryTableEnvName+"="+registryTable,
		styleTestLogEnvName+"="+logPath,
	)
	output, err := runBashScriptWithEnv(
		harness.scriptPath("entrypoints/check-style.sh"),
		environment,
		"--scope",
		registryScopeTools,
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

	if !strings.Contains(logOutput, "/tools/style/ast") {
		t.Fatalf("expected executor to run inside tools/style/ast, got log:\n%s", logOutput)
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
