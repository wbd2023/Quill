package cli

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/wbd2023/quill/internal/engine"
	"github.com/wbd2023/quill/internal/report"
)

/* ---------------------------------------- List Command ---------------------------------------- */

func runList(ctx context.Context, tool Tool, options listOptions) (exitCode int) {
	engineInstance, err := tool.buildEngine(options.repoRoot)
	if err != nil {
		return tool.reportCommandError(ctx, "list", options.format, err)
	}

	snapshot, err := engineInstance.Metadata(ctx)
	if err != nil {
		return tool.reportCommandError(ctx, "list", options.format, err)
	}

	result := newListResult(snapshot, options.selector)
	if err := report.WriteList(tool.stdout, "list", options.format, result); err != nil {
		return tool.reportCommandError(ctx, "list", options.format, err)
	}

	return 0
}

func parseListOptionsWithResolver(
	resolve repositoryRootResolver,
	arguments []string,
) (options listOptions, err error) {
	const summary = "list packs, rules, tools, or scopes"
	var format string
	flagSet := newListFlagSet(&options, &format)

	positional, err := parseFlags(flagSet, summary, arguments)
	if err != nil {
		return options, err
	}

	options.format, err = parseFormat(format)
	if err != nil {
		return options, err
	}

	options.selector, err = parseListSelector(positional)
	if err != nil {
		return options, err
	}

	options.repoRoot, err = resolve(options.repoRoot)
	return options, err
}

func parseListSelector(positional []string) (selector string, err error) {
	if len(positional) != 1 {
		return "", fmt.Errorf(
			"expected one selector (packs|rules|tools|scopes), got %s",
			describePositional(positional),
		)
	}

	selector = positional[0]
	if !report.IsValidListSelector(selector) {
		return "", fmt.Errorf(
			"invalid selector %q: must be one of packs, rules, tools, scopes",
			selector,
		)
	}

	return selector, nil
}

func newListFlagSet(options *listOptions, format *string) (flagSet *flag.FlagSet) {
	flagSet = newFlagSet("list")
	flagSet.StringVar(
		&options.repoRoot,
		"repo-root",
		"",
		"repository root (auto-detected when omitted)",
	)
	flagSet.StringVar(format, "format", string(report.FormatText), "format: text|json")
	return flagSet
}

func listUsageText() (usage string) {
	const summary = "list packs, rules, tools, or scopes"
	var options listOptions
	var format string
	return commandUsage("list", summary, newListFlagSet(&options, &format))
}

func listMachineMode(arguments []string) (requested bool) {
	var options listOptions
	var format string
	return machineModeRequested(arguments, newListFlagSet(&options, &format), &format)
}

/* -------------------------------------- Snapshot Mapping -------------------------------------- */

func newListResult(snapshot engine.MetadataSnapshot, selector string) (result report.ListResult) {
	result = report.ListResult{Selector: selector}
	switch selector {
	case report.ListPacks:
		result.Packs = mapListPacks(snapshot.Packs)
	case report.ListRules:
		result.Rules = mapListRules(snapshot.Rules)
	case report.ListTools:
		result.Tools = mapListTools(snapshot.Tools)
	case report.ListScopes:
		result.Scopes = mapListScopes(snapshot.Scopes)
	}

	return result
}

func mapListPacks(packs []engine.PackMetadata) (rows []report.ListPack) {
	rows = make([]report.ListPack, 0, len(packs))
	for _, pack := range packs {
		rows = append(rows, report.ListPack{
			ID:       pack.ID,
			Name:     pack.Name,
			Active:   pack.Active,
			External: pack.External,
			Rules:    len(pack.RuleIDs),
			Tools:    len(pack.ToolIDs),
		})
	}

	return rows
}

func mapListRules(rules []engine.RuleMetadata) (rows []report.ListRule) {
	rows = make([]report.ListRule, 0, len(rules))
	for _, rule := range rules {
		row := report.ListRule{
			ID:     rule.ID,
			Pack:   rule.PackID,
			Name:   rule.Name,
			Active: rule.Active,
			Fix:    rule.HasFix,
		}
		if rule.Active {
			row.Enforcement = string(rule.Enforcement)
			row.Scope = string(rule.Scope)
		}
		rows = append(rows, row)
	}

	return rows
}

func mapListTools(tools []engine.ToolMetadata) (rows []report.ListTool) {
	rows = make([]report.ListTool, 0, len(tools))
	for _, tool := range tools {
		rows = append(rows, report.ListTool{
			ID:       tool.ID,
			Name:     tool.Name,
			Command:  tool.Command,
			Pin:      tool.PinnedVersion,
			Packs:    tool.PackIDs,
			External: tool.External,
		})
	}

	return rows
}

func mapListScopes(scopes []engine.ScopeMetadata) (rows []report.ListScope) {
	rows = make([]report.ListScope, 0, len(scopes))
	for _, scope := range scopes {
		rows = append(rows, report.ListScope{
			Name:    string(scope.Name),
			Roots:   scope.Roots,
			Default: scope.Default,
		})
	}

	return rows
}

func describePositional(positional []string) (description string) {
	if len(positional) == 0 {
		return "none"
	}

	return strings.Join(positional, ", ")
}
