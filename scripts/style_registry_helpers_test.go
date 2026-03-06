package scripts_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const (
	styleRegistryFieldCount = 7

	styleRegistryTableEnvName = "STYLE_REGISTRY_TABLE_FILE"
	styleTestLogEnvName       = "STYLE_TEST_LOG"

	registryTierOne   = "tier1"
	registryTierTwo   = "tier2"
	registryTierThree = "tier3"

	registryLevelRequired       = "required"
	registryLevelRecommendation = "recommendation"
	styleProfileAll             = "all"

	registryScopeApp   = "app"
	registryScopeTools = "tools"
	registryScopeAll   = "all"

	registryRunnerScript      = "script"
	registryRunnerScriptScope = "script_scope"
	registryRunnerExecutor    = "runner"

	registryTargetGolangciApp   = "golangci_app"
	registryTargetGolangciTools = "golangci_tools"
	registryTargetASTApp        = "ast_app"
	registryTargetASTTools      = "ast_tools"

	checkPathBashLineLength       = "checks/bash/line-length.sh"
	checkPathBashMagicValues      = "checks/bash/magic-values.sh"
	checkPathBashShellcheck       = "checks/bash/shellcheck.sh"
	checkPathBashShfmt            = "checks/bash/shfmt.sh"
	checkPathBashStyle            = "checks/bash/style.sh"
	checkPathGeneralASCII         = "checks/general/ascii.sh"
	checkPathGeneralMarkdown      = "checks/general/markdown.sh"
	checkPathGeneralNaming        = "checks/general/naming.sh"
	checkPathGeneralHeaders       = "checks/general/section-headers.sh"
	checkPathGeneralSpelling      = "checks/general/spelling.sh"
	checkPathGoArchitecture       = "checks/go/architecture-imports.sh"
	checkPathGoLineLength         = "checks/go/line-length.sh"
	checkPathGoMagicValues        = "checks/go/magic-values.sh"
	checkPathGoVerticalSpacing    = "checks/go/vertical-spacing.sh"
	missingCheckScriptPath        = "checks/test/missing.sh"
	requiredCheckScriptPath       = "checks/test/required.sh"
	appOnlyCheckScriptPath        = "checks/test/app-only.sh"
	recommendationCheckScriptPath = "checks/test/recommendation.sh"
)

type styleRegistryRow struct {
	Tier   string
	Level  string
	Rule   string
	Name   string
	Scope  string
	Runner string
	Target string
}

/* ----------------------------------------- Row Helpers ---------------------------------------- */

func newStyleRegistryRow(
	tier string,
	level string,
	rule string,
	name string,
	scope string,
	runner string,
	target string,
) (row styleRegistryRow) {
	return styleRegistryRow{
		Tier:   tier,
		Level:  level,
		Rule:   rule,
		Name:   name,
		Scope:  scope,
		Runner: runner,
		Target: target,
	}
}

func registryTableRow(
	tier string,
	level string,
	rule string,
	name string,
	scope string,
	runner string,
	target string,
) (registryRowText string) {
	fields := []string{tier, level, rule, name, scope, runner, target}
	return strings.Join(fields, "|")
}

func requiredRegistryRow(
	tier string,
	rule string,
	name string,
	scope string,
	runner string,
	target string,
) (registryRowText string) {
	return registryTableRow(tier, registryLevelRequired, rule, name, scope, runner, target)
}

func recommendationRegistryRow(
	tier string,
	rule string,
	name string,
	scope string,
	runner string,
	target string,
) (registryRowText string) {
	return registryTableRow(
		tier,
		registryLevelRecommendation,
		rule,
		name,
		scope,
		runner,
		target,
	)
}

func expectedStyleRegistryRows() (rows []styleRegistryRow) {
	rows = []styleRegistryRow{
		newStyleRegistryRow(
			registryTierOne,
			registryLevelRequired,
			"1-2",
			"golangci-lint (app)",
			registryScopeApp,
			registryRunnerExecutor,
			registryTargetGolangciApp,
		),
		newStyleRegistryRow(
			registryTierOne,
			registryLevelRequired,
			"1-2",
			"golangci-lint (tools/stylecheck)",
			registryScopeTools,
			registryRunnerExecutor,
			registryTargetGolangciTools,
		),
		newStyleRegistryRow(
			registryTierTwo,
			registryLevelRequired,
			"1.1",
			"Line length (Go)",
			registryScopeAll,
			registryRunnerScriptScope,
			checkPathGoLineLength,
		),
		newStyleRegistryRow(
			registryTierTwo,
			registryLevelRequired,
			"1.1",
			"Line length (Bash scripts)",
			registryScopeAll,
			registryRunnerScriptScope,
			checkPathBashLineLength,
		),
		newStyleRegistryRow(
			registryTierTwo,
			registryLevelRequired,
			"1.2",
			"Language spelling (non-Go, UK locale)",
			registryScopeAll,
			registryRunnerScriptScope,
			checkPathGeneralSpelling,
		),
		newStyleRegistryRow(
			registryTierTwo,
			registryLevelRequired,
			"3.2",
			"Markdown style (markdownlint)",
			registryScopeAll,
			registryRunnerScriptScope,
			checkPathGeneralMarkdown,
		),
		newStyleRegistryRow(
			registryTierTwo,
			registryLevelRequired,
			"1.3",
			"Architecture imports",
			registryScopeApp,
			registryRunnerScript,
			checkPathGoArchitecture,
		),
		newStyleRegistryRow(
			registryTierTwo,
			registryLevelRequired,
			"1.4",
			"ASCII-only characters",
			registryScopeAll,
			registryRunnerScriptScope,
			checkPathGeneralASCII,
		),
		newStyleRegistryRow(
			registryTierTwo,
			registryLevelRequired,
			"2.2",
			"Naming conventions",
			registryScopeAll,
			registryRunnerScriptScope,
			checkPathGeneralNaming,
		),
		newStyleRegistryRow(
			registryTierTwo,
			registryLevelRequired,
			"2.4",
			"Section header format",
			registryScopeAll,
			registryRunnerScriptScope,
			checkPathGeneralHeaders,
		),
		newStyleRegistryRow(
			registryTierTwo,
			registryLevelRequired,
			"2.10",
			"No magic values (Go)",
			registryScopeAll,
			registryRunnerScriptScope,
			checkPathGoMagicValues,
		),
		newStyleRegistryRow(
			registryTierTwo,
			registryLevelRequired,
			"2.11",
			"Bash script style",
			registryScopeAll,
			registryRunnerScriptScope,
			checkPathBashStyle,
		),
		newStyleRegistryRow(
			registryTierTwo,
			registryLevelRequired,
			"2.11",
			"Bash static analysis (shellcheck)",
			registryScopeAll,
			registryRunnerScriptScope,
			checkPathBashShellcheck,
		),
		newStyleRegistryRow(
			registryTierTwo,
			registryLevelRequired,
			"2.11",
			"Bash formatting (shfmt)",
			registryScopeAll,
			registryRunnerScriptScope,
			checkPathBashShfmt,
		),
		newStyleRegistryRow(
			registryTierTwo,
			registryLevelRecommendation,
			"R2.10",
			"Magic values (Bash)",
			registryScopeAll,
			registryRunnerScriptScope,
			checkPathBashMagicValues,
		),
		newStyleRegistryRow(
			registryTierTwo,
			registryLevelRecommendation,
			"R2.12",
			"Guard-clause spacing (Go)",
			registryScopeAll,
			registryRunnerScriptScope,
			checkPathGoVerticalSpacing,
		),
		newStyleRegistryRow(
			registryTierThree,
			registryLevelRequired,
			"2.2+",
			"AST rules (app)",
			registryScopeApp,
			registryRunnerExecutor,
			registryTargetASTApp,
		),
		newStyleRegistryRow(
			registryTierThree,
			registryLevelRequired,
			"2.2+",
			"AST rules (tools/stylecheck)",
			registryScopeTools,
			registryRunnerExecutor,
			registryTargetASTTools,
		),
	}

	return rows
}

/* -------------------------------------- Registry Helpers -------------------------------------- */

func loadStyleRegistryRows(tableOverride string) (rows []styleRegistryRow, err error) {
	output, err := runStyleRegistryCommand(
		tableOverride,
		[]string{
			"set -euo pipefail",
			"source \"$SCRIPTS_DIR/lib/style-common.sh\"",
			"source \"$SCRIPTS_DIR/lib/style-registry.sh\"",
			"style_register_default_checks",
			"for index in \"${!CHECK_TIERS[@]}\"; do",
			"\tprintf '%s|%s|%s|%s|%s|%s|%s\\n' \\",
			"\t\t\"${CHECK_TIERS[$index]}\" \\",
			"\t\t\"${CHECK_LEVELS[$index]}\" \\",
			"\t\t\"${CHECK_RULES[$index]}\" \\",
			"\t\t\"${CHECK_NAMES[$index]}\" \\",
			"\t\t\"${CHECK_SCOPES[$index]}\" \\",
			"\t\t\"${CHECK_RUNNERS[$index]}\" \\",
			"\t\t\"${CHECK_TARGETS[$index]}\"",
			"done",
		},
	)
	if err != nil {
		return nil, fmt.Errorf("run registry probe: %w (output: %s)", err, output)
	}

	if strings.TrimSpace(output) == "" {
		return nil, nil
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	rows = make([]styleRegistryRow, 0, len(lines))
	for lineNumber, line := range lines {
		fields := strings.Split(line, "|")
		if len(fields) != styleRegistryFieldCount {
			return nil, fmt.Errorf(
				"unexpected field count on line %d: got %d",
				lineNumber+1,
				len(fields),
			)
		}

		rows = append(rows, styleRegistryRow{
			Tier:   fields[0],
			Level:  fields[1],
			Rule:   fields[2],
			Name:   fields[3],
			Scope:  fields[4],
			Runner: fields[5],
			Target: fields[6],
		})
	}

	return rows, nil
}

func runStyleRegistryCommand(
	tableOverride string,
	scriptLines []string,
) (output string, err error) {
	command := exec.Command("bash", "-c", strings.Join(scriptLines, "\n"))
	command.Env = append(
		os.Environ(),
		"SCRIPTS_DIR="+currentScriptsDirectory(),
		styleRegistryTableEnvName+"="+tableOverride,
	)

	rawOutput, err := command.CombinedOutput()
	return string(rawOutput), err
}

func assertRegistryLoadRejected(t *testing.T, tableFile string, expectedError string) {
	t.Helper()

	output, err := runStyleRegistryCommand(
		tableFile,
		[]string{
			"set -euo pipefail",
			"source \"$SCRIPTS_DIR/lib/style-common.sh\"",
			"source \"$SCRIPTS_DIR/lib/style-registry.sh\"",
			"if style_register_default_checks; then",
			"\techo \"expected registry load failure\" >&2",
			"\texit 1",
			"fi",
		},
	)
	if err != nil {
		t.Fatalf("expected registry load to fail gracefully, output:\n%s", output)
	}

	if !strings.Contains(output, expectedError) {
		t.Fatalf("expected error %q, got:\n%s", expectedError, output)
	}
}

func writeRegistryTable(t *testing.T, lines []string) (tableFile string) {
	t.Helper()

	tableFile = filepath.Join(t.TempDir(), "style-registry.table")
	contents := strings.Join(lines, "\n") + "\n"
	if err := os.WriteFile(tableFile, []byte(contents), 0o600); err != nil {
		t.Fatalf("write registry table %s: %v", tableFile, err)
	}

	return tableFile
}
