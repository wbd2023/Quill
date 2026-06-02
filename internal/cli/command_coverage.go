package cli

import (
	"flag"
	"io"

	"ciphera/tools/internal/coverage"
	"ciphera/tools/internal/pack/builtin"
	"ciphera/tools/internal/profile"
	"ciphera/tools/internal/report"
	"ciphera/tools/internal/styleguide"
)

/* -------------------------------------- Coverage Command -------------------------------------- */

func runCoverage(tool CLI, options coverageOptions) (exitCode int) {
	coverageReport, err := loadCoverageReport(options.repoRoot)
	if err != nil {
		tool.writeError(err)
		return 1
	}

	if err = writeCoverageResult(tool.stdout, coverageReport, options); err != nil {
		tool.writeError(err)
		return 1
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

/* -------------------------------------- Coverage Loading -------------------------------------- */

func loadCoverageReport(repoRoot string) (coverageReport coverage.Report, err error) {
	config, err := profile.Load(repoRoot)
	if err != nil {
		return coverage.Report{}, err
	}

	document, err := styleguide.Load(repoRoot, styleguide.Config{
		Filename: config.StyleGuide.Path,
		IDScheme: config.StyleGuide.IDScheme,
	})
	if err != nil {
		return coverage.Report{}, err
	}

	registry, err := builtin.DefaultRegistry(config.EnabledPacks)
	if err != nil {
		return coverage.Report{}, err
	}

	compiled, err := profile.Compile(config, registry)
	if err != nil {
		return coverage.Report{}, err
	}

	return coverage.Build(document, compiled.Effective.Rules), nil
}

/* ------------------------------------------ Rendering ----------------------------------------- */

func writeCoverageResult(
	writer io.Writer,
	coverageReport coverage.Report,
	options coverageOptions,
) (err error) {
	view := report.NewCoverageView(coverageReport)
	return report.WriteCoverage(writer, options.format, view, options.verbose)
}
