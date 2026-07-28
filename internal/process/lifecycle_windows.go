//go:build windows

package process

import (
	"errors"
	"os"
	"os/exec"
)

// configureProcessTree applies the best process-tree termination available without a
// platform-specific dependency. exec.CommandContext kills the direct child when the request
// context is cancelled or the timeout fires; WaitDelay bounds the wait for the child to exit
// after the kill signal. The returned killState records whether the cancellation actually killed a
// running process, so termination cause is attributed from the kill rather than from the exit code
// (a Windows kill exits with code 1). Grandchildren are not tracked in this build; a full Windows
// job-object tree kill would require a platform dependency that is not justified for the current
// Pack set.
func configureProcessTree(command *exec.Cmd) (state *killState) {
	state = &killState{}
	command.Cancel = cancelProcess(command, state)
	command.WaitDelay = childWaitDelay
	return state
}

// cancelProcess returns the exec.Cancel function for the child. A successful kill is recorded on
// state; a process that has already exited yields os.ErrProcessDone so Run does not inject a
// context error after a natural exit and does not record a kill.
func cancelProcess(command *exec.Cmd, state *killState) (cancel func() (err error)) {
	return func() (err error) {
		if command.Process == nil || command.ProcessState != nil {
			return os.ErrProcessDone
		}

		if err = command.Process.Kill(); err != nil {
			if errors.Is(err, os.ErrProcessDone) {
				return os.ErrProcessDone
			}

			return err
		}

		state.markKilled()
		return nil
	}
}
