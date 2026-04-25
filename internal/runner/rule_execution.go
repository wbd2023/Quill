package runner

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/runtime"
)

/* ------------------------------------------- Errors ------------------------------------------- */

var (
	errRuleBlocked   = errors.New("rule blocked by toolchain")
	errRuleViolation = errors.New("rule violations found")
)

type Executor func(
	context Context,
	spec contract.ExecutionSpec,
	toolStatuses map[string]runtime.ToolStatus,
) (output string, err error)

type ExecutorRegistry map[string]Executor

func IsBlocked(err error) (blocked bool) {
	return errors.Is(err, errRuleBlocked)
}

func ToolIDsForRules(rules []contract.Rule) (toolIDs []string) {
	seen := make(map[string]bool)
	for _, rule := range rules {
		for _, toolID := range rule.ToolIDs() {
			if seen[toolID] {
				continue
			}

			seen[toolID] = true
			toolIDs = append(toolIDs, toolID)
		}
	}

	sort.Strings(toolIDs)
	return toolIDs
}

func ToolIDsForFixes(rules []contract.Rule) (toolIDs []string) {
	seen := make(map[string]bool)
	for _, rule := range rules {
		for _, toolID := range rule.FixToolIDs() {
			if seen[toolID] {
				continue
			}

			seen[toolID] = true
			toolIDs = append(toolIDs, toolID)
		}
	}

	sort.Strings(toolIDs)
	return toolIDs
}

func RunRule(
	rule contract.Rule,
	context Context,
	toolStatuses map[string]runtime.ToolStatus,
	executors ExecutorRegistry,
) (output string, err error) {
	return runRuleSpec(rule.ID, rule.Spec, rule.ToolIDs(), context, toolStatuses, executors)
}

func RunFix(
	rule contract.Rule,
	context Context,
	toolStatuses map[string]runtime.ToolStatus,
	executors ExecutorRegistry,
) (output string, err error) {
	return runRuleSpec(rule.ID, rule.FixSpec, rule.FixToolIDs(), context, toolStatuses, executors)
}

func runRuleSpec(
	ruleID string,
	spec contract.ExecutionSpec,
	toolIDs []string,
	context Context,
	toolStatuses map[string]runtime.ToolStatus,
	executors ExecutorRegistry,
) (output string, err error) {
	if spec.Executor == "" {
		return "", nil
	}

	if len(toolIDs) > 0 && !runtime.AllToolsValid(toolIDs, toolStatuses) {
		return runtime.ExplainToolIssues(toolIDs, toolStatuses), errRuleBlocked
	}

	executor, found := executors[spec.Executor]
	if !found {
		return "", fmt.Errorf("rule %s uses unknown executor %q", ruleID, spec.Executor)
	}

	return executor(context, spec, toolStatuses)
}

func ToolchainExecutor(
	_ Context,
	spec contract.ExecutionSpec,
	toolStatuses map[string]runtime.ToolStatus,
) (output string, err error) {
	var builder strings.Builder
	foundFailure := false
	for _, toolID := range spec.ToolIDs {
		status, found := toolStatuses[toolID]
		if !found || status.Valid {
			continue
		}

		foundFailure = true
		builder.WriteString(runtime.ExplainToolIssues([]string{toolID}, toolStatuses))
		builder.WriteString("\n")
	}

	if !foundFailure {
		return "", nil
	}

	return strings.TrimSpace(builder.String()), errRuleViolation
}
