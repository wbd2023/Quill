package cli

import (
	"errors"
	"fmt"
	"io"
)

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
