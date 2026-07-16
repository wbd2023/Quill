package cli

import (
	"context"
	"flag"
	"io"

	"ciphera/tools/internal/engine"
	"ciphera/tools/internal/report"
	"ciphera/tools/internal/style"
)

/* ---------------------------------------- Check Command --------------------------------------- */

func runCheck(tool Tool, options checkOptions) (exitCode int) {
	checker, err := engine.New(options.repoRoot)
	if err != nil {
		tool.writeError(err)
		return 1
	}

	result, err := checker.Check(context.Background(), engine.CheckOptions{
		Scope:                 options.scope,
		Mode:                  options.mode,
		StrictRecommendations: options.strictRecommendations,
	})
	if err != nil {
		tool.writeError(err)
		return 1
	}

	checkResult := report.CheckResult{
		Entries: make([]report.CheckEntry, 0, len(result.Rules)),
	}
	for _, ruleResult := range result.Rules {
		checkResult.Entries = append(
			checkResult.Entries,
			report.NewCheckEntry(
				ruleResult.Rule,
				ruleResult.Status,
				ruleResult.Execution,
			),
		)
	}

	summary, err := writeCheckResult(tool.stdout, checkResult, options)
	if err != nil {
		tool.writeError(err)
		return 1
	}

	if summary.Failed > 0 || summary.Errored > 0 {
		return 1
	}

	return 0
}

func parseCheckOptions(arguments []string) (options checkOptions, err error) {
	return parseCheckOptionsWithResolver(resolveRepoRoot, arguments)
}

func parseCheckOptionsWithResolver(
	resolve repositoryRootResolver,
	arguments []string,
) (options checkOptions, err error) {
	const summary = "run STYLE.md checks"
	var scope string
	var mode string
	var format string
	flagSet := newCheckFlagSet(&options, &scope, &mode, &format)

	if err = parseArguments(flagSet, summary, arguments); err != nil {
		return options, err
	}

	options.scope, err = parseScope(scope)
	if err != nil {
		return options, err
	}

	options.mode, err = parseCheckMode(mode)
	if err != nil {
		return options, err
	}

	options.format, err = parseFormat(format)
	if err != nil {
		return options, err
	}

	options.repoRoot, err = resolve(options.repoRoot)
	return options, err
}

func newCheckFlagSet(
	options *checkOptions,
	scope *string,
	mode *string,
	format *string,
) (flagSet *flag.FlagSet) {
	flagSet = newFlagSet("check")
	flagSet.StringVar(
		&options.repoRoot,
		"repo-root",
		"",
		"repository root (auto-detected when omitted)",
	)
	flagSet.StringVar(scope, "scope", "", "configured scope (profile default when omitted)")
	flagSet.StringVar(
		mode,
		"mode",
		string(style.CheckModeRequired),
		"mode: required|all",
	)
	flagSet.BoolVar(
		&options.strictRecommendations,
		"strict-recommendations",
		false,
		"fail on recommendation findings",
	)
	flagSet.StringVar(format, "format", string(report.FormatText), "format: text|json")
	flagSet.BoolVar(&options.verbose, "verbose", false, "print failing output")
	return flagSet
}

func checkUsageText() (usage string) {
	const summary = "run STYLE.md checks"
	var options checkOptions
	var scope string
	var mode string
	var format string
	return commandUsage("check", summary, newCheckFlagSet(&options, &scope, &mode, &format))
}

/* ------------------------------------------ Rendering ----------------------------------------- */

func writeCheckResult(
	writer io.Writer,
	result report.CheckResult,
	options checkOptions,
) (summary report.CheckSummary, err error) {
	view := report.NewCheckView(result)
	return report.WriteCheck(writer, options.format, view, options.verbose)
}
