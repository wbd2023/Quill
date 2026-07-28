//go:build !windows

package process

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/wbd2023/quill/internal/testutil"
)

// TestCancelProcessGroupMarksKilledOnSuccessfulKill is the deterministic regression for
// kill-based attribution: when the cancellation actually kills a running process, the killState
// must record it so TimedOut/Canceled are attributed from the kill rather than from a Unix exit
// code.
func TestCancelProcessGroupMarksKilledOnSuccessfulKill(t *testing.T) {
	dir := t.TempDir()
	testutil.WriteExecutable(t, dir, "sleeper", "#!/bin/sh\nsleep 30\n")

	command := exec.Command(filepath.Join(dir, "sleeper"))
	command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := command.Start(); err != nil {
		t.Fatalf("start helper: %v", err)
	}

	state := &killState{}
	if err := cancelProcessGroup(command, state)(); err != nil {
		t.Fatalf("cancel: %v", err)
	}

	_ = command.Wait()

	if !state.wasKilled() {
		t.Fatal("expected kill to be recorded on state")
	}
}

// TestCancelProcessGroupDoesNotMarkKilledAfterNaturalExit confirms a child that exits on its own
// is not recorded as killed: the group is gone (ESRCH), os.ErrProcessDone is returned, and Run
// will not inject a context error or attribute termination.
func TestCancelProcessGroupDoesNotMarkKilledAfterNaturalExit(t *testing.T) {
	command := exec.Command("true")
	command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := command.Run(); err != nil {
		t.Fatalf("run helper: %v", err)
	}

	state := &killState{}
	err := cancelProcessGroup(command, state)()
	if !errors.Is(err, os.ErrProcessDone) {
		t.Fatalf("expected os.ErrProcessDone after natural exit, got %v", err)
	}

	if state.wasKilled() {
		t.Fatal("natural exit must not record a kill")
	}
}

// TestCancelProcessGroupReportsDoneWithoutProcess confirms the nil-process guard maps to
// os.ErrProcessDone without recording a kill.
func TestCancelProcessGroupReportsDoneWithoutProcess(t *testing.T) {
	state := &killState{}
	err := cancelProcessGroup(&exec.Cmd{}, state)()
	if !errors.Is(err, os.ErrProcessDone) {
		t.Fatalf("expected os.ErrProcessDone when process is absent, got %v", err)
	}

	if state.wasKilled() {
		t.Fatal("absent process must not record a kill")
	}
}

// TestSignalProcessGroupReportsDoneAfterNaturalExit guards the ESRCH mapping directly.
func TestSignalProcessGroupReportsDoneAfterNaturalExit(t *testing.T) {
	command := exec.Command("true")
	command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := command.Run(); err != nil {
		t.Fatalf("run helper: %v", err)
	}

	if err := signalProcessGroup(command.Process.Pid); !errors.Is(err, os.ErrProcessDone) {
		t.Fatalf("expected os.ErrProcessDone after natural exit, got %v", err)
	}
}
