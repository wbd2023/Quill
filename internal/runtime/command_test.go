package runtime

import (
	"errors"
	"os"
	"strings"
	"testing"
	"time"

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
		Name:        "test-tool",
		Environment: map[string]string{"PATH": tempDir},
		Directory:   tempDir,
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
		Name:             "slow-tool",
		Environment:      map[string]string{"PATH": commandSearchPath(tempDir)},
		Directory:        tempDir,
		Timeout:          time.Second,
		OutputLimitBytes: 1024,
	})
	if err == nil {
		t.Fatal("expected timeout error")
	}

	var commandErr CommandError
	if !errors.As(err, &commandErr) {
		t.Fatalf("expected CommandError, got %T", err)
	}

	if !result.TimedOut {
		t.Fatalf("expected timeout result, got %+v", result)
	}

	if !strings.Contains(commandErr.Error(), "timed out") {
		t.Fatalf("expected timeout error message, got %q", commandErr.Error())
	}
}

func commandSearchPath(tempDir string) (path string) {
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
		Name:             "loud-tool",
		Environment:      map[string]string{"PATH": tempDir},
		Directory:        tempDir,
		OutputLimitBytes: 4,
	})
	if err != nil {
		t.Fatalf("RunCommand: %v", err)
	}

	if result.Output != "1234" || !result.Truncated {
		t.Fatalf("expected truncated output, got %+v", result)
	}
}

func TestRunCommandReturnsExitCodeAndOutput(t *testing.T) {
	tempDir := t.TempDir()
	testutil.WriteExecutable(
		t,
		tempDir,
		"bad-tool",
		"#!/bin/sh\necho failure\nexit 7\n",
	)

	result, err := RunCommand(CommandRequest{
		Name:             "bad-tool",
		Environment:      map[string]string{"PATH": tempDir},
		Directory:        tempDir,
		OutputLimitBytes: 1024,
	})
	if err == nil {
		t.Fatal("expected command error")
	}

	if result.ExitCode != 7 || !strings.Contains(result.Output, "failure") {
		t.Fatalf("unexpected command result: %+v", result)
	}
}
