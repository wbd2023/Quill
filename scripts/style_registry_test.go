package scripts_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

/* -------------------------------------------- Tests ------------------------------------------- */

func TestStyleRegistryDefaultRows(t *testing.T) {
	t.Parallel()

	rows, err := loadStyleRegistryRows("")
	if err != nil {
		t.Fatalf("load default style registry rows: %v", err)
	}

	expectedRows := expectedStyleRegistryRows()
	if len(rows) != len(expectedRows) {
		t.Fatalf("row count mismatch: got %d, want %d", len(rows), len(expectedRows))
	}

	for index := range expectedRows {
		if rows[index] != expectedRows[index] {
			t.Fatalf(
				"row %d mismatch:\n got:  %#v\n want: %#v",
				index,
				rows[index],
				expectedRows[index],
			)
		}
	}

	seen := map[styleRegistryRow]struct{}{}
	for _, row := range rows {
		if _, exists := seen[row]; exists {
			t.Fatalf("duplicate registry row detected: %#v", row)
		}
		seen[row] = struct{}{}
	}
}

func TestStyleRegistryValidateDetectsInconsistentArrays(t *testing.T) {
	t.Parallel()

	output, err := runStyleRegistryCommand(
		"",
		[]string{
			"set -euo pipefail",
			"source \"$SCRIPTS_DIR/lib/style-common.sh\"",
			"source \"$SCRIPTS_DIR/lib/style-registry.sh\"",
			"style_register_default_checks",
			"unset 'CHECK_TARGETS[0]'",
			"if style_registry_validate; then",
			"	echo \"expected validation failure\" >&2",
			"	exit 1",
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
		[]string{
			"# tier|level|rule|name|scope|runner|target",
			"tier1|required|1-2|broken row|app|runner",
		},
	)

	assertRegistryLoadRejected(t, malformedTableFile, "invalid style check registry row:")
}

func TestStyleRegistryRejectsEmptyTable(t *testing.T) {
	t.Parallel()

	tableFile := writeRegistryTable(
		t,
		[]string{
			"# tier|level|rule|name|scope|runner|target",
		},
	)

	assertRegistryLoadRejected(t, tableFile, "style check registry is inconsistent")
}

func TestStyleRegistryRejectsInvalidLevel(t *testing.T) {
	t.Parallel()

	tableFile := writeRegistryTable(
		t,
		[]string{
			"# tier|level|rule|name|scope|runner|target",
			"tier1|required-ish|1-2|golangci-lint (app)|app|runner|golangci_app",
		},
	)

	assertRegistryLoadRejected(t, tableFile, "invalid style check registry row:")
}

func TestStyleRegistryRejectsRequiredRecommendationRule(t *testing.T) {
	t.Parallel()

	tableFile := writeRegistryTable(
		t,
		[]string{
			"# tier|level|rule|name|scope|runner|target",
			"tier2|required|R2.10|Magic values (Bash)|all|script_scope|check-bash-magic-values.sh",
		},
	)

	assertRegistryLoadRejected(t, tableFile, "invalid style check registry row:")
}

func TestStyleRegistryRejectsRecommendationRequiredRule(t *testing.T) {
	t.Parallel()

	tableFile := writeRegistryTable(
		t,
		[]string{
			"# tier|level|rule|name|scope|runner|target",
			"tier2|recommendation|2.10|No magic values (Go)|all|script_scope|" +
				"check-go-magic-values.sh",
		},
	)

	assertRegistryLoadRejected(t, tableFile, "invalid style check registry row:")
}

func TestStyleRegistryRejectsInvalidRunnerTargetPair(t *testing.T) {
	t.Parallel()

	tableFile := writeRegistryTable(
		t,
		[]string{
			"# tier|level|rule|name|scope|runner|target",
			"tier2|required|1.1|Line length (Go)|all|script_scope|lll",
		},
	)

	assertRegistryLoadRejected(t, tableFile, "invalid style check registry row:")
}

func TestStyleRegistryRejectsUnknownExecutorTarget(t *testing.T) {
	t.Parallel()

	tableFile := writeRegistryTable(
		t,
		[]string{
			"# tier|level|rule|name|scope|runner|target",
			"tier1|required|1-2|golangci-lint (app)|app|runner|golangci_unknown",
		},
	)

	assertRegistryLoadRejected(t, tableFile, "invalid style check registry row:")
}

func TestStyleRegistryRejectsMissingScriptTarget(t *testing.T) {
	t.Parallel()

	tableFile := writeRegistryTable(
		t,
		[]string{
			"# tier|level|rule|name|scope|runner|target",
			"tier2|required|2.11|Bash style|all|script_scope|check-not-here.sh",
		},
	)

	assertRegistryLoadRejected(t, tableFile, "invalid style check registry row:")
}

/* ------------------------------------------- Helpers ------------------------------------------ */

const styleRegistryFieldCount = 7

type styleRegistryRow struct {
	Tier   string
	Level  string
	Rule   string
	Name   string
	Scope  string
	Runner string
	Target string
}

func expectedStyleRegistryRows() []styleRegistryRow {
	row := func(
		tier string,
		level string,
		rule string,
		name string,
		scope string,
		runner string,
		target string,
	) styleRegistryRow {
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

	return []styleRegistryRow{
		row("tier1", "required", "1-2", "golangci-lint (app)", "app", "runner", "golangci_app"),
		row(
			"tier1",
			"required",
			"1-2",
			"golangci-lint (tools/stylecheck)",
			"tools",
			"runner",
			"golangci_tools",
		),
		row(
			"tier2",
			"required",
			"1.1",
			"Line length (Go)",
			"all",
			"script_scope",
			"check-go-line-length.sh",
		),
		row(
			"tier2",
			"required",
			"1.1",
			"Line length (Bash scripts)",
			"all",
			"script_scope",
			"check-bash-line-length.sh",
		),
		row(
			"tier2",
			"required",
			"1.2",
			"Language spelling (non-Go, UK locale)",
			"all",
			"script_scope",
			"check-spelling.sh",
		),
		row(
			"tier2",
			"required",
			"3.2",
			"Markdown style (markdownlint)",
			"all",
			"script_scope",
			"check-markdown.sh",
		),
		row(
			"tier2",
			"required",
			"1.3",
			"Architecture imports",
			"app",
			"script",
			"check-go-architecture-imports.sh",
		),
		row(
			"tier2",
			"required",
			"1.4",
			"ASCII-only characters",
			"all",
			"script_scope",
			"check-ascii.sh",
		),
		row(
			"tier2",
			"required",
			"2.1",
			"Error handling",
			"app",
			"script",
			"check-go-error-style.sh",
		),
		row(
			"tier2",
			"required",
			"2.1",
			"Adapter error wrapping",
			"app",
			"script",
			"check-adapter-errors.sh",
		),
		row(
			"tier2",
			"required",
			"2.2",
			"Naming conventions",
			"all",
			"script_scope",
			"check-naming.sh",
		),
		row(
			"tier2",
			"required",
			"2.3",
			"Inline comment style",
			"app",
			"script",
			"check-go-inline-comments.sh",
		),
		row(
			"tier2",
			"required",
			"2.4",
			"Section header format",
			"all",
			"script_scope",
			"check-section-headers.sh",
		),
		row(
			"tier2",
			"required",
			"2.10",
			"No magic values (Go)",
			"all",
			"script_scope",
			"check-go-magic-values.sh",
		),
		row(
			"tier2",
			"required",
			"2.11",
			"Bash script style",
			"all",
			"script_scope",
			"check-bash-style.sh",
		),
		row(
			"tier2",
			"required",
			"2.11",
			"Bash static analysis (shellcheck)",
			"all",
			"script_scope",
			"check-shellcheck.sh",
		),
		row(
			"tier2",
			"required",
			"2.11",
			"Bash formatting (shfmt)",
			"all",
			"script_scope",
			"check-shfmt.sh",
		),
		row(
			"tier2",
			"recommendation",
			"R2.10",
			"Magic values (Bash)",
			"all",
			"script_scope",
			"check-bash-magic-values.sh",
		),
		row(
			"tier2",
			"recommendation",
			"R2.12",
			"Guard-clause spacing (Go)",
			"all",
			"script_scope",
			"check-vertical-spacing.sh",
		),
		row("tier3", "required", "2.2+", "AST rules (app)", "app", "runner", "ast_app"),
		row(
			"tier3",
			"required",
			"2.2+",
			"AST rules (tools/stylecheck)",
			"tools",
			"runner",
			"ast_tools",
		),
	}
}

func loadStyleRegistryRows(tableOverride string) ([]styleRegistryRow, error) {
	output, err := runStyleRegistryCommand(
		tableOverride,
		[]string{
			"set -euo pipefail",
			"source \"$SCRIPTS_DIR/lib/style-common.sh\"",
			"source \"$SCRIPTS_DIR/lib/style-registry.sh\"",
			"style_register_default_checks",
			"for index in \"${!CHECK_TIERS[@]}\"; do",
			"	printf '%s|%s|%s|%s|%s|%s|%s\\n' \\",
			"		\"${CHECK_TIERS[$index]}\" \\",
			"		\"${CHECK_LEVELS[$index]}\" \\",
			"		\"${CHECK_RULES[$index]}\" \\",
			"		\"${CHECK_NAMES[$index]}\" \\",
			"		\"${CHECK_SCOPES[$index]}\" \\",
			"		\"${CHECK_RUNNERS[$index]}\" \\",
			"		\"${CHECK_TARGETS[$index]}\"",
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
	rows := make([]styleRegistryRow, 0, len(lines))
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

func runStyleRegistryCommand(tableOverride string, scriptLines []string) (string, error) {
	command := exec.Command("bash", "-c", strings.Join(scriptLines, "\n"))
	command.Env = append(
		os.Environ(),
		"SCRIPTS_DIR="+currentScriptsDirectory(),
		"STYLE_REGISTRY_TABLE_FILE="+tableOverride,
	)

	rawOutput, runErr := command.CombinedOutput()
	return string(rawOutput), runErr
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
			"	echo \"expected registry load failure\" >&2",
			"	exit 1",
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
