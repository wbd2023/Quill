package process

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/wbd2023/quill/internal/testutil"
)

// readyPollInterval paces the ready-marker poll so the loop does not busy-spin.
const readyPollInterval = 2 * time.Millisecond

/* -------------------------------------- Command Execution ------------------------------------- */

func TestRunCommandResolvesCommandsFromProvidedPath(t *testing.T) {
	tempDir := t.TempDir()
	testutil.WriteExecutable(
		t,
		tempDir,
		"test-tool",
		"#!/bin/sh\necho resolved\n",
	)

	result, err := RunCommand(context.Background(), CommandRequest{
		Name:        "test-tool",
		Environment: EnvironmentInherit,
		Variables:   map[string]string{"PATH": tempDir},
		Directory:   tempDir,
	})
	if err != nil {
		t.Fatalf("RunCommand: %v", err)
	}

	if result.Output != "resolved\n" {
		t.Fatalf("unexpected output %q", result.Output)
	}

	if result.Stdout != "resolved\n" || result.Stderr != "" {
		t.Fatalf("unexpected separated streams stdout=%q stderr=%q",
			result.Stdout, result.Stderr)
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

	result, err := RunCommand(context.Background(), CommandRequest{
		Name:             "slow-tool",
		Environment:      EnvironmentInherit,
		Variables:        map[string]string{"PATH": commandSearchPath(tempDir)},
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

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected error to wrap context.DeadlineExceeded, got %v", err)
	}

	if !result.TimedOut {
		t.Fatalf("expected timeout result, got %+v", result)
	}

	if result.Canceled {
		t.Fatalf("timeout must not be reported as cancellation, got %+v", result)
	}

	if !strings.Contains(commandErr.Error(), "timed out") {
		t.Fatalf("expected timeout error message, got %q", commandErr.Error())
	}
}

func TestRunCommandReportsParentCancellation(t *testing.T) {
	tempDir := t.TempDir()
	marker := filepath.Join(tempDir, "ready")
	testutil.WriteExecutable(
		t,
		tempDir,
		"cancellable-tool",
		"#!/bin/sh\n> ready\nsleep 30\n",
	)

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	resultChan := make(chan commandOutcome, 1)
	go func() {
		result, err := RunCommand(ctx, CommandRequest{
			Name:             "cancellable-tool",
			Environment:      EnvironmentInherit,
			Variables:        map[string]string{"PATH": commandSearchPath(tempDir)},
			Directory:        tempDir,
			Timeout:          5 * time.Second,
			OutputLimitBytes: 1024,
		})
		resultChan <- commandOutcome{result: result, err: err}
	}()

	waitForFile(t, marker, 5*time.Second)
	cancel()

	outcome := <-resultChan
	if outcome.err == nil {
		t.Fatal("expected cancellation error")
	}

	if !outcome.result.Canceled {
		t.Fatalf("expected canceled result, got %+v", outcome.result)
	}

	if outcome.result.TimedOut {
		t.Fatalf("cancellation must not be reported as timeout, got %+v", outcome.result)
	}

	if !errors.Is(outcome.err, context.Canceled) {
		t.Fatalf("expected error to wrap context.Canceled, got %v", outcome.err)
	}
}

func TestRunCommandCapsOutput(t *testing.T) {
	tempDir := t.TempDir()
	testutil.WriteExecutable(
		t,
		tempDir,
		"loud-tool",
		"#!/bin/sh\nprintf 1234567890\n",
	)

	result, err := RunCommand(context.Background(), CommandRequest{
		Name:             "loud-tool",
		Environment:      EnvironmentInherit,
		Variables:        map[string]string{"PATH": tempDir},
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

	result, err := RunCommand(context.Background(), CommandRequest{
		Name:             "bad-tool",
		Environment:      EnvironmentInherit,
		Variables:        map[string]string{"PATH": tempDir},
		Directory:        tempDir,
		OutputLimitBytes: 1024,
	})
	if err == nil {
		t.Fatal("expected command error")
	}

	var commandErr CommandError
	if !errors.As(err, &commandErr) {
		t.Fatalf("expected CommandError, got %T", err)
	}

	// A non-zero exit is an ordinary failure, not cancellation or timeout.
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("non-zero exit must not wrap a context sentinel, got %v", err)
	}

	if result.ExitCode != 7 || !strings.Contains(result.Output, "failure") {
		t.Fatalf("unexpected command result: %+v", result)
	}
}

/* --------------------------------------- Stream Capture --------------------------------------- */

func TestRunCommandCapturesStdoutAndStderrSeparately(t *testing.T) {
	tempDir := t.TempDir()
	testutil.WriteExecutable(
		t,
		tempDir,
		"split-tool",
		"#!/bin/sh\necho out\necho err 1>&2\n",
	)

	result, err := RunCommand(context.Background(), CommandRequest{
		Name:        "split-tool",
		Environment: EnvironmentInherit,
		Variables:   map[string]string{"PATH": tempDir},
		Directory:   tempDir,
	})
	if err != nil {
		t.Fatalf("RunCommand: %v", err)
	}

	if result.Stdout != "out\n" {
		t.Fatalf("Stdout = %q, want %q", result.Stdout, "out\n")
	}

	if result.Stderr != "err\n" {
		t.Fatalf("Stderr = %q, want %q", result.Stderr, "err\n")
	}

	if !strings.Contains(result.Output, "out") || !strings.Contains(result.Output, "err") {
		t.Fatalf("combined Output must preserve both streams, got %q", result.Output)
	}
}

func TestRunCommandTruncatesStreamsSeparately(t *testing.T) {
	tempDir := t.TempDir()
	testutil.WriteExecutable(
		t,
		tempDir,
		"both-streams",
		"#!/bin/sh\nprintf AAAAAAAAAA\nprintf BBBBBBBBBB 1>&2\n",
	)

	result, err := RunCommand(context.Background(), CommandRequest{
		Name:             "both-streams",
		Environment:      EnvironmentInherit,
		Variables:        map[string]string{"PATH": tempDir},
		Directory:        tempDir,
		OutputLimitBytes: 4,
	})
	if err != nil {
		t.Fatalf("RunCommand: %v", err)
	}

	if result.Stdout != "AAAA" || result.Stderr != "BBBB" {
		t.Fatalf("separate streams must each be bounded: stdout=%q stderr=%q",
			result.Stdout, result.Stderr)
	}

	if !result.Truncated {
		t.Fatalf("expected truncated result when both streams overflow, got %+v", result)
	}
}

func TestRunCommandPipesByteStdin(t *testing.T) {
	tempDir := t.TempDir()
	testutil.WriteExecutable(
		t,
		tempDir,
		"echo-tool",
		"#!/bin/sh\ncat\n",
	)

	result, err := RunCommand(context.Background(), CommandRequest{
		Name:        "echo-tool",
		Environment: EnvironmentInherit,
		Variables:   map[string]string{"PATH": commandSearchPath(tempDir)},
		Directory:   tempDir,
		Stdin:       []byte("hello stdin\n"),
	})
	if err != nil {
		t.Fatalf("RunCommand: %v", err)
	}

	if result.Stdout != "hello stdin\n" {
		t.Fatalf("Stdout = %q, want %q", result.Stdout, "hello stdin\n")
	}
}

/* ----------------------------------- Environment Resolution ----------------------------------- */

func TestRunCommandInheritedEnvironmentOverridesParent(t *testing.T) {
	t.Setenv("QUILL_INHERIT", "parent")

	tempDir := t.TempDir()
	testutil.WriteExecutable(
		t,
		tempDir,
		"print-var",
		"#!/bin/sh\necho \"${QUILL_INHERIT}\"\n",
	)

	result, err := RunCommand(context.Background(), CommandRequest{
		Name:        "print-var",
		Environment: EnvironmentInherit,
		Variables: map[string]string{
			"PATH":          tempDir,
			"QUILL_INHERIT": "child",
		},
		Directory: tempDir,
	})
	if err != nil {
		t.Fatalf("RunCommand: %v", err)
	}

	if result.Stdout != "child\n" {
		t.Fatalf("inherited variable must be overridden by Variables: %q", result.Stdout)
	}
}

func TestRunCommandIsolatedEnvironmentExcludesParent(t *testing.T) {
	t.Setenv("QUILL_SECRET", "leaked")

	tempDir := t.TempDir()
	testutil.WriteExecutable(
		t,
		tempDir,
		"probe-secret",
		"#!/bin/sh\nif [ -n \"${QUILL_SECRET}\" ]; then echo leaked; else echo clean; fi\n",
	)

	result, err := RunCommand(context.Background(), CommandRequest{
		Name:        "probe-secret",
		Environment: EnvironmentIsolated,
		Variables:   map[string]string{"PATH": tempDir},
		Directory:   tempDir,
	})
	if err != nil {
		t.Fatalf("RunCommand: %v", err)
	}

	if result.Stdout != "clean\n" {
		t.Fatalf("isolated environment must not inherit parent secrets: %q", result.Stdout)
	}
}

func TestBuildEnvironmentIsolatedEmitsOnlyVariables(t *testing.T) {
	t.Setenv("QUILL_PARENT_LEAK", "present")

	environment := buildEnvironment(EnvironmentIsolated, map[string]string{
		"B": "2",
		"A": "1",
	})

	if len(environment) != 2 || environment[0] != "A=1" || environment[1] != "B=2" {
		t.Fatalf("isolated environment must emit only sorted Variables, got %v", environment)
	}

	for _, entry := range environment {
		if strings.HasPrefix(entry, "QUILL_PARENT_LEAK=") {
			t.Fatalf("isolated environment leaked parent entry %q", entry)
		}
	}
}

func TestBuildEnvironmentInheritDeduplicatesOverrides(t *testing.T) {
	t.Setenv("QUILL_BUILD_TEST", "parent")

	environment := buildEnvironment(EnvironmentInherit, map[string]string{
		"QUILL_BUILD_TEST": "child",
	})

	var seen int
	for _, entry := range environment {
		if entry == "QUILL_BUILD_TEST=child" {
			seen++
		}

		if entry == "QUILL_BUILD_TEST=parent" {
			t.Fatalf("inherited parent entry must be overridden, not duplicated: %v", environment)
		}
	}

	if seen != 1 {
		t.Fatalf("override entry must appear exactly once, got %d in %v", seen, environment)
	}
}

func TestExecutableExtensionsParsesPathExt(t *testing.T) {
	extensions := executableExtensions(".COM;.EXE;.BAT")
	if len(extensions) != 3 || extensions[0] != ".COM" || extensions[2] != ".BAT" {
		t.Fatalf("unexpected extensions %v", extensions)
	}

	if got := executableExtensions(""); got != nil {
		t.Fatalf("empty PATHEXT must yield no extensions, got %v", got)
	}
}

func TestResolveExecutableRejectsEmptyName(t *testing.T) {
	if _, err := ResolveExecutable(map[string]string{"PATH": "/bin"}, ""); err == nil {
		t.Fatal("expected error resolving empty command name")
	}
}

/* ----------------------------------- Concurrency Regression ----------------------------------- */

// TestRunCommandConcurrentStreamsAreRaceFree drives heavy concurrent stdout and stderr traffic so
// the shared combined buffer is reached from two goroutines at once. Under -race this fails without
// the boundedBuffer mutex and passes with it.
func TestRunCommandConcurrentStreamsAreRaceFree(t *testing.T) {
	tempDir := t.TempDir()
	testutil.WriteExecutable(
		t,
		tempDir,
		"chatty-tool",
		"#!/bin/sh\ni=0\nwhile [ $i -lt 300 ]; do\n"+
			"printf STDOUT%d $i\necho ERR$i 1>&2\ni=$((i+1))\ndone\n",
	)

	result, err := RunCommand(context.Background(), CommandRequest{
		Name:             "chatty-tool",
		Environment:      EnvironmentInherit,
		Variables:        map[string]string{"PATH": commandSearchPath(tempDir)},
		Directory:        tempDir,
		OutputLimitBytes: 1 << 16,
	})
	if err != nil {
		t.Fatalf("RunCommand: %v", err)
	}

	if !strings.Contains(result.Stdout, "STDOUT") {
		t.Fatalf("stdout stream lost under concurrent writes: %q", result.Stdout)
	}

	if !strings.Contains(result.Stderr, "ERR") {
		t.Fatalf("stderr stream lost under concurrent writes: %q", result.Stderr)
	}

	if !strings.Contains(result.Output, "STDOUT") || !strings.Contains(result.Output, "ERR") {
		t.Fatalf("combined output must contain both streams: %q", result.Output)
	}
}

/* ------------------------------------------- Helpers ------------------------------------------ */

type commandOutcome struct {
	result CommandResult
	err    error
}

func commandSearchPath(tempDir string) (path string) {
	return tempDir + string(os.PathListSeparator) + os.Getenv("PATH")
}

// waitForFile polls for path until it exists or the timeout elapses. It waits for a real signal
// (the child writing its ready marker) rather than an arbitrary sleep, so cancellation is exercised
// against a genuinely running child.
func waitForFile(t *testing.T, path string, timeout time.Duration) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(path); err == nil {
			return
		}

		<-time.After(readyPollInterval)
	}

	t.Fatalf("ready marker %s never appeared within %s", path, timeout)
}
