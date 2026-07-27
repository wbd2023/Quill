package cli

import (
	"context"
	"flag"
	"io"

	"github.com/wbd2023/quill/internal/coverage"
	"github.com/wbd2023/quill/internal/report"
)

/* -------------------------------------- Coverage Command -------------------------------------- */

func runCoverage(ctx context.Context, tool Tool, options coverageOptions) (exitCode int) {
	engineInstance, err := tool.buildEngine(options.repoRoot)
	if err != nil {
		return tool.reportCommandError(ctx, "coverage", options.format, err)
	}

	coverageReport, err := engineInstance.Coverage(ctx)
	if err != nil {
		return tool.reportCommandError(ctx, "coverage", options.format, err)
	}

	if err = writeCoverageResult(tool.stdout, "coverage", coverageReport, options); err != nil {
		return tool.reportCommandError(ctx, "coverage", options.format, err)
	}

	return 0
}

/* --------------------------------------- Option Parsing --------------------------------------- */

func parseCoverageOptions(arguments []string) (options coverageOptions, err error) {
	return parseCoverageOptionsWithResolver(resolveRepoRoot, arguments)
}

func parseCoverageOptionsWithResolver(
	resolve repositoryRootResolver,
	arguments []string,
) (options coverageOptions, err error) {
	const summary = "show STYLE.md automation coverage"
	var format string
	flagSet := newCoverageFlagSet(&options, &format)
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

func newCoverageFlagSet(
	options *coverageOptions,
	format *string,
) (flagSet *flag.FlagSet) {
	flagSet = newFlagSet("coverage")
	flagSet.StringVar(
		&options.repoRoot,
		"repo-root",
		"",
		"repository root (auto-detected when omitted)",
	)
	flagSet.StringVar(format, "format", string(report.FormatText), "format: text|json")
	flagSet.BoolVar(&options.verbose, "verbose", false, "print requirement-level detail")
	return flagSet
}

func coverageUsageText() (usage string) {
	const summary = "show STYLE.md automation coverage"
	var options coverageOptions
	var format string
	return commandUsage("coverage", summary, newCoverageFlagSet(&options, &format))
}

/* ------------------------------------------ Rendering ----------------------------------------- */

func writeCoverageResult(
	writer io.Writer,
	command string,
	coverageReport coverage.Report,
	options coverageOptions,
) (err error) {
	view := report.NewCoverageView(coverageReport)
	return report.WriteCoverage(writer, command, options.format, view, options.verbose)
}
