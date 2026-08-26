package process

import (
	"bytes"
	"context"
	"errors"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"time"
)

/* ------------------------------------------ Constants ----------------------------------------- */

// Default limits applied when a request omits its own timeout or output limit.
const (
	defaultCommandTimeout   = 120 * time.Second
	defaultOutputLimitBytes = 1 << 20
)

const anyExecutePermission os.FileMode = 0o111

/* ------------------------------------- Environment Policy ------------------------------------- */

// defaultWindowsPathExt mirrors the PATHEXT fallback that os/exec applies on Windows when the
// variable is absent, so bare-name resolution stays consistent with the operating system.
const defaultWindowsPathExt = ".COM;.EXE;.BAT;.CMD;.VBS;.VBE;.JS;.JSE;.WSF;.WSH;.MSC"

const (
	// EnvironmentInherit starts from the parent process environment and applies Variables as
	// explicit overrides (Variables win). Use when the child needs ambient operating-system state.
	EnvironmentInherit EnvironmentPolicy = iota
	// EnvironmentIsolated exposes only Variables to the child; no parent environment is inherited.
	// Use for hermetic, deliberate tool execution that must not leak sensitive ambient state, such
	// as the external Pack subprocess protocol.
	EnvironmentIsolated
)

// EnvironmentPolicy controls how a request's Variables combine with the parent process
// environment when launching a child process. Making the policy explicit removes the prior
// ambiguity where every map silently inherited the full parent environment.
type EnvironmentPolicy int

/* -------------------------------------------- Types ------------------------------------------- */

// CommandRequest describes a command to execute: the executable name, direct arguments,
// environment policy and variables, working directory, standard input, and execution limits.
// Arguments are passed directly to the child process; no shell is ever invoked.
type CommandRequest struct {
	Name      string
	Arguments []string

	// Environment declares how Variables combine with the parent environment. The zero value
	// (EnvironmentInherit) preserves ambient operating-system state for the child.
	Environment EnvironmentPolicy
	// Variables are the explicit environment entries, interpreted according to Environment.
	Variables map[string]string

	// Directory is the working directory for the child. Empty means the current directory.
	Directory string

	// Stdin is the bytes piped to the child's standard input. Empty means the child reads from
	// the null device; the parent's standard input is never inherited.
	Stdin []byte

	Timeout          time.Duration
	OutputLimitBytes int64
}

// CommandResult is the outcome of running a command. Stdout, Stderr, and Output are each bounded
// to the request's OutputLimitBytes.
type CommandResult struct {
	// Output is the combined standard-output and standard-error stream, bounded to
	// OutputLimitBytes. It preserves the interleaved legacy interpretation relied on by existing
	// diagnostic interpreters.
	Output string
	// Stdout is the bounded standard-output stream, captured separately from standard error.
	Stdout string
	// Stderr is the bounded standard-error stream, captured separately from standard output.
	Stderr string

	ExitCode int

	// TimedOut reports whether the request timeout fired. It is distinct from Canceled.
	TimedOut bool
	// Canceled reports whether the caller's context was cancelled before the child exited.
	Canceled bool
	// Truncated reports whether any captured stream (stdout, stderr, or combined) reached the
	// output limit.
	Truncated bool
}

/* -------------------------------------- Command Execution ------------------------------------- */

// RunCommand resolves and executes the command described by request, applying its timeout and
// output limit, and returns the captured streams and exit status.
//
// The caller's context is honoured: parent cancellation propagates to the child and surfaces as
// context.Canceled, while the request timeout surfaces as context.DeadlineExceeded. The two are
// never collapsed into a generic command error. A non-zero exit is reported as a CommandError
// wrapping the underlying exit error. The executable is resolved once, before start, using an
// OS-aware PATH and PATHEXT search; the resolved path is passed directly to the child.
func RunCommand(
	parent context.Context,
	request CommandRequest,
) (result CommandResult, err error) {
	path, err := ResolveExecutable(request.Variables, request.Name)
	if err != nil {
		return CommandResult{}, err
	}

	timeoutCtx, cancel := context.WithTimeout(parent, resolveTimeout(request.Timeout))
	defer cancel()

	command := exec.CommandContext(timeoutCtx, path, request.Arguments...)
	command.Dir = resolveDirectory(request.Directory)
	command.Env = buildEnvironment(request.Environment, request.Variables)
	if len(request.Stdin) > 0 {
		command.Stdin = bytes.NewReader(request.Stdin)
	}
	tracker, afterStart, release, err := configureProcessTree(command)
	if err != nil {
		// The process tree cannot be tracked safely (for example the Windows Job Object could not
		// be established). Fail closed before the child is ever started, so no untrusted
		// instruction runs outside the managed tree.
		result = CommandResult{ExitCode: resolveExitCode(err)}
		return result, classifyError(request, err, nil, result)
	}

	limit := resolveOutputLimit(request.OutputLimitBytes)
	stdout, stderr, combined := newOutputBuffers(limit)
	command.Stdout = &streamSink{stream: stdout, combined: combined}
	command.Stderr = &streamSink{stream: stderr, combined: combined}

	// Start and Wait are split (rather than command.Run) so the Windows build can bind the
	// suspended child to its Job Object and resume it between the two calls; that ordering is what
	// makes descendant termination race-safe. POSIX returns no afterStart hook, so the split is a
	// no-op there.
	runErr := command.Start()
	if runErr != nil {
		release()
		result = CommandResult{ExitCode: resolveExitCode(runErr)}
		return result, classifyError(request, runErr, parent.Err(), result)
	}

	if afterStart != nil {
		if bindErr := afterStart(command); bindErr != nil {
			// The suspended child could not be bound to the managed tree. It has not executed, so
			// terminate it and fail closed before any untrusted instruction runs.
			_ = command.Process.Kill()
			_ = command.Wait()
			release()
			result = CommandResult{ExitCode: resolveExitCode(bindErr)}
			return result, classifyError(request, bindErr, nil, result)
		}
	}
	defer release()

	runErr = command.Wait()

	// Termination cause is attributed from whether the cancellation function actually killed a
	// running process, not from a process exit code: a Unix signal death reports a negative exit
	// code while a Windows kill exits with code 1, so the exit code is not a reliable
	// cross-platform signal. A natural exit (success or a real exit code) is never attributed
	// to the context, even if the timer or cancellation arrived moments later.
	forced := tracker.wasKilled()
	parentErr := parent.Err()
	timeoutErr := timeoutCtx.Err()

	result = CommandResult{
		Output:    combined.String(),
		Stdout:    stdout.String(),
		Stderr:    stderr.String(),
		ExitCode:  resolveExitCode(runErr),
		TimedOut:  forced && errors.Is(timeoutErr, context.DeadlineExceeded) && parentErr == nil,
		Canceled:  forced && errors.Is(parentErr, context.Canceled),
		Truncated: stdout.truncated || stderr.truncated || combined.truncated,
	}

	if runErr == nil {
		return result, nil
	}

	return result, classifyError(request, runErr, parentErr, result)
}

// classifyError wraps the failure in a CommandError whose underlying error distinguishes parent
// cancellation, Quill's request timeout, the caller's own context expiry, and ordinary exit
// failure. The context sentinels are preserved so callers can errors.Is against context.Canceled
// and context.DeadlineExceeded.
func classifyError(
	request CommandRequest,
	runErr error,
	parentErr error,
	result CommandResult,
) (err error) {
	switch {
	case result.Canceled:
		err = context.Canceled

	case result.TimedOut:
		err = context.DeadlineExceeded

	case parentErr != nil:
		// The caller's context expired on its own (for example a parent deadline) and terminated
		// the child before Quill's request timer. Surface that cause rather than
		// the injected error.
		err = parentErr

	default:
		err = runErr
	}

	return CommandError{
		Name:      request.Name,
		Arguments: slices.Clone(request.Arguments),
		Err:       err,
	}
}

/* ------------------------------------ Executable Resolution ----------------------------------- */

// ResolveExecutable resolves command to an executable path, honouring the PATH in variables when
// set and otherwise falling back to the operating system's native resolution. A direct path
// (absolute or containing a separator) is verified to be executable and, when it lacks an
// extension, completed through PATHEXT rather than handed to os/exec to requalify at start. Bare
// names are searched across PATH, applying PATHEXT on Windows so they resolve to a qualified
// executable.
func ResolveExecutable(
	variables map[string]string,
	command string,
) (resolved string, err error) {
	if command == "" {
		return "", exec.ErrNotFound
	}

	extensions := executableExtensions(resolvePathExt(lookupVariable(variables, "PATHEXT")))

	searchPath := lookupVariable(variables, "PATH")
	if searchPath == "" {
		return exec.LookPath(command)
	}

	if isDirectPath(command) {
		return resolveDirectPath(command, extensions)
	}

	for _, directory := range filepath.SplitList(searchPath) {
		if directory == "" {
			continue
		}

		if resolved = findExecutable(directory, command, extensions); resolved != "" {
			return resolved, nil
		}
	}

	return "", exec.ErrNotFound
}

func isDirectPath(command string) (direct bool) {
	return filepath.IsAbs(command) || strings.ContainsAny(command, `/\`)
}

// resolveDirectPath verifies that a direct command path is executable and returns the final path.
// A bare path that exists and is executable is returned as-is; an extensionless path is completed
// through PATHEXT (so /tools/go resolves to /tools/go.exe and os/exec does not requalify it at
// start). A path that cannot be resolved to an executable is reported as not found rather than
// returned unchecked.
func resolveDirectPath(command string, extensions []string) (resolved string, err error) {
	if resolved = executableCandidate(command); resolved != "" {
		return resolved, nil
	}

	if filepath.Ext(command) == "" {
		for _, extension := range extensions {
			if resolved = executableCandidate(command + extension); resolved != "" {
				return resolved, nil
			}
		}
	}

	return "", exec.ErrNotFound
}

// executableExtensions parses PATHEXT into a list of extensions, each normalised to carry a leading
// dot so that malformed entries such as EXE or cmd resolve go.exe and tool.cmd instead of the
// concatenation goEXE.
func executableExtensions(pathExt string) (extensions []string) {
	if pathExt == "" {
		return nil
	}

	for _, extension := range strings.Split(pathExt, ";") {
		if extension = strings.TrimSpace(extension); extension != "" {
			if !strings.HasPrefix(extension, ".") {
				extension = "." + extension
			}
			extensions = append(extensions, extension)
		}
	}

	return extensions
}

// resolvePathExt returns the PATHEXT value to search with. An explicit value is honoured as-is;
// an empty value falls back to the Windows default so resolution does not silently degrade, and to
// the empty string on POSIX where PATHEXT is meaningless.
func resolvePathExt(pathExt string) (resolved string) {
	if pathExt != "" {
		return pathExt
	}

	if runtime.GOOS == "windows" {
		return defaultWindowsPathExt
	}

	return ""
}

func findExecutable(directory string, command string, extensions []string) (resolved string) {
	base := filepath.Join(directory, command)

	// A command that already carries an extension is fully qualified, and a search with no
	// PATHEXT extensions (POSIX) has nothing else to try, so the bare candidate is acceptable
	// in those cases. On Windows a bare name is resolved only through PATHEXT so that an
	// extensionless file is never preferred over a qualified executable such as go.exe,
	// matching os/exec.
	if len(extensions) == 0 || filepath.Ext(command) != "" {
		if resolved = executableCandidate(base); resolved != "" {
			return resolved
		}
	}

	for _, extension := range extensions {
		if resolved = executableCandidate(base + extension); resolved != "" {
			return resolved
		}
	}

	return ""
}

func executableCandidate(path string) (resolved string) {
	if isExecutableFile(path) {
		return path
	}

	return ""
}

/* ------------------------------------------- Helpers ------------------------------------------ */

// buildEnvironment materialises the child environment for policy. Entries are merged under their
// canonical key (see envKey) so Variables override inherited values deterministically. On Windows
// this is case-insensitive: a Variables PATH overrides an inherited Path. On POSIX, where names
// are case-sensitive, the key is preserved verbatim.
func buildEnvironment(
	policy EnvironmentPolicy,
	variables map[string]string,
) (environment []string) {
	combined := make(map[string]string, len(variables))

	if policy != EnvironmentIsolated {
		for _, entry := range os.Environ() {
			if key, value, ok := splitEnvEntry(entry); ok {
				combined[envKey(key)] = value
			}
		}
	}

	for key, value := range variables {
		combined[envKey(key)] = value
	}

	return sortedEntries(combined)
}

// envKey canonicalises an environment-variable name for merge deduplication. Windows environment
// names are case-insensitive, so PATH and Path refer to the same variable; folding to upper case
// lets a Variables entry override an inherited entry regardless of the parent's casing. POSIX
// names are case-sensitive and are returned untouched.
func envKey(key string) (canonical string) {
	if runtime.GOOS == "windows" {
		return strings.ToUpper(key)
	}

	return key
}

func sortedEntries(variables map[string]string) (entries []string) {
	entries = make([]string, 0, len(variables))
	for _, key := range slices.Sorted(maps.Keys(variables)) {
		entries = append(entries, key+"="+variables[key])
	}

	return entries
}

// splitEnvEntry splits an environment entry into key and value. Entries that begin with '=' are
// Windows per-drive current-directory pseudo-variables (for example "=C:=C:\work"); their key runs
// to the second '=', so each drive stays distinct instead of collapsing to an empty key. This
// mirrors how os/exec parses environment entries for deduplication.
func splitEnvEntry(entry string) (key string, value string, ok bool) {
	if entry == "" {
		return "", "", false
	}

	if entry[0] == '=' {
		separator := strings.IndexByte(entry[1:], '=')
		if separator < 0 {
			return "", "", false
		}

		separator++
		return entry[:separator], entry[separator+1:], true
	}

	separator := strings.IndexByte(entry, '=')
	if separator < 0 {
		return "", "", false
	}

	return entry[:separator], entry[separator+1:], true
}

// lookupVariable returns the named variable from variables when present, otherwise from the parent
// process environment. The parent fallback aids executable resolution only; it never leaks into the
// child environment under EnvironmentIsolated.
func lookupVariable(variables map[string]string, name string) (value string) {
	if variables != nil {
		if value = variables[name]; value != "" {
			return value
		}
	}

	return os.Getenv(name)
}

func resolveDirectory(directory string) (resolved string) {
	if directory == "" {
		return "."
	}

	return directory
}

func resolveTimeout(timeout time.Duration) (duration time.Duration) {
	if timeout <= 0 {
		return defaultCommandTimeout
	}

	return timeout
}

func resolveOutputLimit(limit int64) (resolved int64) {
	if limit <= 0 {
		return defaultOutputLimitBytes
	}

	return limit
}

func resolveExitCode(err error) (code int) {
	if err == nil {
		return 0
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}

	return -1
}

func isExecutableFile(path string) (found bool) {
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return false
	}

	if runtime.GOOS == "windows" {
		return true
	}

	return info.Mode()&anyExecutePermission != 0
}
