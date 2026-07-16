package cli

import (
	"flag"
	"fmt"

	"ciphera/tools/internal/installer"
	"ciphera/tools/internal/lockfile"
	"ciphera/tools/internal/report"
	"ciphera/tools/internal/runtime"
)

func runInstall(tool Tool, options installOptions) (exitCode int) {
	context, err := loadContext(options.repoRoot, "")
	if err != nil {
		tool.writeError(err)
		return 1
	}

	layout := runtime.NewLayout(context.RepoRoot)
	loaded, err := lockfile.Load(context.RepoRoot)
	if err != nil {
		tool.writeError(err)
		return 1
	}

	if err := installer.Install(
		layout,
		tool.stdout,
		sortedTools(context.Tools),
		loaded,
	); err != nil {
		tool.writeError(err)
		return 1
	}

	statuses, allValid := inspectToolchain(
		runtime.Runner{},
		context.Tools,
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
