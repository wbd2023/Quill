package effective_test

import (
	"testing"

	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/profile/internal/effective"
	"ciphera/tools/internal/profile/internal/fixture"
	"ciphera/tools/internal/style"
)

/* --------------------------------------- Rule Executions -------------------------------------- */

func TestCompileRejectsIncompleteFileCommandExecution(t *testing.T) {
	t.Parallel()

	err := compileRuleDefinition(t, style.RuleDefinition{
		ID: "test/bad-file-command",
		Check: style.ExecutionSpec{
			Kind: style.ExecutionFileCommand,
			Detail: style.FileCommandExecution{
				ToolID: fixture.Tool,
			},
		},
	})
	requireErrorContains(t, err, "must define a file set")
}

func TestCompileRejectsMissingRuleCheck(t *testing.T) {
	t.Parallel()

	err := compileRuleDefinition(t, style.RuleDefinition{
		ID: "test/missing-check",
	})
	requireErrorContains(t, err, "must define check execution")
}

func TestCompileRejectsMismatchedExecutionKind(t *testing.T) {
	t.Parallel()

	err := compileRuleDefinition(t, style.RuleDefinition{
		ID: "test/mismatched-kind",
		Check: style.ExecutionSpec{
			Kind: style.ExecutionFileCommand,
			Detail: style.ToolchainExecution{
				ToolIDs: []string{fixture.Tool},
			},
		},
	})
	requireErrorContains(t, err, `expected "toolchain"`)
}

func TestCompileRejectsBlankRuleToolReference(t *testing.T) {
	t.Parallel()

	err := compileRuleDefinition(t, style.RuleDefinition{
		ID: "test/blank-tool",
		Check: style.ExecutionSpec{
			Kind: style.ExecutionToolchain,
			Detail: style.ToolchainExecution{
				ToolIDs: []string{" "},
			},
		},
	})
	requireErrorContains(t, err, "empty tool ID")
}

func TestCompileRejectsDuplicateRuleToolReference(t *testing.T) {
	t.Parallel()

	err := compileRuleDefinition(t, style.RuleDefinition{
		ID: "test/duplicate-tool",
		Check: style.ExecutionSpec{
			Kind: style.ExecutionToolchain,
			Detail: style.ToolchainExecution{
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

	err := compileRuleDefinition(t, style.RuleDefinition{
		ID: "test/unknown-tool",
		Check: style.ExecutionSpec{
			Kind: style.ExecutionToolchain,
			Detail: style.ToolchainExecution{
				ToolIDs: []string{"unknown"},
			},
		},
	})
	requireErrorContains(t, err, "references unknown tool")
}

/* ------------------------------------------- Support ------------------------------------------ */

func compileRuleDefinition(t *testing.T, definition style.RuleDefinition) (err error) {
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
			Enforcement:    style.EnforcementRequired,
			Scope:          config.Repository.DefaultScope,
			RequirementIDs: []string{fixture.Requirement},
		},
	}
	config.Tools = []policy.PinnedTool{
		{ID: fixture.Tool, Version: "1.0.0"},
	}

	_, err = effective.Compile(config, style.Definitions{
		Tools: []style.Tool{
			{ID: fixture.Tool, Name: "Test tool"},
		},
		Rules: []style.RuleDefinition{definition},
	})
	return err
}
