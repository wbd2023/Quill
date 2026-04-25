package cli

import (
	"fmt"
	"io"
	"strings"
)

func (tool CLI) runHelp(arguments []string) (exitCode int) {
	if len(arguments) == 0 {
		_, _ = io.WriteString(tool.stdout, rootUsageText())
		return 0
	}

	command, found := findCommand(arguments[0])
	if !found {
		tool.writeUsageError(rootUsageText(), fmt.Errorf("unknown command %q", arguments[0]))
		return usageExitCode
	}

	if len(arguments) > 1 {
		tool.writeUsageError(
			command.usage(),
			fmt.Errorf("unexpected arguments for help: %s", strings.Join(arguments[1:], ", ")),
		)
		return usageExitCode
	}

	_, _ = io.WriteString(tool.stdout, command.usage())
	return 0
}

func isHelpRequest(argument string) (help bool) {
	return argument == helpCommand || argument == "-h" || argument == "--help"
}

func rootUsageText() (usage string) {
	lines := []string{
		"usage:",
		"  style <command> [flags]",
		"",
		"commands:",
	}
	for _, command := range commands {
		lines = append(lines, fmt.Sprintf("  %-9s %s", command.name, command.summary))
	}

	lines = append(
		lines,
		"",
		"run `style help <command>` or `style <command> -h` to see command-specific flags",
		"",
	)
	return strings.Join(lines, "\n")
}
