package effective_test

import (
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/profile/effective"
	"ciphera/tools/internal/profile/internal/fixture"
)

/* --------------------------------------- Rule Executions -------------------------------------- */

func TestCompileRejectsIncompleteFileCommandExecution(t *testing.T) {
	t.Parallel()

	err := compileRuleDefinition(t, contract.RuleDefinition{
		ID: "test/bad-file-command",
		Check: contract.ExecutionSpec{
			Kind: contract.ExecutorFileCommand,
			Detail: contract.FileCommandExecution{
				ToolID: fixture.Tool,
			},
		},
	})
	requireErrorContains(t, err, "must define a file set")
}

func TestCompileRejectsMissingRuleCheck(t *testing.T) {
	t.Parallel()

	err := compileRuleDefinition(t, contract.RuleDefinition{
		ID: "test/missing-check",
	})
	requireErrorContains(t, err, "must define check execution")
}

func TestCompileRejectsMismatchedExecutionKind(t *testing.T) {
	t.Parallel()

	err := compileRuleDefinition(t, contract.RuleDefinition{
		ID: "test/mismatched-kind",
		Check: contract.ExecutionSpec{
			Kind: contract.ExecutorFileCommand,
			Detail: contract.ToolchainExecution{
				ToolIDs: []string{fixture.Tool},
			},
		},
	})
	requireErrorContains(t, err, `expected "toolchain"`)
}

func TestCompileRejectsBlankRuleToolReference(t *testing.T) {
	t.Parallel()

	err := compileRuleDefinition(t, contract.RuleDefinition{
		ID: "test/blank-tool",
		Check: contract.ExecutionSpec{
			Kind: contract.ExecutorToolchain,
			Detail: contract.ToolchainExecution{
				ToolIDs: []string{" "},
			},
		},
	})
	requireErrorContains(t, err, "empty tool ID")
}

func TestCompileRejectsDuplicateRuleToolReference(t *testing.T) {
	t.Parallel()

	err := compileRuleDefinition(t, contract.RuleDefinition{
		ID: "test/duplicate-tool",
		Check: contract.ExecutionSpec{
			Kind: contract.ExecutorToolchain,
			Detail: contract.ToolchainExecution{
				ToolIDs: []string{
					fixture.Tool,
					fixture.Tool,
				},
			},
		},
	})
	requireErrorContains(t, err, "duplicates tool")
}

func TestCompileRejectsUnknownRuleToolReference(t *testing.T) {
	t.Parallel()

	err := compileRuleDefinition(t, contract.RuleDefinition{
		ID: "test/unknown-tool",
		Check: contract.ExecutionSpec{
			Kind: contract.ExecutorToolchain,
			Detail: contract.ToolchainExecution{
				ToolIDs: []string{"unknown"},
			},
		},
	})
	requireErrorContains(t, err, "references unknown tool")
}

/* ------------------------------------------- Support ------------------------------------------ */

func compileRuleDefinition(t *testing.T, definition contract.RuleDefinition) (err error) {
	t.Helper()

	if definition.Name == "" {
		definition.Name = "Test rule"
	}
	if definition.Group == "" {
		definition.Group = "test"
	}

	config := fixture.Config()
	config.Rules = []policy.RuleBinding{
		{
			RuleID:         definition.ID,
			Enforcement:    contract.EnforcementRequired,
			Scope:          config.Repository.DefaultScope,
			RequirementIDs: []string{fixture.Requirement},
		},
	}
	config.Tools = []policy.PinnedTool{
		{ID: fixture.Tool, Version: "1.0.0"},
	}

	_, err = effective.Compile(config, contract.Definitions{
		Tools: []contract.Tool{
			{ID: fixture.Tool, Name: "Test tool"},
		},
		Rules: []contract.RuleDefinition{definition},
	})
	return err
}
