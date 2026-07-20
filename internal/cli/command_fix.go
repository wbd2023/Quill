package cli

import (
	"context"
	"flag"

	"github.com/wbd2023/Quill/internal/engine"
	"github.com/wbd2023/Quill/internal/report"
)

/* ----------------------------------------- Fix Command ---------------------------------------- */

func runFix(tool Tool, options fixOptions) (exitCode int) {
	fixer, err := engine.New(options.repoRoot)
	if err != nil {
		tool.writeError(err)
		return 1
	}

	result, err := fixer.Fix(context.Background(), engine.FixOptions{
		Scope: options.scope,
	})
	if err != nil {
		tool.writeError(err)
		return 1
	}

	if len(result.Rules) == 0 {
		return 0
	}

	toolchainResult := report.ToolchainResult{Statuses: result.Toolchain.Statuses}
	if _, err := renderToolchainStatus(
		tool.stderr, report.FormatText, toolchainResult,
	); err != nil {
		tool.writeError(err)
		return 1
	}

	if !result.Toolchain.AllValid {
		return 1
	}

	for _, ruleResult := range result.Rules {
		if ruleResult.ExecutionError != nil {
			if ruleResult.Execution.Output != "" {
				tool.writeCommandOutput(ruleResult.Execution.Output)
			} else {
				tool.writeError(ruleResult.ExecutionError)
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
