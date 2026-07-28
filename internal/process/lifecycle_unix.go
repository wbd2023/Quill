//go:build !windows

package process

import (
	"errors"
	"os"
	"os/exec"
	"syscall"
)

// configureProcessTree isolates the child in its own process group so the whole tree can be
// terminated when the request context is cancelled or the timeout fires. On POSIX systems the
// child becomes its own process-group leader (its process id is the group id), so signalling the
// negative process id delivers the signal to every member of the group, including grandchildren
// the child may have spawned. The returned killState records whether the cancellation actually
// killed a running process, so termination cause is attributed from the kill rather than from a
// Unix-specific exit code.
func configureProcessTree(command *exec.Cmd) (state *killState) {
	state = &killState{}
	command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	command.Cancel = cancelProcessGroup(command, state)
	command.WaitDelay = childWaitDelay
	return state
}

// cancelProcessGroup returns the exec.Cancel function for the child's process group. A successful
// kill is recorded on state; a group that has already exited yields os.ErrProcessDone so Run does
// not inject a context error after a natural exit and does not record a kill.
func cancelProcessGroup(command *exec.Cmd, state *killState) (cancel func() (err error)) {
	return func() (err error) {
		if command.Process == nil {
			return os.ErrProcessDone
		}

		if err = signalProcessGroup(command.Process.Pid); err != nil {
			return err
		}

		state.markKilled()
		return nil
	}
}

// signalProcessGroup sends SIGKILL to the child's process group and maps the result so that a
// group that has already exited is reported as os.ErrProcessDone. This prevents Run from injecting
// a context error after the child exits naturally.
func signalProcessGroup(pid int) (err error) {
	if err = syscall.Kill(-pid, syscall.SIGKILL); err != nil {
		if errors.Is(err, syscall.ESRCH) {
			return os.ErrProcessDone
		}

		return err
	}

	return nil
}
