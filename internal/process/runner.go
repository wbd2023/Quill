package process

import "context"

// Runner resolves and executes commands using the operating system. It satisfies
// toolchain.CommandRunner without importing toolchain (structural typing).
type Runner struct{}

// ResolvePath finds the full path to command, searching the PATH in environment.
func (Runner) ResolvePath(
	ctx context.Context,
	environment map[string]string,
	command string,
) (path string, err error) {
	return ResolveCommandPath(environment, command)
}

// Run executes the binary at path with arguments, using environment, and returns its combined
// output.
func (Runner) Run(
	ctx context.Context,
	environment map[string]string,
	path string,
	arguments []string,
) (output string, err error) {
	result, err := RunCommand(ctx, CommandRequest{
		Name:        path,
		Arguments:   arguments,
		Environment: environment,
	})
	if err != nil {
		return "", err
	}

	return result.Output, nil
}
