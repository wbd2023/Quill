package cli

import "fmt"

const versionSummary = "print the Quill version"

type versionOptions struct{}

func parseVersionOptions(
	_ repositoryRootResolver,
	arguments []string,
) (options versionOptions, err error) {
	flagSet := newFlagSet("version")
	return options, parseArguments(flagSet, versionSummary, arguments)
}

func versionUsageText() (usage string) {
	return commandUsage("version", versionSummary, newFlagSet("version"))
}

func runVersion(tool Tool, _ versionOptions) (exitCode int) {
	_, _ = fmt.Fprintln(tool.stdout, tool.version)
	return 0
}
