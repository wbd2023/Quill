//go:build windows

package process

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"golang.org/x/sys/windows"
)

/* --------------------------------------- Re-exec dispatch -------------------------------------- */

// TestMain turns the test binary into a cooperating fixture for the Windows job-object regression.
// When QUILL_PROCESS_ROLE is set, the binary does not run the suite: it acts as the managed child
// (spawner) or its descendant (grandchild) instead, so RunCommand can launch a real process tree
// and the test can prove cancellation reaches the descendant. With no role the normal suite runs.
func TestMain(m *testing.M) {
	switch role := os.Getenv(processRoleEnv); role {
	case "spawner":
		runSpawner(os.Getenv(processDirEnv))
		os.Exit(0)
	case "grandchild":
		runGrandchild(os.Getenv(processDirEnv))
		os.Exit(0)
	}

	os.Exit(m.Run())
}

const (
	processRoleEnv = "QUILL_PROCESS_ROLE"
	processDirEnv  = "QUILL_PROCESS_DIR"

	// survivalTickInterval paces the grandchild's survival markers so a survivor is unambiguous
	// within the sampling window while keeping the test fast.
	survivalTickInterval = 25 * time.Millisecond
	// descendantLinger lets any in-flight tick and the job termination settle before the test
	// samples the survival marker, then samples again to detect a survivor.
	descendantLinger = 600 * time.Millisecond
)

// runSpawner is the managed child launched by RunCommand. It starts a genuine descendant (the
// grandchild), waits for the grandchild to confirm it is running, announces readiness, and then
// holds until the Job Object terminates it. The descendant inherits spawner's job membership, so a
// correct tree kill reaches it; a direct-child-only kill (the bug) leaves it alive.
func runSpawner(dir string) {
	executable, err := os.Executable()
	if err != nil {
		fmt.Fprintln(os.Stderr, "spawner: resolve executable:", err)
		os.Exit(1)
	}

	grandchild := exec.Command(executable)
	grandchild.Env = mergeEnv(os.Environ(), map[string]string{
		processRoleEnv: "grandchild",
		processDirEnv:  dir,
	})
	grandchild.Stdout = io.Discard
	grandchild.Stderr = io.Discard
	if err := grandchild.Start(); err != nil {
		fmt.Fprintln(os.Stderr, "spawner: start descendant:", err)
		os.Exit(1)
	}

	if !waitUntilExists(filepath.Join(dir, "started"), 10*time.Second) {
		fmt.Fprintln(os.Stderr, "spawner: descendant never reported started")
		os.Exit(1)
	}

	if err := os.WriteFile(filepath.Join(dir, "ready"), []byte("ready"), 0o600); err != nil {
		fmt.Fprintln(os.Stderr, "spawner: write ready:", err)
		os.Exit(1)
	}

	// Hold until the Job Object kills us. Reaching os.Exit means the tree was not terminated.
	time.Sleep(time.Minute)
	os.Exit(0)
}

// runGrandchild is the descendant. It signals that it is running, then appends a survival marker on
// a steady cadence until the Job Object kills it. A survivor keeps appending; the test samples the
// marker twice and fails if it keeps growing after cancellation.
func runGrandchild(dir string) {
	survival := filepath.Join(dir, "survival")
	if err := os.WriteFile(filepath.Join(dir, "started"), []byte("started"), 0o600); err != nil {
		os.Exit(1)
	}

	for range 4000 {
		appendSurvivalTick(survival)
		time.Sleep(survivalTickInterval)
	}
}

func appendSurvivalTick(path string) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return
	}
	defer file.Close()
	_, _ = file.WriteString("x")
}

// mergeEnv returns parent with the overrides applied (last value wins), folding names
// case-insensitively so the role override actually replaces the inherited entry on Windows.
func mergeEnv(parent []string, overrides map[string]string) []string {
	superseded := make(map[string]bool, len(overrides))
	for name := range overrides {
		superseded[strings.ToUpper(name)] = true
	}

	merged := make([]string, 0, len(parent)+len(overrides))
	for _, entry := range parent {
		name, _, ok := splitEnvEntry(entry)
		if ok && superseded[strings.ToUpper(name)] {
			continue
		}
		merged = append(merged, entry)
	}

	for name, value := range overrides {
		merged = append(merged, name+"="+value)
	}

	return merged
}

func waitUntilExists(path string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(path); err == nil {
			return true
		}
		time.Sleep(readyPollInterval)
	}
	return false
}

/* ----------------------------------------- Guard tests ---------------------------------------- */

// TestCancelJobTreeReportsDoneWithoutProcess is the Windows guard symmetric to the Unix lifecycle
// test: an absent process yields os.ErrProcessDone and records no kill, so Run does not inject a
// context error after a natural exit.
func TestCancelJobTreeReportsDoneWithoutProcess(t *testing.T) {
	job, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		t.Fatalf("create job object: %v", err)
	}
	defer windows.CloseHandle(job)

	state := &killState{}
	if err := cancelJobTree(&exec.Cmd{}, job, state)(); !errors.Is(err, os.ErrProcessDone) {
		t.Fatalf("expected os.ErrProcessDone when process is absent, got %v", err)
	}

	if state.wasKilled() {
		t.Fatal("absent process must not record a kill")
	}
}

/* --------------------------- Windows job-object descendant regression -------------------------- */

// TestWindowsJobObjectKillsDescendantsOnCancel is the native regression for QUILL-TRUST-004. It
// launches a managed child that spawns a real descendant, cancels after the tree reports readiness,
// and proves the descendant never survives: its survival marker must stop growing once the Job
// Object is terminated. Without the job tree kill, cancellation would only reach the direct child
// and the descendant would keep appending, failing this test.
func TestWindowsJobObjectKillsDescendantsOnCancel(t *testing.T) {
	dir := t.TempDir()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	executable, err := os.Executable()
	if err != nil {
		t.Fatalf("locate test binary: %v", err)
	}

	outcome := make(chan commandOutcome, 1)
	go func() {
		result, runErr := RunCommand(ctx, CommandRequest{
			Name:        executable,
			Environment: EnvironmentInherit,
			Variables: map[string]string{
				processRoleEnv: "spawner",
				processDirEnv:  dir,
			},
			Timeout:          time.Minute,
			OutputLimitBytes: 1 << 20,
		})
		outcome <- commandOutcome{result: result, err: runErr}
	}()

	waitForFile(t, filepath.Join(dir, "ready"), 15*time.Second)
	waitForFile(t, filepath.Join(dir, "survival"), 5*time.Second)
	cancel()

	done := <-outcome
	if !done.result.Canceled {
		t.Fatalf("expected canceled result, got %+v (err=%v)", done.result, done.err)
	}

	assertDescendantDidNotSurvive(t, filepath.Join(dir, "survival"))
}

// assertDescendantDidNotSurvive samples the survival marker twice, separated by a settle window. A
// descendant killed by the Job Object stops writing, so the two samples are equal; a survivor keeps
// appending and the second sample is strictly larger.
func assertDescendantDidNotSurvive(t *testing.T, survival string) {
	t.Helper()

	time.Sleep(descendantLinger)
	first := survivalSize(t, survival)

	time.Sleep(descendantLinger)
	second := survivalSize(t, survival)

	if second > first {
		t.Fatalf("descendant survived cancellation: survival marker grew %d -> %d bytes", first, second)
	}
}

func survivalSize(t *testing.T, path string) int {
	t.Helper()

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat survival marker %s: %v", path, err)
	}
	return int(info.Size())
}
