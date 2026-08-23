package profile_test

import (
	"testing"

	"github.com/wbd2023/quill/internal/profile"
	"github.com/wbd2023/quill/internal/profile/internal/profiletest"
	"github.com/wbd2023/quill/internal/style"
)

/* --------------------------------------- Rule Executions -------------------------------------- */

func TestCompileRejectsIncompleteFileCommandExecution(t *testing.T) {
	t.Parallel()

	err := compileRuleDefinition(t, style.RuleDefinition{
		ID: "test/bad-file-command",
		Check: style.FileCommand{
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
		Check: style.ToolchainCheck{
			ToolIDs: []string{" "},
		},
	})
	requireErrorContains(t, err, "empty tool ID")
}

func TestCompileRejectsDuplicateRuleToolReference(t *testing.T) {
	t.Parallel()

	err := compileRuleDefinition(t, style.RuleDefinition{
		ID: "test/duplicate-tool",
		Check: style.ToolchainCheck{
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
		Check: style.ToolchainCheck{
			ToolIDs: []string{"unknown"},
		},
	})
	requireErrorContains(t, err, "references unknown tool")
}

func TestCompileCarriesPackIDProvenance(t *testing.T) {
	t.Parallel()

	definition := style.RuleDefinition{
		ID:     "test/provenance",
		PackID: "test",
		Name:   "Provenance rule",
		Group:  "test",
		Check:  style.ProfileCheck{Check: "commands"},
	}

	config := profiletest.Config()
	config.Rules = []profile.RuleBinding{
		{
			RuleID:         definition.ID,
			Enforcement:    style.EnforcementRequired,
			Scope:          config.Repository.DefaultScope,
			RequirementIDs: []string{profiletest.Requirement},
		},
	}

	plan, err := profile.Compile(config, style.Definitions{
		ToolIDs: []string{profiletest.Tool},
		Rules:   []style.RuleDefinition{definition},
	})
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}
	if len(plan.Rules) != 1 {
		t.Fatalf("rules = %d, want 1", len(plan.Rules))
	}

	// Pack provenance is carried by the RuleDefinition/Rule alone; execution values carry none.
	rule := plan.Rules[0]
	if rule.PackID != "test" {
		t.Fatalf("rule PackID = %q, want test", rule.PackID)
	}

	job, ok := rule.Check.(style.ProfileCheck)
	if !ok {
		t.Fatalf("expected ProfileCheck, got %T", rule.Check)
	}
	if job.Check != "commands" {
		t.Fatalf("job Check = %q, want commands", job.Check)
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
	config.Rules = []profile.RuleBinding{
		{
			RuleID:         definition.ID,
			Enforcement:    style.EnforcementRequired,
			Scope:          config.Repository.DefaultScope,
			RequirementIDs: []string{profiletest.Requirement},
		},
	}
	config.Tools = []profile.PinnedTool{
		{ID: profiletest.Tool, Version: "1.0.0"},
	}

	_, err = profile.Compile(config, style.Definitions{
		ToolIDs: []string{profiletest.Tool},
		Rules:   []style.RuleDefinition{definition},
	})
	return err
}
