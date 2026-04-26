package runner

import (
	"errors"
	"fmt"
	"sort"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/toolchain"
)

/* ------------------------------------------- Errors ------------------------------------------- */

var (
	errRuleBlocked   = errors.New("rule blocked by toolchain")
	errRuleViolation = errors.New("rule violations found")
)

type Executor func(
	context Context,
	spec contract.ExecutionSpec,
	toolStatuses map[string]toolchain.Status,
) (result contract.ExecutionResult, err error)

type ExecutorRegistry map[contract.ExecutorKind]Executor

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
	toolStatuses map[string]toolchain.Status,
	executors ExecutorRegistry,
) (result contract.ExecutionResult, err error) {
	return runRuleSpec(rule.ID, rule.Spec, rule.ToolIDs(), context, toolStatuses, executors)
}

func RunFix(
	rule contract.Rule,
	context Context,
	toolStatuses map[string]toolchain.Status,
	executors ExecutorRegistry,
) (result contract.ExecutionResult, err error) {
	return runRuleSpec(rule.ID, rule.FixSpec, rule.FixToolIDs(), context, toolStatuses, executors)
}

func runRuleSpec(
	ruleID string,
	spec contract.ExecutionSpec,
	toolIDs []string,
	context Context,
	toolStatuses map[string]toolchain.Status,
	executors ExecutorRegistry,
) (result contract.ExecutionResult, err error) {
	if spec.Empty() {
		return contract.ExecutionResult{}, nil
	}

	if len(toolIDs) > 0 && !toolchain.AllToolsValid(toolIDs, toolStatuses) {
		return contract.ExecutionResult{
			Diagnostics: []contract.Diagnostic{
				{
					Code:    "toolchain/blocked",
					Message: toolchain.ExplainToolIssues(toolIDs, toolStatuses),
				},
			},
		}, errRuleBlocked
	}

	executor, found := executors[spec.Kind]
	if !found {
		return contract.ExecutionResult{}, fmt.Errorf(
			"rule %s uses unknown executor %q",
			ruleID,
			spec.Executor(),
		)
	}

	return executor(context, spec, toolStatuses)
}

func ToolchainExecutor(
	_ Context,
	spec contract.ExecutionSpec,
	toolStatuses map[string]toolchain.Status,
) (result contract.ExecutionResult, err error) {
	detail, found := spec.ToolchainExecution()
	if !found {
		return contract.ExecutionResult{}, fmt.Errorf("toolchain executor received empty spec")
	}

	diagnostics := make([]contract.Diagnostic, 0, len(detail.ToolIDs))
	foundFailure := false
	for _, toolID := range detail.ToolIDs {
		status, found := toolStatuses[toolID]
		if !found || status.Valid {
			continue
		}

		foundFailure = true
		diagnostics = append(diagnostics, contract.Diagnostic{
			Code:    "toolchain/invalid",
			Message: toolchain.ExplainToolIssues([]string{toolID}, toolStatuses),
		})
	}

	if !foundFailure {
		return contract.ExecutionResult{}, nil
	}

	return contract.ExecutionResult{Diagnostics: diagnostics}, errRuleViolation
}
