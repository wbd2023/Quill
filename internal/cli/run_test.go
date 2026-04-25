package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

/* --------------------------------------------- Run -------------------------------------------- */

func TestRunRejectsMissingCommand(t *testing.T) {
	tool, stdout, stderr := newTestCLI()

	exitCode := tool.Run(nil)
	if exitCode != usageExitCode {
		t.Fatalf("expected usage exit code, got %d", exitCode)
	}

	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout for missing command, got %q", stdout.String())
	}

	if !strings.Contains(stderr.String(), "style <command> [flags]") {
		t.Fatalf("expected root usage on stderr, got %q", stderr.String())
	}
}

func TestRunRejectsUnknownCommand(t *testing.T) {
	tool, _, stderr := newTestCLI()

	exitCode := tool.Run([]string{"unknown"})
	if exitCode != usageExitCode {
		t.Fatalf("expected usage exit code, got %d", exitCode)
	}

	if !strings.Contains(stderr.String(), `unknown command "unknown"`) {
		t.Fatalf("expected unknown-command error, got %q", stderr.String())
	}
}

func TestRunTreatsRootHelpAsSuccess(t *testing.T) {
	tool, stdout, stderr := newTestCLI()

	exitCode := tool.Run([]string{"help"})
	if exitCode != 0 {
		t.Fatalf("expected success exit code for help, got %d", exitCode)
	}

	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr for help, got %q", stderr.String())
	}

	if stdout.String() != readGoldenOutput(t, "root_help.txt") {
		t.Fatalf("unexpected root help output:\n%s", stdout.String())
	}
}

func TestRunTreatsCommandHelpAsSuccess(t *testing.T) {
	tool, stdout, stderr := newTestCLI()

	exitCode := tool.Run([]string{"help", "check"})
	if exitCode != 0 {
		t.Fatalf("expected success exit code for command help, got %d", exitCode)
	}

	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr for command help, got %q", stderr.String())
	}

	if stdout.String() != readGoldenOutput(t, "check_help.txt") {
		t.Fatalf("unexpected command help output:\n%s", stdout.String())
	}
}

func TestRunTreatsFlagHelpAsSuccess(t *testing.T) {
	tool, stdout, stderr := newTestCLI()

	exitCode := tool.Run([]string{"check", "-h"})
	if exitCode != 0 {
		t.Fatalf("expected success exit code for flag help, got %d", exitCode)
	}

	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr for flag help, got %q", stderr.String())
	}

	if stdout.String() != readGoldenOutput(t, "check_help.txt") {
		t.Fatalf("unexpected flag help output:\n%s", stdout.String())
	}
}

/* -------------------------------------------- Usage ------------------------------------------- */

func TestUsageTextListsCommands(t *testing.T) {
	usage := rootUsageText()
	requiredSnippets := []string{
		"style <command> [flags]",
		"check",
		"fix",
		"doctor",
		"coverage",
		"install",
	}

	for _, snippet := range requiredSnippets {
		if !strings.Contains(usage, snippet) {
			t.Fatalf("usage text missing %q:\n%s", snippet, usage)
		}
	}
}

/* ---------------------------------------- Test Harness ---------------------------------------- */

func newTestCLI() (tool CLI, stdout *bytes.Buffer, stderr *bytes.Buffer) {
	stdout = &bytes.Buffer{}
	stderr = &bytes.Buffer{}
	tool = New(stdout, stderr)
	return tool, stdout, stderr
}

/* ---------------------------------------- Golden Output --------------------------------------- */

func readGoldenOutput(t *testing.T, name string) (output string) {
	t.Helper()

	path := filepath.Join("testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden output %q: %v", name, err)
	}

	return string(data)
}
