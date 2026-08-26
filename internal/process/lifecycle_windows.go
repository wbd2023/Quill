//go:build windows

package process

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

/* ------------------------------------------ Constants ----------------------------------------- */
// childProcessAccess is the access mask used to re-open the suspended child by PID so it can be
// bound to the Job Object. AssignProcessToJobObject requires PROCESS_SET_QUOTA; PROCESS_TERMINATE
// lets the caller fail closed by killing a child that could not be bound.
const childProcessAccess = windows.PROCESS_SET_QUOTA | windows.PROCESS_TERMINATE

// jobExitCode is the exit status reported for every process the managed job terminates on
// cancellation or timeout. It matches the code os.Process.Kill reports on Windows (1), keeping
// exit-code consumers consistent across the kill paths.
const jobExitCode = 1

/* -------------------------------------- Job Object Setup -------------------------------------- */
// configureProcessTree isolates the child and every descendant it spawns inside a Windows Job
// Object, so the whole tree is terminated when the request context is cancelled or the timeout
// fires. The job carries JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE, so if Quill itself dies the operating
// system reaps the child tree, and closing the handle after a natural exit clears any descendant
// that lingers past the child.
//
// Race safety: the child is created suspended (CREATE_SUSPENDED), so exec.Cmd.Start returns before
// the child executes a single instruction. The returned afterStart hook binds the suspended child
// to the job and only then resumes its main thread, so any descendant the child later spawns is
// born inside the job and cannot escape cancellation. RunCommand must call afterStart between Start
// and Wait.
//
// The returned error is non-nil only when the job cannot be established; in that case RunCommand
// fails closed and never starts the child. A failure from afterStart (job binding or resume) means
// the child never ran: RunCommand kills the suspended child and fails closed.
func configureProcessTree(
	command *exec.Cmd,
) (
	state *killState,
	afterStart func(cmd *exec.Cmd) error,
	release func(),
	err error,
) {
	job, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("create job object: %w", err)
	}

	info := windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{}
	info.BasicLimitInformation.LimitFlags = windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE
	if _, err = windows.SetInformationJobObject(
		job,
		windows.JobObjectExtendedLimitInformation,
		uintptr(unsafe.Pointer(&info)),
		uint32(unsafe.Sizeof(info)),
	); err != nil {
		// The job was never handed to the child, so releasing the unusable handle is all that
		// remains; nothing else references it.
		_ = windows.CloseHandle(job)
		return nil, nil, nil, fmt.Errorf("configure job object: %w", err)
	}

	state = &killState{}
	command.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: windows.CREATE_SUSPENDED,
	}
	command.Cancel = cancelJobTree(command, job, state)
	command.WaitDelay = childWaitDelay

	release = func() {
		// Closing the kill-on-close job releases any lingering descendant; if the close itself
		// fails, Quill exiting closes the last handle anyway and the tree is reaped then.
		_ = windows.CloseHandle(job)
	}
	return state, bindSuspendedChild(job, state), release, nil
}

/* ---------------------------------------- Cancellation ---------------------------------------- */
// cancelJobTree returns the exec.Cancel function for the managed job. A successful termination is
// recorded on state; a tree whose child has already exited yields os.ErrProcessDone so Run does not
// inject a context error after a natural exit and does not record a kill.
func cancelJobTree(
	command *exec.Cmd,
	job windows.Handle,
	state *killState,
) (cancel func() (err error)) {
	return func() (err error) {
		if command.Process == nil || command.ProcessState != nil {
			return os.ErrProcessDone
		}

		// TerminateJobObject kills every process currently in the job. If afterStart has not bound
		// the suspended child yet this is a no-op; afterStart re-checks wasKilled after binding and
		// refuses to resume a cancelled child, so the kill is never lost to the bind race.
		if err = windows.TerminateJobObject(job, jobExitCode); err != nil {
			return err
		}

		state.markKilled()
		return nil
	}
}

/* ----------------------------------- Suspended Child Binding ---------------------------------- */
// bindSuspendedChild returns the afterStart hook that makes the suspended child part of the job
// before it is allowed to run. It opens the child by PID, assigns it to the job, then resumes the
// child's threads. Because assignment happens before resume, any descendant the child spawns is
// born inside the job.
func bindSuspendedChild(job windows.Handle, state *killState) (hook func(cmd *exec.Cmd) error) {
	return func(cmd *exec.Cmd) error {
		if err := assignSuspendedToJob(job, cmd.Process.Pid); err != nil {
			return err
		}

		// Cancellation may have raced with binding: cancelJobTree ran while the job was still empty
		// (a no-op terminate) and recorded a kill. The suspended child is now in the job, so
		// terminate it for real and never resume, so the cancelled child does not execute. Wait
		// reaps the dead child and the normal killState path attributes the cause.
		if state.wasKilled() {
			_ = windows.TerminateJobObject(job, jobExitCode)
			return nil
		}

		if err := resumeProcessThreads(uint32(cmd.Process.Pid)); err != nil {
			// Cancellation may instead have landed between the check above and resume. If so the
			// child is bound and has been killed by the job; let the normal flow attribute it.
			if state.wasKilled() {
				_ = windows.TerminateJobObject(job, jobExitCode)
				return nil
			}
			return err
		}

		return nil
	}
}

// assignSuspendedToJob opens the suspended child by PID and assigns it to the job. Opening by PID
// avoids relying on the unexported handle os.Process keeps on Windows. The child is still suspended
// on return; resumeProcessThreads starts it.
func assignSuspendedToJob(job windows.Handle, pid int) (err error) {
	handle, err := windows.OpenProcess(childProcessAccess, false, uint32(pid))
	if err != nil {
		return fmt.Errorf("open child process: %w", err)
	}
	defer windows.CloseHandle(handle) // the duplicate exists only for the assignment call

	if err = windows.AssignProcessToJobObject(job, handle); err != nil {
		return fmt.Errorf("assign process to job: %w", err)
	}

	return nil
}

/* -------------------------------------- Thread Resumption ------------------------------------- */
// resumeProcessThreads resumes every thread belonging to pid. A freshly created suspended process
// has exactly one thread, but resuming all of its threads is safe and correct: assignment has
// already placed the process in the job, so anything it does from here is captured.
func resumeProcessThreads(pid uint32) (err error) {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPTHREAD, 0)
	if err != nil {
		return fmt.Errorf("snapshot threads: %w", err)
	}
	defer windows.CloseHandle(snapshot) // read-only enumeration data: closing cannot affect results

	var entry windows.ThreadEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))

	if err = windows.Thread32First(snapshot, &entry); err != nil {
		if errors.Is(err, windows.ERROR_NO_MORE_FILES) {
			return fmt.Errorf("no threads for process %d: %w", pid, err)
		}
		return fmt.Errorf("enumerate threads: %w", err)
	}

	resumed := 0
	for {
		if entry.OwnerProcessID == pid {
			if err = resumeThread(entry.ThreadID); err != nil {
				return err
			}
			resumed++
		}

		if err = windows.Thread32Next(snapshot, &entry); err != nil {
			if errors.Is(err, windows.ERROR_NO_MORE_FILES) {
				break
			}
			return fmt.Errorf("enumerate threads: %w", err)
		}
	}

	if resumed == 0 {
		// The suspended child's main thread was not present in the snapshot. Returning an error
		// makes RunCommand fail closed (kill the suspended child) instead of leaving it suspended
		// and hanging Wait.
		return fmt.Errorf("no threads found for process %d", pid)
	}

	return nil
}

// resumeThread opens a thread by ID with the right to alter its suspend state and resumes it.
func resumeThread(threadID uint32) (err error) {
	handle, err := windows.OpenThread(windows.THREAD_SUSPEND_RESUME, false, threadID)
	if err != nil {
		return fmt.Errorf("open thread: %w", err)
	}
	defer windows.CloseHandle(handle) // the handle exists only for the resume call above

	if _, err = windows.ResumeThread(handle); err != nil {
		return fmt.Errorf("resume thread: %w", err)
	}

	return nil
}
