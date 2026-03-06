package style_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

/* -------------------------------------- Registry Helpers -------------------------------------- */

func loadStyleRegistryRows(tableOverride string) (rows []styleRegistryRow, err error) {
	output, err := runStyleRegistryCommand(
		tableOverride,
		[]string{
			"set -euo pipefail",
			"source \"$STYLE_DIR/internal/common.sh\"",
			"source \"$STYLE_DIR/internal/registry.sh\"",
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
		"STYLE_DIR="+currentStyleDirectory(),
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
			"source \"$STYLE_DIR/internal/common.sh\"",
			"source \"$STYLE_DIR/internal/registry.sh\"",
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

	tableFile = filepath.Join(t.TempDir(), "registry.table")
	contents := strings.Join(lines, "\n") + "\n"
	if err := os.WriteFile(tableFile, []byte(contents), 0o600); err != nil {
		t.Fatalf("write registry table %s: %v", tableFile, err)
	}

	return tableFile
}
