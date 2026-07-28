package cli

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/wbd2023/quill/internal/engine"
	"github.com/wbd2023/quill/internal/report"
)

/* --------------------------------------- Explain Command -------------------------------------- */

func runExplain(ctx context.Context, tool Tool, options explainOptions) (exitCode int) {
	engineInstance, err := tool.buildEngine(options.repoRoot)
	if err != nil {
		return tool.reportCommandError(ctx, "explain", options.format, err)
	}

	snapshot, err := engineInstance.Metadata(ctx)
	if err != nil {
		return tool.reportCommandError(ctx, "explain", options.format, err)
	}

	explanation, err := buildExplanation(snapshot, options.subject)
	if err != nil {
		return tool.reportCommandError(ctx, "explain", options.format, err)
	}

	if err := report.WriteExplain(tool.stdout, "explain", options.format, explanation); err != nil {
		return tool.reportCommandError(ctx, "explain", options.format, err)
	}

	return 0
}

func parseExplainOptionsWithResolver(
	resolve repositoryRootResolver,
	arguments []string,
) (options explainOptions, err error) {
	const summary = "explain an active rule"
	var format string
	flagSet := newExplainFlagSet(&options, &format)

	positional, err := parseFlags(flagSet, summary, arguments)
	if err != nil {
		return options, err
	}

	options.format, err = parseFormat(format)
	if err != nil {
		return options, err
	}

	options.subject, err = parseExplainSubject(positional)
	if err != nil {
		return options, err
	}

	options.repoRoot, err = resolve(options.repoRoot)
	return options, err
}

// parseExplainSubject accepts exactly one subject in "<kind>:<id>" form. Only the "rule" kind is
// supported for the MVP; an empty or missing kind or id is an invalid-argument error.
func parseExplainSubject(positional []string) (subject string, err error) {
	if len(positional) != 1 {
		return "", fmt.Errorf(
			"expected one subject (rule:<id>), got %s",
			describePositional(positional),
		)
	}

	subject = positional[0]
	kind, id, found := strings.Cut(subject, ":")
	if !found {
		return "", fmt.Errorf("invalid subject %q: expected rule:<id>", subject)
	}

	if kind != "rule" {
		return "", fmt.Errorf("unsupported subject %q: only rule:<id> is supported", subject)
	}

	if strings.TrimSpace(id) == "" {
		return "", fmt.Errorf("invalid subject %q: rule id must not be empty", subject)
	}

	return subject, nil
}

func newExplainFlagSet(options *explainOptions, format *string) (flagSet *flag.FlagSet) {
	flagSet = newFlagSet("explain")
	flagSet.StringVar(
		&options.repoRoot,
		"repo-root",
		"",
		"repository root (auto-detected when omitted)",
	)
	flagSet.StringVar(format, "format", string(report.FormatText), "format: text|json")
	return flagSet
}

func explainUsageText() (usage string) {
	const summary = "explain an active rule"
	var options explainOptions
	var format string
	return commandUsage("explain", summary, newExplainFlagSet(&options, &format))
}

func explainMachineMode(arguments []string) (requested bool) {
	var options explainOptions
	var format string
	return machineModeRequested(arguments, newExplainFlagSet(&options, &format), &format)
}

/* ------------------------------------ Explanation Building ------------------------------------ */

// explainRuleID extracts the rule identifier from a validated "rule:<id>" subject.
func explainRuleID(subject string) (ruleID string) {
	_, id, _ := strings.Cut(subject, ":")
	return id
}

func buildExplanation(
	snapshot engine.MetadataSnapshot,
	subject string,
) (explanation report.ExplainResult, err error) {
	ruleID := explainRuleID(subject)

	for _, rule := range snapshot.Rules {
		if rule.ID != ruleID {
			continue
		}

		if !rule.Active {
			return report.ExplainResult{}, fmt.Errorf(
				"rule %q is declared by pack %q but is not active in this profile",
				ruleID,
				rule.PackID,
			)
		}

		return report.ExplainResult{Rule: mapExplainRule(rule, snapshot)}, nil
	}

	return report.ExplainResult{}, fmt.Errorf(
		"unknown rule %q: no matching rule is declared",
		ruleID,
	)
}

func mapExplainRule(
	rule engine.RuleMetadata,
	snapshot engine.MetadataSnapshot,
) (explained report.ExplainRule) {
	explained = report.ExplainRule{
		ID:           rule.ID,
		Pack:         rule.PackID,
		Name:         rule.Name,
		Group:        string(rule.Group),
		External:     packExternal(rule.PackID, snapshot),
		Enforcement:  string(rule.Enforcement),
		Scope:        string(rule.Scope),
		Requirements: rule.RequirementIDs,
		Check:        mapExecution(rule.Check),
	}

	if rule.HasFix {
		fix := mapExecution(rule.Fix)
		explained.Fix = &fix
	}

	if config, found := snapshot.PackConfigs.Lookup(rule.PackID); found {
		explained.PackConfig = map[string]any(config)
	}

	return explained
}

func mapExecution(detail engine.ExecutionDetail) (execution report.ExplainExecution) {
	return report.ExplainExecution{
		Category: detail.Category,
		Tools:    detail.ToolIDs,
		FileSet:  detail.FileSet,
		Language: detail.Language,
		Detail:   detail.Detail,
	}
}

func packExternal(packID string, snapshot engine.MetadataSnapshot) (external bool) {
	for _, pack := range snapshot.Packs {
		if pack.ID == packID {
			return pack.External
		}
	}

	return false
}
