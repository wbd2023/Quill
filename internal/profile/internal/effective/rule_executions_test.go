package effective_test

import (
	"testing"

	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/profile/internal/effective"
	"ciphera/tools/internal/profile/internal/profiletest"
	"ciphera/tools/internal/style"
)

/* --------------------------------------- Rule Executions -------------------------------------- */

func TestCompileRejectsIncompleteFileCommandExecution(t *testing.T) {
	t.Parallel()

	err := compileRuleDefinition(t, style.RuleDefinition{
		ID: "test/bad-file-command",
		Check: style.FileCommandExecution{
			ToolID: profiletest.Tool,
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

func TestCompileRejectsUnknownExecutionDetail(t *testing.T) {
	// The sealed Template interface prevents external types from
	// satisfying it, so the default case in the validator switch is
	// unreachable from outside the style package. This test documents
	// that the guard exists but cannot be exercised from external tests.
	t.Skip("sealed interface prevents constructing unknown detail types")
}

func TestCompileRejectsBlankRuleToolReference(t *testing.T) {
	t.Parallel()

	err := compileRuleDefinition(t, style.RuleDefinition{
		ID: "test/blank-tool",
		Check: style.ToolchainExecution{
			ToolIDs: []string{" "},
		},
	})
	requireErrorContains(t, err, "empty tool ID")
}

func TestCompileRejectsDuplicateRuleToolReference(t *testing.T) {
	t.Parallel()

	err := compileRuleDefinition(t, style.RuleDefinition{
		ID: "test/duplicate-tool",
		Check: style.ToolchainExecution{
			ToolIDs: []string{
				profiletest.Tool,
				profiletest.Tool,
			},
		},
	})
	requireErrorContains(t, err, "duplicates tool")
}

func TestCompileRejectsUnknownRuleToolReference(t *testing.T) {
	t.Parallel()

	err := compileRuleDefinition(t, style.RuleDefinition{
		ID: "test/unknown-tool",
		Check: style.ToolchainExecution{
			ToolIDs: []string{"unknown"},
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

	config := profiletest.Config()
	config.Rules = []policy.RuleBinding{
		{
			RuleID:         definition.ID,
			Enforcement:    style.EnforcementRequired,
			Scope:          config.Repository.DefaultScope,
			RequirementIDs: []string{profiletest.Requirement},
		},
	}
	config.Tools = []policy.PinnedTool{
		{ID: profiletest.Tool, Version: "1.0.0"},
	}

	_, err = effective.Compile(config, style.Definitions{
		ToolIDs: []string{profiletest.Tool},
		Rules:   []style.RuleDefinition{definition},
	})
	return err
}
