package runtime

import (
	"errors"
	"os"
	"strings"
	"testing"

	"ciphera/tools/internal/testutil"
)

/* -------------------------------------- Command Execution ------------------------------------- */

func TestRunCommandResolvesCommandsFromProvidedPath(t *testing.T) {
	tempDir := t.TempDir()
	testutil.WriteExecutable(
		t,
		tempDir,
		"test-tool",
		"#!/bin/sh\necho resolved\n",
	)

	result, err := RunCommand(CommandRequest{
		Directory:   tempDir,
		Environment: map[string]string{"PATH": tempDir},
		Name:        "test-tool",
	})
	if err != nil {
		t.Fatalf("RunCommand: %v", err)
	}

	if result.Output != "resolved\n" {
		t.Fatalf("unexpected output %q", result.Output)
	}
}

func TestRunCommandTimesOut(t *testing.T) {
	tempDir := t.TempDir()
	testutil.WriteExecutable(
		t,
		tempDir,
		"slow-tool",
		"#!/bin/sh\nsleep 5\n",
	)

	result, err := RunCommand(CommandRequest{
		Directory:        tempDir,
		Environment:      map[string]string{"PATH": commandSearchPath(tempDir)},
		Name:             "slow-tool",
		TimeoutSeconds:   1,
		OutputLimitBytes: 1024,
	})
	if err == nil {
		t.Fatal("expected timeout error")
	}

	var commandErr CommandError
	if !errors.As(err, &commandErr) {
		t.Fatalf("expected CommandError, got %T", err)
	}

	if !result.TimedOut || !commandErr.Result.TimedOut {
		t.Fatalf("expected timeout result, got %+v", result)
	}
}

func commandSearchPath(tempDir string) (value string) {
	return tempDir + string(os.PathListSeparator) + os.Getenv("PATH")
}

func TestRunCommandCapsOutput(t *testing.T) {
	tempDir := t.TempDir()
	testutil.WriteExecutable(
		t,
		tempDir,
		"loud-tool",
		"#!/bin/sh\nprintf 1234567890\n",
	)

	result, err := RunCommand(CommandRequest{
		Directory:        tempDir,
		Environment:      map[string]string{"PATH": tempDir},
		Name:             "loud-tool",
		OutputLimitBytes: 4,
	})
	if err != nil {
		t.Fatalf("RunCommand: %v", err)
	}

	if result.Output != "1234" || !result.Truncated {
		t.Fatalf("expected truncated output, got %+v", result)
	}
}

func TestCommandErrorIncludesExitCodeAndOutput(t *testing.T) {
	tempDir := t.TempDir()
	testutil.WriteExecutable(
		t,
		tempDir,
		"bad-tool",
		"#!/bin/sh\necho failure\nexit 7\n",
	)

	result, err := RunCommand(CommandRequest{
		Directory:        tempDir,
		Environment:      map[string]string{"PATH": tempDir},
		Name:             "bad-tool",
		OutputLimitBytes: 1024,
	})
	if err == nil {
		t.Fatal("expected command error")
	}

	if result.ExitCode != 7 || !strings.Contains(result.Output, "failure") {
		t.Fatalf("unexpected command result: %+v", result)
	}
}
