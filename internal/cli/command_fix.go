package cli

import (
	"context"
	"flag"

	"github.com/wbd2023/quill/internal/engine"
	"github.com/wbd2023/quill/internal/report"
)

/* ----------------------------------------- Fix Command ---------------------------------------- */

func runFix(ctx context.Context, tool Tool, options fixOptions) (exitCode int) {
	engineInstance, err := tool.buildEngine(options.repoRoot)
	if err != nil {
		return tool.reportCommandError(ctx, "fix", options.format, err)
	}

	result, err := engineInstance.Fix(ctx, engine.FixOptions{Scope: options.scope})
	if err != nil {
		return tool.reportCommandError(ctx, "fix", options.format, err)
	}

	if options.format == report.FormatJSON {
		return renderFixResult(tool, "fix", result)
	}

	return renderFixText(tool, result)
}

func renderFixText(tool Tool, result engine.FixResult) (exitCode int) {
	if len(result.Rules) == 0 {
		return 0
	}

	toolchainResult := report.ToolchainResult{Statuses: result.Toolchain.Statuses}
	if _, err := renderToolchainStatus(
		tool.stderr, "fix", report.FormatText, toolchainResult,
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

// renderFixResult writes the machine-mode fix envelope and maps blocking outcomes to a nonzero
// exit status. An invalid toolchain or a per-fixer execution error completes the operation
// (status "ok") but returns exit 1, matching the documented "status ok + nonzero exit" rule.
func renderFixResult(tool Tool, command string, result engine.FixResult) (exitCode int) {
	view := report.NewFixView(
		result.Scope,
		result.Toolchain.AllValid,
		result.Toolchain.Statuses,
		toFixEntries(result.Rules),
	)

	if err := report.WriteFix(tool.stdout, command, view); err != nil {
		tool.writeError(err)
		return 1
	}

	if !result.Toolchain.AllValid || view.HasExecutionError() {
		return 1
	}

	return 0
}

func toFixEntries(rules []engine.RuleFixResult) (entries []report.FixEntry) {
	entries = make([]report.FixEntry, 0, len(rules))
	for _, rule := range rules {
		entries = append(entries, report.FixEntry{
			Rule:           report.NewRuleSummary(rule.Rule),
			Execution:      rule.Execution,
			ExecutionError: rule.ExecutionError,
		})
	}

	return entries
}

/* --------------------------------------- Option Parsing --------------------------------------- */

func parseFixOptions(arguments []string) (options fixOptions, err error) {
	return parseFixOptionsWithResolver(resolveRepoRoot, arguments)
}

func parseFixOptionsWithResolver(
	resolve repositoryRootResolver,
	arguments []string,
) (options fixOptions, err error) {
	const summary = "run safe style auto-fixes"
	var scope string
	var format string
	flagSet := newFixFlagSet(&options, &scope, &format)
	if err = parseArguments(flagSet, summary, arguments); err != nil {
		return options, err
	}

	options.scope, err = parseScope(scope)
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

func newFixFlagSet(options *fixOptions, scope *string, format *string) (flagSet *flag.FlagSet) {
	flagSet = newFlagSet("fix")
	flagSet.StringVar(
		&options.repoRoot,
		"repo-root",
		"",
		"repository root (auto-detected when omitted)",
	)
	flagSet.StringVar(scope, "scope", "", "configured scope (profile default when omitted)")
	flagSet.StringVar(format, "format", string(report.FormatText), "format: text|json")
	return flagSet
}

func fixUsageText() (usage string) {
	const summary = "run safe style auto-fixes"
	var options fixOptions
	var scope string
	var format string
	return commandUsage("fix", summary, newFixFlagSet(&options, &scope, &format))
}
