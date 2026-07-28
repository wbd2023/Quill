package process

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
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

func TestExecutableExtensionsParsesPathext(t *testing.T) {
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

/* --------------------------------- Environment Key Regression --------------------------------- */

// TestEnvKeyFoldsCaseOnlyOnWindows documents and pins the canonicalisation rule: Windows env names
// fold to upper case so PATH and Path merge; POSIX names are case-sensitive and stay verbatim.
func TestEnvKeyFoldsCaseOnlyOnWindows(t *testing.T) {
	if runtime.GOOS == "windows" {
		if envKey("Path") != "PATH" || envKey("path") != "PATH" {
			t.Fatalf("Windows envKey must fold names to upper case: %q %q",
				envKey("Path"), envKey("path"))
		}
		return
	}

	if envKey("Path") != "Path" {
		t.Fatalf("POSIX envKey must preserve case, got %q", envKey("Path"))
	}
}

// TestBuildEnvironmentMergesCaseInsensitiveKeysOnWindows guards the Windows override bug: a
// Variables PATH must replace an inherited Path rather than sit beside it. Skipped on POSIX, where
// the two names are genuinely distinct variables.
func TestBuildEnvironmentMergesCaseInsensitiveKeysOnWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("case-insensitive environment keys are a Windows behaviour")
	}

	t.Setenv("Path", "parent-value")

	environment := buildEnvironment(EnvironmentInherit, map[string]string{
		"PATH": "child-value",
	})

	var pathEntries int
	for _, entry := range environment {
		name, _, _ := strings.Cut(entry, "=")
		if strings.EqualFold(name, "PATH") {
			pathEntries++
		}

		if entry == "Path=parent-value" {
			t.Fatalf("inherited Path must be overridden by Variables PATH, got %v", environment)
		}
	}

	if pathEntries != 1 {
		t.Fatalf("PATH must appear exactly once after case-insensitive merge, got %d in %v",
			pathEntries, environment)
	}
}

/* -------------------------------- Termination Cause Regression -------------------------------- */

// TestRunCommandNaturalSuccessNotFlaggedAsTerminated pins the contract that a process which exits
// on its own is never marked TimedOut or Canceled, even though a request timeout is configured: the
// flags derive from Run's error, not from the context state observed after the fact.
func TestRunCommandNaturalSuccessNotFlaggedAsTerminated(t *testing.T) {
	tempDir := t.TempDir()
	testutil.WriteExecutable(t, tempDir, "quick-tool", "#!/bin/sh\necho done\n")

	result, err := RunCommand(context.Background(), CommandRequest{
		Name:        "quick-tool",
		Environment: EnvironmentInherit,
		Variables:   map[string]string{"PATH": tempDir},
		Directory:   tempDir,
		Timeout:     5 * time.Second,
	})
	if err != nil {
		t.Fatalf("RunCommand: %v", err)
	}

	if result.TimedOut || result.Canceled {
		t.Fatalf("natural success must not be flagged as terminated: %+v", result)
	}
}

// TestRunCommandParentDeadlineNotReportedAsRequestTimeout is the deterministic regression for a
// caller-supplied deadline that fires before Quill's request timer: the child is terminated by the
// parent deadline, so TimedOut must stay false (it is not Quill's timer) while the error still
// surfaces context.DeadlineExceeded.
func TestRunCommandParentDeadlineNotReportedAsRequestTimeout(t *testing.T) {
	tempDir := t.TempDir()
	testutil.WriteExecutable(t, tempDir, "slow-tool", "#!/bin/sh\nsleep 5\n")

	parent, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	result, err := RunCommand(parent, CommandRequest{
		Name:             "slow-tool",
		Environment:      EnvironmentInherit,
		Variables:        map[string]string{"PATH": commandSearchPath(tempDir)},
		Directory:        tempDir,
		Timeout:          5 * time.Second,
		OutputLimitBytes: 1024,
	})
	if err == nil {
		t.Fatal("expected error from parent deadline")
	}

	if result.TimedOut {
		t.Fatalf("parent deadline must not be reported as request timeout: %+v", result)
	}

	if result.Canceled {
		t.Fatalf("parent deadline must not be reported as cancellation: %+v", result)
	}

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("parent deadline must surface as context.DeadlineExceeded, got %v", err)
	}
}

/* ------------------------------ Executable Resolution Regression ------------------------------ */

// TestFindExecutableDoesNotPreferExtensionlessCandidate guards the Windows resolution rule: when
// PATHEXT extensions apply, a bare command name resolves through the qualified executable rather
// than an extensionless file, matching os/exec.
func TestFindExecutableDoesNotPreferExtensionlessCandidate(t *testing.T) {
	dir := t.TempDir()
	testutil.WriteExecutable(t, dir, "tool", "#!/bin/sh\n")
	testutil.WriteExecutable(t, dir, "tool.exe", "#!/bin/sh\n")

	resolved := findExecutable(dir, "tool", []string{".exe"})
	if want := filepath.Join(dir, "tool.exe"); resolved != want {
		t.Fatalf("findExecutable = %q, want %q (must not prefer extensionless)", resolved, want)
	}
}

// TestFindExecutableUsesBareCandidateWhenAlreadyQualified ensures a command carrying an extension
// resolves to itself rather than appending a further PATHEXT suffix.
func TestFindExecutableUsesBareCandidateWhenAlreadyQualified(t *testing.T) {
	dir := t.TempDir()
	testutil.WriteExecutable(t, dir, "tool.exe", "#!/bin/sh\n")

	resolved := findExecutable(dir, "tool.exe", []string{".exe"})
	if want := filepath.Join(dir, "tool.exe"); resolved != want {
		t.Fatalf("findExecutable = %q, want %q", resolved, want)
	}
}

// TestResolvePathExtAppliesWindowsDefault pins the PATHEXT fallback: an explicit value is
// honoured, an empty value uses the Windows default (so bare-name resolution does not degrade)
// and stays empty on POSIX.
func TestResolvePathExtAppliesWindowsDefault(t *testing.T) {
	if got := resolvePathExt(".FOO;.BAR"); got != ".FOO;.BAR" {
		t.Fatalf("explicit PATHEXT must be honoured, got %q", got)
	}

	if runtime.GOOS == "windows" {
		if got := resolvePathExt(""); got != defaultWindowsPathExt {
			t.Fatalf("empty PATHEXT must fall back to Windows default, got %q", got)
		}
		return
	}

	if got := resolvePathExt(""); got != "" {
		t.Fatalf("POSIX empty PATHEXT must stay empty, got %q", got)
	}
}

// TestExecutableExtensionsNormalisesMissingDot ensures malformed PATHEXT entries without a leading
// dot still resolve correctly: EXE becomes .EXE so go.exe is found rather than goEXE.
func TestExecutableExtensionsNormalisesMissingDot(t *testing.T) {
	extensions := executableExtensions("EXE;.CMD;bat")
	want := []string{".EXE", ".CMD", ".bat"}
	if len(extensions) != len(want) {
		t.Fatalf("extensions = %v, want %v", extensions, want)
	}
	for index, expected := range want {
		if extensions[index] != expected {
			t.Fatalf("extensions[%d] = %q, want %q", index, extensions[index], expected)
		}
	}
}

// TestResolveDirectPathRejectsNonexistent confirms a direct path that is not executable is reported
// as not found rather than returned unchecked.
func TestResolveDirectPathRejectsNonexistent(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing")
	if _, err := resolveDirectPath(missing, []string{".exe"}); err == nil {
		t.Fatal("expected error for non-existent direct path")
	}
}

// TestResolveDirectPathCompletesExtensionlessViaPathExt ensures an extensionless direct path is
// completed through PATHEXT and verified executable, so os/exec never has to requalify it at start.
func TestResolveDirectPathCompletesExtensionlessViaPathExt(t *testing.T) {
	dir := t.TempDir()
	testutil.WriteExecutable(t, dir, "tool.exe", "#!/bin/sh\n")

	resolved, err := resolveDirectPath(filepath.Join(dir, "tool"), []string{".exe"})
	if err != nil {
		t.Fatalf("resolveDirectPath: %v", err)
	}

	if want := filepath.Join(dir, "tool.exe"); resolved != want {
		t.Fatalf("resolved = %q, want %q", resolved, want)
	}
}

// TestResolveDirectPathReturnsVerifiedExecutable ensures a direct path that already exists and is
// executable is returned verbatim.
func TestResolveDirectPathReturnsVerifiedExecutable(t *testing.T) {
	dir := t.TempDir()
	testutil.WriteExecutable(t, dir, "tool", "#!/bin/sh\n")

	resolved, err := resolveDirectPath(filepath.Join(dir, "tool"), nil)
	if err != nil {
		t.Fatalf("resolveDirectPath: %v", err)
	}

	if want := filepath.Join(dir, "tool"); resolved != want {
		t.Fatalf("resolved = %q, want %q", resolved, want)
	}
}

// TestSplitEnvEntryPreservesWindowsPerDriveEntries is the deterministic seam for the per-drive
// collapse bug: entries beginning with '=' must split at the second '=' so distinct drives
// ("=C:" and "=D:") stay distinct instead of collapsing to an empty key, and values containing '='
// are preserved verbatim.
func TestSplitEnvEntryPreservesWindowsPerDriveEntries(t *testing.T) {
	cases := []struct {
		entry string
		key   string
		value string
	}{
		{"=C:=C:\\work", "=C:", "C:\\work"},
		{"=D:=D:\\workdir", "=D:", "D:\\workdir"},
		{"PATH=/usr/bin", "PATH", "/usr/bin"},
		{"EMPTY=", "EMPTY", ""},
		{"TOKEN=a=b=c", "TOKEN", "a=b=c"},
	}

	for _, example := range cases {
		key, value, ok := splitEnvEntry(example.entry)
		if !ok || key != example.key || value != example.value {
			t.Fatalf("splitEnvEntry(%q) = (%q, %q, %v), want (%q, %q, true)",
				example.entry, key, value, ok, example.key, example.value)
		}
	}

	for _, malformed := range []string{"", "NO_EQUALS"} {
		if _, _, ok := splitEnvEntry(malformed); ok {
			t.Fatalf("splitEnvEntry(%q) must reject malformed entries", malformed)
		}
	}
}

/* -------------------------------- Termination Attribution Seam -------------------------------- */

// TestKillStateRecordsKill is the cross-platform deterministic seam for kill-based attribution:
// termination cause derives from whether the cancellation killed a running process, not from a
// platform-specific exit code.
func TestKillStateRecordsKill(t *testing.T) {
	state := &killState{}
	if state.wasKilled() {
		t.Fatal("fresh state must not report killed")
	}

	state.markKilled()
	if !state.wasKilled() {
		t.Fatal("markKilled must record the kill")
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
