package cli

import (
	"context"
	"flag"
	"fmt"

	"github.com/wbd2023/quill/internal/engine"
	"github.com/wbd2023/quill/internal/report"
)

func runInstall(ctx context.Context, tool Tool, options installOptions) (exitCode int) {
	progressWriter := tool.stdout
	if options.format == report.FormatJSON {
		// Machine mode reserves stdout for the single envelope; route install progress to stderr.
		progressWriter = tool.stderr
	}

	engineInstance, err := tool.buildEngine(
		options.repoRoot, engine.WithProgressWriter(progressWriter),
	)
	if err != nil {
		return tool.reportCommandError(ctx, "install", options.format, err)
	}

	result, err := engineInstance.Install(ctx)
	if err != nil {
		return tool.reportCommandError(ctx, "install", options.format, err)
	}

	toolchainResult := report.ToolchainResult{Statuses: result.Toolchain.Statuses}
	if _, err = renderToolchainStatus(
		tool.stdout, "install", options.format, toolchainResult,
	); err != nil {
		return tool.reportCommandError(ctx, "install", options.format, err)
	}

	if !result.Toolchain.AllValid {
		return 1
	}

	if options.format == report.FormatJSON {
		return 0
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
	var format string
	flagSet := newInstallFlagSet(&options, &format)
	if err = parseArguments(flagSet, summary, arguments); err != nil {
		return options, err
	}

	options.format, err = parseFormat(format)
	if err != nil {
		return options, err
	}

	options.repoRoot, err = resolve(options.repoRoot)
	return options, err
}

func newInstallFlagSet(options *installOptions, format *string) (flagSet *flag.FlagSet) {
	flagSet = newFlagSet("install")
	flagSet.StringVar(
		&options.repoRoot,
		"repo-root",
		"",
		"repository root (auto-detected when omitted)",
	)
	flagSet.StringVar(format, "format", string(report.FormatText), "format: text|json")
	return flagSet
}

func installUsageText() (usage string) {
	const summary = "install pinned style tools"
	var options installOptions
	var format string
	return commandUsage("install", summary, newInstallFlagSet(&options, &format))
}
