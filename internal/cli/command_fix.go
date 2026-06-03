package cli

import (
	"flag"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/report"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runner/drivers"
	"ciphera/tools/internal/toolchain"
)

/* ----------------------------------------- Fix Command ---------------------------------------- */

func runFix(tool CLI, options fixOptions) (exitCode int) {
	context, err := loadContext(options.repoRoot, options.scope)
	if err != nil {
		tool.writeError(err)
		return 1
	}

	rules := fixableRules(context.Effective.Rules, context)
	if len(rules) == 0 {
		return 0
	}

	toolIDs := runner.ToolIDsForFixes(rules)
	statuses, allValid := inspectToolchain(
		context.Effective.Tools,
		context.ToolCapabilities,
		toolIDs,
		context.ToolEnvironment,
	)
	result := report.ToolchainResult{Statuses: statuses}
	if _, err := renderToolchainStatus(tool.stderr, report.FormatText, result); err != nil {
		tool.writeError(err)
		return 1
	}

	if !allValid {
		return 1
	}

	statusIndex := toolchain.StatusesByID(statuses)
	fixers := drivers.FixDrivers()
	for _, rule := range rules {
		result, err := runner.RunFix(rule, context, statusIndex, fixers)
		if err != nil {
			tool.writeCommandOutput(result.Output)
			if result.Output == "" {
				tool.writeError(err)
			}
			return 1
		}
	}

	return 0
}

func parseFixOptions(arguments []string) (options fixOptions, err error) {
	return parseFixOptionsWithResolver(resolveRepoRoot, arguments)
}

func parseFixOptionsWithResolver(
	resolve repositoryRootResolver,
	arguments []string,
) (options fixOptions, err error) {
	const summary = "run safe style auto-fixes"
	var scope string
	flagSet := newFixFlagSet(&options, &scope)
	if err = parseArguments(flagSet, summary, arguments); err != nil {
		return options, err
	}

	options.scope, err = parseScope(scope)
	if err != nil {
		return options, err
	}

	options.repoRoot, err = resolve(options.repoRoot)
	return options, err
}

func newFixFlagSet(options *fixOptions, scope *string) (flagSet *flag.FlagSet) {
	flagSet = newFlagSet("fix")
	flagSet.StringVar(
		&options.repoRoot,
		"repo-root",
		"",
		"repository root (auto-detected when omitted)",
	)
	flagSet.StringVar(scope, "scope", "", "configured scope (profile default when omitted)")
	return flagSet
}

func fixUsageText() (usage string) {
	const summary = "run safe style auto-fixes"
	var options fixOptions
	var scope string
	return commandUsage("fix", summary, newFixFlagSet(&options, &scope))
}

/* --------------------------------------- Rule Selection --------------------------------------- */

func fixableRules(
	available []contract.Rule,
	context runner.Context,
) (rules []contract.Rule) {
	for _, rule := range available {
		if !context.Profile.Repository.HasScopeOverlap(context.Scope, rule.Scope) {
			continue
		}

		if rule.Fix.Empty() {
			continue
		}

		rules = append(rules, rule)
	}

	return rules
}
