package toolchain

type CommandRequest struct {
	Name        string
	Arguments   []string
	Environment map[string]string
}

type CommandRunner func(request CommandRequest) (output string, err error)
