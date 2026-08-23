package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

/* --------------------------------------------- Run -------------------------------------------- */

func TestRunRejectsMissingCommand(t *testing.T) {
	tool, stdout, stderr := newTestCLI()

	exitCode := tool.Run(context.Background(), nil)
	if exitCode != usageExitCode {
		t.Fatalf("expected usage exit code, got %d", exitCode)
	}

	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout for missing command, got %q", stdout.String())
	}

	if !strings.Contains(stderr.String(), "quill <command> [flags]") {
		t.Fatalf("expected root usage on stderr, got %q", stderr.String())
	}
}

func TestRunRejectsUnknownCommand(t *testing.T) {
	tool, _, stderr := newTestCLI()

	exitCode := tool.Run(context.Background(), []string{"unknown"})
	if exitCode != usageExitCode {
		t.Fatalf("expected usage exit code, got %d", exitCode)
	}

	if !strings.Contains(stderr.String(), "unexpected argument") {
		t.Fatalf("expected unknown-command error, got %q", stderr.String())
	}
}

func TestRunRejectsFormerRepositoryRootFlag(t *testing.T) {
	tool, stdout, stderr := newTestCLI()

	exitCode := tool.Run(context.Background(), []string{"check", "--repo-root", "."})
	if exitCode != usageExitCode {
		t.Fatalf("expected usage exit code, got %d", exitCode)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout for invalid flag, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "--repo-root") {
		t.Fatalf("expected invalid flag error, got %q", stderr.String())
	}
}

func TestRunRejectsLegacySingleDashLongFlag(t *testing.T) {
	tool, stdout, stderr := newTestCLI()

	exitCode := tool.Run(context.Background(), []string{"check", "-format", "json"})
	if exitCode != usageExitCode {
		t.Fatalf("expected usage exit code, got %d", exitCode)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout for invalid flag, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "-format") {
		t.Fatalf("expected invalid flag error, got %q", stderr.String())
	}
}

func TestRunRejectsHelpArgumentsBeyondCommand(t *testing.T) {
	testCases := []struct {
		name      string
		arguments []string
		usage     string
	}{
		{
			name:      "flag after command",
			arguments: []string{"help", "check", "--format", "json"},
			usage:     "--mode",
		},
		{
			name:      "positional after command",
			arguments: []string{"help", "list", "rules"},
			usage:     "<selector>",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			tool, stdout, stderr := newTestCLI()

			exitCode := tool.Run(context.Background(), testCase.arguments)
			if exitCode != usageExitCode {
				t.Fatalf("expected usage exit code, got %d", exitCode)
			}
			if stdout.Len() != 0 {
				t.Fatalf("expected no stdout for invalid help, got %q", stdout.String())
			}
			if !strings.Contains(stderr.String(), "unexpected argument") {
				t.Fatalf("expected invalid-help error, got %q", stderr.String())
			}
			if !strings.Contains(stderr.String(), testCase.usage) {
				t.Fatalf("expected command usage %q, got %q", testCase.usage, stderr.String())
			}
		})
	}
}

func TestRunTreatsRootHelpAsSuccess(t *testing.T) {
	tool, stdout, stderr := newTestCLI()

	exitCode := tool.Run(context.Background(), []string{"help"})
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
	testCases := []struct {
		name      string
		arguments []string
		golden    string
	}{
		{name: "check", arguments: []string{"help", "check"}, golden: "check_help.txt"},
		{name: "list", arguments: []string{"help", "list"}, golden: "list_help.txt"},
		{name: "explain", arguments: []string{"help", "explain"}, golden: "explain_help.txt"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			tool, stdout, stderr := newTestCLI()

			exitCode := tool.Run(context.Background(), testCase.arguments)
			if exitCode != 0 {
				t.Fatalf("expected success exit code for help, got %d", exitCode)
			}
			if stderr.Len() != 0 {
				t.Fatalf("expected no stderr for help, got %q", stderr.String())
			}
			if stdout.String() != readGoldenOutput(t, testCase.golden) {
				t.Fatalf("unexpected command help output:\n%s", stdout.String())
			}
		})
	}
}

func TestRunTreatsFlagHelpAsSuccess(t *testing.T) {
	tool, stdout, stderr := newTestCLI()

	exitCode := tool.Run(context.Background(), []string{"check", "-h"})
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

func TestRunPrintsVersion(t *testing.T) {
	tool, stdout, stderr := newTestCLI()

	exitCode := tool.Run(context.Background(), []string{"version"})
	if exitCode != 0 {
		t.Fatalf("expected success exit code for version, got %d", exitCode)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr for version, got %q", stderr.String())
	}
	if stdout.String() != "test-version\n" {
		t.Fatalf("unexpected version output: %q", stdout.String())
	}
}

/* -------------------------------------------- Usage ------------------------------------------- */

func TestUsageTextListsCommands(t *testing.T) {
	usage := rootUsageText()
	requiredSnippets := []string{
		"quill <command> [flags]",
		"check",
		"fix",
		"doctor",
		"coverage",
		"install",
		"lock",
		"version",
		"init",
		"list",
		"explain",
	}

	for _, snippet := range requiredSnippets {
		if !strings.Contains(usage, snippet) {
			t.Fatalf("usage text missing %q:\n%s", snippet, usage)
		}
	}
}

/* ---------------------------------------- Test Harness ---------------------------------------- */

func newTestCLI() (runner Runner, stdout *bytes.Buffer, stderr *bytes.Buffer) {
	stdout = &bytes.Buffer{}
	stderr = &bytes.Buffer{}
	runner = New(stdout, stderr, "test-version")
	return runner, stdout, stderr
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
