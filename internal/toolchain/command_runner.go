package toolchain

// CommandRequest is command request.
type CommandRequest struct {
	Name        string
	Arguments   []string
	Environment map[string]string
}

// CommandRunner is command runner.
type CommandRunner func(request CommandRequest) (output string, err error)
