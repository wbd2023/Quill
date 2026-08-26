package process

import (
	"context"
	"errors"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/wbd2023/quill/internal/testutil"
)

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
