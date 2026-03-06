package style_test

import "testing"

/* -------------------------------------------- Tests ------------------------------------------- */

func TestStyleRegistryValidateDetectsInconsistentArrays(t *testing.T) {
	t.Parallel()

	output, err := runStyleRegistryCommand(
		"",
		[]string{
			"set -euo pipefail",
			"source \"$STYLE_DIR/internal/common.sh\"",
			"source \"$STYLE_DIR/internal/registry.sh\"",
			"style_register_default_checks",
			"unset 'CHECK_TARGETS[0]'",
			"if style_registry_validate; then",
			"\techo \"expected validation failure\" >&2",
			"\texit 1",
			"fi",
		},
	)
	if err != nil {
		t.Fatalf("expected validation mismatch to be detected, output:\n%s", output)
	}
}

func TestStyleRegistryRejectsMalformedRows(t *testing.T) {
	t.Parallel()

	malformedTableFile := writeRegistryTable(
		t,
		registryTableLines(
			"tier1|required|1-2|broken row|app|runner",
		),
	)

	assertRegistryLoadRejected(t, malformedTableFile, "invalid style check registry row:")
}

func TestStyleRegistryRejectsEmptyTable(t *testing.T) {
	t.Parallel()

	tableFile := writeRegistryTable(
		t,
		registryTableLines(),
	)

	assertRegistryLoadRejected(t, tableFile, "style check registry is inconsistent")
}

func TestStyleRegistryRejectsInvalidLevel(t *testing.T) {
	t.Parallel()

	tableFile := writeRegistryTable(
		t,
		registryTableLines(
			registryTableRow(
				registryTierOne,
				"required-ish",
				"1-2",
				"golangci-lint (app)",
				registryScopeApp,
				registryRunnerExecutor,
				registryTargetGolangciApp,
			),
		),
	)

	assertRegistryLoadRejected(t, tableFile, "invalid style check registry row:")
}

func TestStyleRegistryRejectsRequiredRecommendationRule(t *testing.T) {
	t.Parallel()

	tableFile := writeRegistryTable(
		t,
		registryTableLines(
			requiredRegistryRow(
				registryTierTwo,
				"R2.10",
				"Magic values (Bash)",
				registryScopeAll,
				registryRunnerScriptScope,
				checkPathBashMagicValues,
			),
		),
	)

	assertRegistryLoadRejected(t, tableFile, "invalid style check registry row:")
}

func TestStyleRegistryRejectsRecommendationRequiredRule(t *testing.T) {
	t.Parallel()

	tableFile := writeRegistryTable(
		t,
		registryTableLines(
			recommendationRegistryRow(
				registryTierTwo,
				"2.10",
				"No magic values (Go)",
				registryScopeAll,
				registryRunnerScriptScope,
				checkPathGoMagicValues,
			),
		),
	)

	assertRegistryLoadRejected(t, tableFile, "invalid style check registry row:")
}

func TestStyleRegistryRejectsInvalidRunnerTargetPair(t *testing.T) {
	t.Parallel()

	tableFile := writeRegistryTable(
		t,
		registryTableLines(
			requiredRegistryRow(
				registryTierTwo,
				"1.1",
				"Line length (Go)",
				registryScopeAll,
				registryRunnerScriptScope,
				"lll",
			),
		),
	)

	assertRegistryLoadRejected(t, tableFile, "invalid style check registry row:")
}

func TestStyleRegistryRejectsUnknownExecutorTarget(t *testing.T) {
	t.Parallel()

	tableFile := writeRegistryTable(
		t,
		registryTableLines(
			requiredRegistryRow(
				registryTierOne,
				"1-2",
				"golangci-lint (app)",
				registryScopeApp,
				registryRunnerExecutor,
				"golangci_unknown",
			),
		),
	)

	assertRegistryLoadRejected(t, tableFile, "invalid style check registry row:")
}

func TestStyleRegistryRejectsMissingScriptTarget(t *testing.T) {
	t.Parallel()

	tableFile := writeRegistryTable(
		t,
		registryTableLines(
			requiredRegistryRow(
				registryTierTwo,
				"2.11",
				"Bash style",
				registryScopeAll,
				registryRunnerScriptScope,
				missingCheckScriptPath,
			),
		),
	)

	assertRegistryLoadRejected(t, tableFile, "invalid style check registry row:")
}
