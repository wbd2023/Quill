package cli

import (
	"flag"
	"io"

	"ciphera/tools/internal/report"
	"ciphera/tools/internal/styleguide"
)

func runCoverage(tool CLI, options coverageOptions) (exitCode int) {
	coverage, err := styleguide.Coverage(options.repoRoot)
	if err != nil {
		tool.writeError(err)
		return 1
	}

	if err = writeCoverageResult(tool.stdout, coverage, options); err != nil {
		tool.writeError(err)
		return 1
	}

	return 0
}

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

func writeCoverageResult(
	writer io.Writer,
	coverage styleguide.CoverageReport,
	options coverageOptions,
) (err error) {
	view := report.NewCoverageView(coverage)
	return report.WriteCoverage(writer, options.format, view, options.verbose)
}
