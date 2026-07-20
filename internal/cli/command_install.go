package cli

import (
	"context"
	"flag"
	"fmt"

	"github.com/wbd2023/Quill/internal/engine"
	"github.com/wbd2023/Quill/internal/report"
)

func runInstall(tool Tool, options installOptions) (exitCode int) {
	installer, err := engine.New(
		options.repoRoot,
		engine.WithProgressWriter(tool.stdout),
	)
	if err != nil {
		tool.writeError(err)
		return 1
	}

	result, err := installer.Install(context.Background())
	if err != nil {
		tool.writeError(err)
		return 1
	}

	toolchainResult := report.ToolchainResult{Statuses: result.Toolchain.Statuses}
	if _, err = renderToolchainStatus(tool.stdout, report.FormatText, toolchainResult); err != nil {
		tool.writeError(err)
		return 1
	}

	if !result.Toolchain.AllValid {
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
