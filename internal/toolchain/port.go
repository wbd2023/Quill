package toolchain

// CommandRunner resolves and executes commands for tool version inspection.
type CommandRunner interface {
	// ResolvePath finds the full path to command, searching the PATH in environment.
	ResolvePath(environment map[string]string, command string) (path string, err error)

	// Run executes the binary at path with arguments, using environment, and returns its
	// combined output.
	Run(environment map[string]string, path string, arguments []string) (output string, err error)
}
