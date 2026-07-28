package profile

import (
	"testing"

	"github.com/wbd2023/quill/internal/policy"
	"github.com/wbd2023/quill/internal/profile/internal/profiletest"
	"github.com/wbd2023/quill/internal/style"
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
	requireErrorContainsInternal(t, err, "must define a file set")
}

func TestCompileRejectsMissingRuleCheck(t *testing.T) {
	t.Parallel()

	err := compileRuleDefinition(t, style.RuleDefinition{
		ID: "test/missing-check",
	})
	requireErrorContainsInternal(t, err, "must define check execution")
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
	requireErrorContainsInternal(t, err, "empty tool ID")
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
	requireErrorContainsInternal(t, err, "duplicates tool")
}

func TestCompileRejectsUnknownRuleToolReference(t *testing.T) {
	t.Parallel()

	err := compileRuleDefinition(t, style.RuleDefinition{
		ID: "test/unknown-tool",
		Check: style.ToolchainExecution{
			ToolIDs: []string{"unknown"},
		},
	})
	requireErrorContainsInternal(t, err, "references unknown tool")
}

func TestCompileCarriesPackIDProvenance(t *testing.T) {
	t.Parallel()

	definition := style.RuleDefinition{
		ID:     "test/provenance",
		PackID: "test",
		Name:   "Provenance rule",
		Group:  "test",
		Check: style.ProfileExecution{
			PackID: "test",
			Check:  "commands",
		},
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

	plan, err := compilePlan(config, style.Definitions{
		ToolIDs: []string{profiletest.Tool},
		Rules:   []style.RuleDefinition{definition},
	})
	if err != nil {
		t.Fatalf("compilePlan: %v", err)
	}
	if len(plan.Rules) != 1 {
		t.Fatalf("rules = %d, want 1", len(plan.Rules))
	}

	rule := plan.Rules[0]
	if rule.PackID != "test" {
		t.Fatalf("rule PackID = %q, want test", rule.PackID)
	}

	job, ok := rule.Check.(style.ProfileExecution)
	if !ok {
		t.Fatalf("expected ProfileExecution, got %T", rule.Check)
	}

	if job.PackID != "test" {
		t.Fatalf("job PackID = %q, want test", job.PackID)
	}
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

	_, err = compilePlan(config, style.Definitions{
		ToolIDs: []string{profiletest.Tool},
		Rules:   []style.RuleDefinition{definition},
	})
	return err
}
