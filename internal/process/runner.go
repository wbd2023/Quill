package process

// Runner resolves and executes commands using the operating system. It satisfies
// toolchain.CommandRunner without importing toolchain (structural typing).
type Runner struct{}

// ResolvePath finds the full path to command, searching the PATH in environment.
func (Runner) ResolvePath(environment map[string]string, command string) (path string, err error) {
	return ResolveCommandPath(environment, command)
}

// Run executes the binary at path with arguments, using environment, and returns its combined
// output.
func (Runner) Run(
	environment map[string]string,
	path string,
	arguments []string,
) (output string, err error) {
	result, err := RunCommand(CommandRequest{
		Name:        path,
		Arguments:   arguments,
		Environment: environment,
	})
	if err != nil {
		return "", err
	}

	return result.Output, nil
}
