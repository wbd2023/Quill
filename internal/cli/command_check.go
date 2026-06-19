package cli

import (
	"flag"
	"io"

	"ciphera/tools/internal/pack/shipped/bindings"
	"ciphera/tools/internal/report"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runner/drivers"
	"ciphera/tools/internal/runtime"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/toolchain"
)

/* ---------------------------------------- Check Command --------------------------------------- */

func runCheck(tool Tool, options checkOptions) (exitCode int) {
	context, err := loadContext(options.repoRoot, options.scope)
	if err != nil {
		tool.writeError(err)
		return 1
	}

	selected, err := selectedRules(context.Effective.Rules, context, options.mode)
	if err != nil {
		tool.writeError(err)
		return 1
	}
	toolStatusList := toolchain.InspectToolsWithEnvironment(
		context.Effective.Tools,
		context.ToolCapabilities,
		runner.ToolIDsForRules(selected),
		context.ToolEnvironment,
		runtime.RunToolchainCommand,
	)
	toolStatuses := toolchain.StatusesByID(toolStatusList)

	result := report.CheckResult{
		Entries: make([]report.CheckEntry, 0, len(selected)),
	}
	checkers := drivers.CheckDrivers(bindings.Build())
	for _, rule := range selected {
		execution, err := runner.RunRule(rule, context, toolStatuses, checkers)
		result.Entries = append(
			result.Entries,
			report.NewCheckEntry(
				rule,
				statusForRuleResult(rule, err, options.strictRecommendations),
				execution,
			),
		)
	}

	summary, err := writeCheckResult(tool.stdout, result, options)
	if err != nil {
		tool.writeError(err)
		return 1
	}

	if summary.Failed > 0 {
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

/* --------------------------------------- Rule Selection --------------------------------------- */

func selectedRules(
	available []style.Rule,
	context runner.Context,
	mode style.CheckMode,
) (rules []style.Rule, err error) {
	for _, rule := range available {
		if !context.Profile.Repository.HasScopeOverlap(context.Scope, rule.Scope) {
			continue
		}

		if mode == style.CheckModeRequired &&
			rule.Enforcement == style.EnforcementRecommendation {
			continue
		}

		rules = append(rules, rule)
	}

	return rules, nil
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

/* --------------------------------------- Status Mapping --------------------------------------- */

func statusForRuleResult(
	rule style.Rule,
	err error,
	strictRecommendations bool,
) (status style.CheckStatus) {
	return runner.CheckStatus(rule, err, strictRecommendations)
}
