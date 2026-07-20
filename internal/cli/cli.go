package cli

import (
	"errors"
	"fmt"
	"io"
)

const (
	helpCommand   = "help"
	usageExitCode = 2
)

type repositoryRootResolver func(string) (string, error)

// Tool runs CLI commands using the configured output streams.
type Tool struct {
	stdout          io.Writer
	stderr          io.Writer
	resolveRepoRoot repositoryRootResolver
	version         string
}

// Action is a parsed CLI command ready to run.
type Action func(Tool) int

// New constructs a CLI tool with the given output streams and build version.
func New(stdout io.Writer, stderr io.Writer, version string) (tool Tool) {
	if stdout == nil {
		stdout = io.Discard
	}
	if stderr == nil {
		stderr = io.Discard
	}

	return Tool{
		stdout:          stdout,
		stderr:          stderr,
		resolveRepoRoot: resolveRepoRoot,
		version:         version,
	}
}

// Run parses and executes one CLI command.
func (tool Tool) Run(arguments []string) (exitCode int) {
	if len(arguments) == 0 {
		tool.writeUsageError(rootUsageText(), nil)
		return usageExitCode
	}

	if isHelpRequest(arguments[0]) {
		return tool.runHelp(arguments[1:])
	}

	command, found := findCommand(arguments[0])
	if !found {
		tool.writeUsageError(rootUsageText(), fmt.Errorf("unknown command %q", arguments[0]))
		return usageExitCode
	}

	action, err := command.prepare(tool.resolveRepoRoot, arguments[1:])
	if err == nil {
		return action(tool)
	}

	var help flagHelpError
	if errors.As(err, &help) {
		_, _ = io.WriteString(tool.stdout, help.message)
		return 0
	}

	tool.writeUsageError(command.usage(), err)
	return usageExitCode
}
