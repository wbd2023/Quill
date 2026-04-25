package cli

import (
	"flag"
	"io"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/executors"
	"ciphera/tools/internal/report"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runtime"
)

/* ---------------------------------------- Check Command --------------------------------------- */

func runCheck(tool CLI, options checkOptions) (exitCode int) {
	context, err := loadContext(options.repoRoot, options.scope)
	if err != nil {
		tool.writeError(err)
		return 1
	}

	selected, err := selectedRules(context.Effective.Rules, options.scope, options.profile)
	if err != nil {
		tool.writeError(err)
		return 1
	}
	toolStatusList := runtime.InspectToolsWithEnvironment(
		context.Effective.Tools,
		runner.ToolIDsForRules(selected),
		context.ToolEnvironment,
	)
	toolStatuses := runtime.StatusesByID(toolStatusList)

	result := report.CheckResult{
		Entries: make([]report.CheckEntry, 0, len(selected)),
	}
	checkers := executors.Checkers()
	for _, rule := range selected {
		output, err := runner.RunRule(rule, context, toolStatuses, checkers)
		entry := report.CheckEntry{
			Rule:   rule,
			Status: statusForRuleResult(rule, err, options.strictRecommendations),
			Output: output,
		}
		result.Entries = append(result.Entries, entry)
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
	var profile string
	var format string
	flagSet := newCheckFlagSet(&options, &scope, &profile, &format)

	if err = parseArguments(flagSet, summary, arguments); err != nil {
		return options, err
	}

	options.scope, err = parseScope(scope)
	if err != nil {
		return options, err
	}

	options.profile, err = parseProfile(profile)
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
	profile *string,
	format *string,
) (flagSet *flag.FlagSet) {
	flagSet = newFlagSet("check")
	flagSet.StringVar(
		&options.repoRoot,
		"repo-root",
		"",
		"repository root (auto-detected when omitted)",
	)
	flagSet.StringVar(scope, "scope", string(contract.ScopeAll), "scope: app|tools|all")
	flagSet.StringVar(
		profile,
		"profile",
		string(contract.CheckProfileRequired),
		"profile: required|all",
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
	var profile string
	var format string
	return commandUsage("check", summary, newCheckFlagSet(&options, &scope, &profile, &format))
}

/* --------------------------------------- Rule Selection --------------------------------------- */

func selectedRules(
	available []contract.Rule,
	scope contract.Scope,
	profile contract.CheckProfile,
) (rules []contract.Rule, err error) {
	for _, rule := range available {
		if !contract.ScopeCovers(scope, rule.Scope) {
			continue
		}

		if profile == contract.CheckProfileRequired && rule.Level == contract.LevelRecommendation {
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
	rule contract.Rule,
	err error,
	strictRecommendations bool,
) (status report.CheckStatus) {
	return runner.CheckStatus(rule, err, strictRecommendations)
}
