package cli

import (
	"flag"
	"fmt"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/report"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runtime"
)

func parseInstallOptions(arguments []string) (options installOptions, err error) {
	return parseInstallOptionsWithResolver(resolveRepoRoot, arguments)
}

func parseInstallOptionsWithResolver(
	resolve repositoryRootResolver,
	arguments []string,
) (options installOptions, err error) {
	const summary = "install pinned style tools"
	flagSet := newInstallFlagSet(&options)
	if err = parseArguments(flagSet, summary, arguments); err != nil {
		return options, err
	}

	options.repoRoot, err = resolve(options.repoRoot)
	return options, err
}

func newInstallFlagSet(options *installOptions) (flagSet *flag.FlagSet) {
	flagSet = newFlagSet("install")
	flagSet.StringVar(
		&options.repoRoot,
		"repo-root",
		"",
		"repository root (auto-detected when omitted)",
	)
	return flagSet
}

func installUsageText() (usage string) {
	const summary = "install pinned style tools"
	var options installOptions
	return commandUsage("install", summary, newInstallFlagSet(&options))
}

func runInstall(tool CLI, options installOptions) (exitCode int) {
	context, err := loadContext(options.repoRoot, contract.ScopeAll)
	if err != nil {
		tool.writeError(err)
		return 1
	}

	if err := runtime.Install(context.Layout, tool.stdout, context.Effective.Tools); err != nil {
		tool.writeError(err)
		return 1
	}

	toolIDs := toolIDsFromTools(context.Effective.Tools)
	statuses, allValid := runner.InspectToolchain(
		context.Effective.Tools,
		toolIDs,
		context.ToolEnvironment,
	)
	result := report.ToolchainResult{Statuses: statuses}
	if _, err = renderToolchainStatus(tool.stdout, report.FormatText, result); err != nil {
		tool.writeError(err)
		return 1
	}

	if !allValid {
		return 1
	}

	if _, err := fmt.Fprintln(tool.stdout, "Style tools installed."); err != nil {
		tool.writeError(err)
		return 1
	}

	return 0
}
