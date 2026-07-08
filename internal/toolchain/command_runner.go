package toolchain

// CommandRequest names a command and its arguments plus the environment to run it in. Passed to
// a CommandRunner for version-detection probes.
type CommandRequest struct {
	Name        string
	Arguments   []string
	Environment map[string]string
}

// CommandRunner executes a CommandRequest and returns its combined stdout. Injected so the
// inspector can probe tool versions without depending on a concrete executor.
type CommandRunner func(request CommandRequest) (output string, err error)
