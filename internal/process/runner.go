package process

import "context"

// Runner resolves and executes commands using the operating system. It satisfies
// toolchain.CommandRunner without importing toolchain (structural typing). It is the single
// injectable OS-process boundary used for tool version inspection.
type Runner struct{}

// ResolvePath finds the full path to command, searching the PATH in environment. The context is
// accepted for interface conformance; path resolution is a fast filesystem operation.
func (Runner) ResolvePath(
	_ context.Context,
	environment map[string]string,
	command string,
) (path string, err error) {
	return ResolveExecutable(environment, command)
}

// Run executes the binary at path with arguments, using environment, and returns its combined
// output. A non-zero exit, timeout, or cancellation is returned as a CommandError that preserves
// the underlying cause.
func (Runner) Run(
	ctx context.Context,
	environment map[string]string,
	path string,
	arguments []string,
) (output string, err error) {
	result, err := RunCommand(ctx, CommandRequest{
		Name:        path,
		Arguments:   arguments,
		Environment: EnvironmentInherit,
		Variables:   environment,
	})
	if err != nil {
		return "", err
	}

	return result.Output, nil
}
